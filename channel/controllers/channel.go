package controllers

import (
	"encoding/json"
	"net/http"

	"strconv"
	"gorm.io/gorm"
	"starter/channel/models"
	//"strconv"
	"starter/auth"
	"github.com/go-chi/chi/v5"
)

type ChannelController struct {
	DB *gorm.DB
}

func NewChannelController(db *gorm.DB) *ChannelController {
	return &ChannelController{DB: db}
}

func (ctrl *ChannelController) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uint)
	if !ok {
		http.Error(w, "Unauthorized: user ID not found", http.StatusUnauthorized)
		return
	}
	var input models.ChannelCreateRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	channel := models.Channel{
	
		UserID:    userID,
		Name: input.Name,
		Logo:   input.Logo,
		Bio:    input.Bio,
	}

	if err := ctrl.DB.Create(&channel).Error; err != nil {
		http.Error(w, "Failed to create channel", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.ChannelResponseDTOFromModel(channel))
}

func (ctrl *ChannelController) FindAll(w http.ResponseWriter, r *http.Request) {
	var channel []models.Channel
	if err := ctrl.DB.Find(&channel).Error; err != nil {
		http.Error(w, "Failed to retrieve channels", http.StatusInternalServerError)
		return
	}

	var response []models.ChannelResponseDTO
	for _, n := range channel {
		response = append(response, models.ChannelResponseDTOFromModel(n))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (ctrl *ChannelController) FindByID(w http.ResponseWriter, r *http.Request) {
	//idStr :=  // Or use chi.URLParam if using Chi routing
	
	idStr := chi.URLParam(r, "id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}
	id := uint(idUint)

	var channel models.Channel
	if err := ctrl.DB.First(&channel, "id = ?", id).Error; err != nil {
		http.Error(w, "Channel not found", http.StatusNotFound)
		return
	}


	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.ChannelResponseDTOFromModel(channel))
}

func (ctrl *ChannelController) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}
	id := uint(idUint)

	var input models.ChannelUpdateRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var channel models.Channel
	if err := ctrl.DB.First(&channel, "id = ?", id).Error; err != nil {
		http.Error(w, "Channel not found", http.StatusNotFound)
		return
	}

	channel.Name = input.Name
	channel.Logo = input.Logo
	channel.Bio = input.Bio



	if err := ctrl.DB.Save(&channel).Error; err != nil {
		http.Error(w, "Failed to update channel", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.ChannelResponseDTOFromModel(channel))
}

func (ctrl *ChannelController) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}
	id := uint(idUint)

	if err := ctrl.DB.Delete(&models.Channel{}, "id = ?", id).Error; err != nil {
		http.Error(w, "Failed to delete channel", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
