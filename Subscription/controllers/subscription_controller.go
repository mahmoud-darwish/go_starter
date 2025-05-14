package controllers

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
	"starter/Subscription/models"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type SubscriptionController struct {
	DB *gorm.DB
}

func NewSubscriptionController(db *gorm.DB) *SubscriptionController {
	return &SubscriptionController{DB: db}
}

func (ctrl *SubscriptionController) Create(w http.ResponseWriter, r *http.Request) {
	var input models.SubscriptionCreateRequestDTO 
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subscription := models.Subscription{
		//ID:        uuid.New(),
		UserID:    input.UserID,
		ChannelID: input.ChannelID,
	}

	if err := ctrl.DB.Create(&subscription).Error; err != nil {
		http.Error(w, "Failed to create subscription", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.SubscriptionResponseDTOFromModel(subscription))
}

func (ctrl *SubscriptionController) FindAll(w http.ResponseWriter, r *http.Request) {
	var subscriptions []models.Subscription
	if err := ctrl.DB.Find(&subscriptions).Error; err != nil {
		http.Error(w, "Failed to retrieve subscriptions", http.StatusInternalServerError)
		return
	}

	var response []models.SubscriptionResponseDTO
	for _, n := range subscriptions {
		response = append(response, models.SubscriptionResponseDTOFromModel(n))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (ctrl *SubscriptionController) FindByID(w http.ResponseWriter, r *http.Request) {
	//idStr :=  // Or use chi.URLParam if using Chi routing
	idStr := chi.URLParam(r, "id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}
	id := uint(idUint)

	var subscription models.Subscription
	if err := ctrl.DB.First(&subscription, "id = ?", id).Error; err != nil {
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.SubscriptionResponseDTOFromModel(subscription))
}



func (ctrl *SubscriptionController) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}
	id := uint(idUint)

	if err := ctrl.DB.Delete(&models.Subscription{}, "id = ?", id).Error; err != nil {
		http.Error(w, "Failed to delete subscription", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
