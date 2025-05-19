package controllers

import (
	"encoding/json"
	"net/http"
	"github.com/go-playground/validator/v10"
	"starter/internal/file"
	"starter/pkg/utils"
	"starter/pkg/logger"
	"strconv"
	"gorm.io/gorm"
	"starter/channel/models"
	//"strconv"
	"starter/auth"
	"github.com/go-chi/chi/v5"

	

	"fmt"



	
)

type ChannelController struct {
	DB *gorm.DB
	Validator *validator.Validate
	UploadSvc *file.UploadService
}

func NewChannelController(db *gorm.DB) *ChannelController {
		return &ChannelController{
		DB:        db,
		Validator: validator.New(),
		UploadSvc: file.NewUploadService(),
	}
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

func (c *ChannelController) UploadImage(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Err(err).Msg("Invalid channel ID")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid channel ID")
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse multipart form")
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve file")
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to retrieve file")
		return
	}
	defer file.Close()

	var channel models.Channel
	if err := c.DB.First(&channel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().Int("id", id).Msg("Channel not found")
			utils.RespondWithError(w, http.StatusNotFound, "Channel not found")
		} else {
			log.Error().Err(err).Int("id", id).Msg("Database error")
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		}
		return
	}

	filename, err := c.UploadSvc.UploadFile(file, handler.Filename, fmt.Sprintf("channel_%d", id))
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("Failed to upload file")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to upload file")
		return
	}

	channel.Logo = filename
	if err := c.DB.Save(&channel).Error; err != nil {
		log.Error().Err(err).Int("id", id).Msg("Failed to update channel")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update channel")
		return
	}

	respDTO := models.ChannelResponseDTO{
		ID:          channel.ID,
		UserID: 	channel.UserID,
		Name:        channel.Name,
		Logo:       channel.Logo,
		Bio:       channel.Bio,
		TimeStamp:   channel.TimeStamp,

	}

	utils.RespondWithJSON(w, http.StatusOK, respDTO)
}
// func (c *ChannelController) GetChannelId(w http.ResponseWriter, r *http.Request) {
// 	userID, ok := r.Context().Value(auth.UserIDKey).(uint)
// 	if !ok {
// 		http.Error(w, "Unauthorized: user ID not found", http.StatusUnauthorized)
// 		return
// 	}

// 	var channel models.Channel
// 	if err := c.DB.First(&channel, "user_id = ?", userID).Error; err ==err  {
// 		if err == err {
// 			// No channel found for this user
// 			resp := map[string]int{"channel_id": -1}
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(resp)
// 			return
// 		}
// 		// Other DB error
// 		http.Error(w, "Database error", http.StatusInternalServerError)
// 		return
// 	}

// 	// Return just the channel ID
// 	resp := map[string]uint{"channel_id": channel.ID}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(resp)
// }

func (c *ChannelController) GetChannelId(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uint)
	if !ok {
		http.Error(w, "Unauthorized: user ID not found", http.StatusUnauthorized)
		return
	}

	var channel models.Channel
	if err := c.DB.First(&channel, "user_id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// No channel found for this user
			resp := map[string]int{"channel_id": -1}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		// Other DB error
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Return just the channel ID
	resp := map[string]uint{"channel_id": channel.ID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}


