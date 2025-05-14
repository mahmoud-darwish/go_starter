package controllers

import (
	"encoding/json"
	"net/http"

	"strconv"
	"gorm.io/gorm"
	"starter/Notifications/models"
	//"strconv"

	"github.com/go-chi/chi/v5"
)

type NotificationController struct {
	DB *gorm.DB
}

func NewNotificationController(db *gorm.DB) *NotificationController {
	return &NotificationController{DB: db}
}

func (ctrl *NotificationController) Create(w http.ResponseWriter, r *http.Request) {
	var input models.NotificationCreateRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	notification := models.Notification{
		//ID:        uuid.New(),
		UserID:    input.UserID,
		ChannelID: input.ChannelID,
		Content:   input.Content,
		Source:    input.Source,
	}

	if err := ctrl.DB.Create(&notification).Error; err != nil {
		http.Error(w, "Failed to create notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.NotificationResponseDTOFromModel(notification))
}

func (ctrl *NotificationController) FindAll(w http.ResponseWriter, r *http.Request) {
	var notifications []models.Notification
	if err := ctrl.DB.Find(&notifications).Error; err != nil {
		http.Error(w, "Failed to retrieve notifications", http.StatusInternalServerError)
		return
	}

	var response []models.NotificationResponseDTO
	for _, n := range notifications {
		response = append(response, models.NotificationResponseDTOFromModel(n))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (ctrl *NotificationController) FindByID(w http.ResponseWriter, r *http.Request) {
	//idStr :=  // Or use chi.URLParam if using Chi routing
	
	idStr := chi.URLParam(r, "id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}
	id := uint(idUint)

	var notification models.Notification
	if err := ctrl.DB.First(&notification, "id = ?", id).Error; err != nil {
		http.Error(w, "Notification not found", http.StatusNotFound)
		return
	}


	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.NotificationResponseDTOFromModel(notification))
}

func (ctrl *NotificationController) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}
	id := uint(idUint)

	var input models.NotificationUpdateRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var notification models.Notification
	if err := ctrl.DB.First(&notification, "id = ?", id).Error; err != nil {
		http.Error(w, "Notification not found", http.StatusNotFound)
		return
	}

	notification.Content = input.Content
	notification.Source = input.Source

	if err := ctrl.DB.Save(&notification).Error; err != nil {
		http.Error(w, "Failed to update notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.NotificationResponseDTOFromModel(notification))
}

func (ctrl *NotificationController) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}
	id := uint(idUint)

	if err := ctrl.DB.Delete(&models.Notification{}, "id = ?", id).Error; err != nil {
		http.Error(w, "Failed to delete notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
