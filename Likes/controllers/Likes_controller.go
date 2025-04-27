package controllers

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
	"starter/Likes/models"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type LikeController struct {
	DB *gorm.DB
}

func NewLikeController(db *gorm.DB) *LikeController {
	return &LikeController{DB: db}
}

func (ctrl *LikeController) Create(w http.ResponseWriter, r *http.Request) {
	var input models.LikeCreateRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	like := models.Like{
		//ID:        uuid.New(),
		UserID:    input.UserID,
		VideoID: input.VideoID,
		
	}

	if err := ctrl.DB.Create(&like).Error; err != nil {
		http.Error(w, "Failed to create like", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.LikeResponseDTOFromModel(like))
}

func (ctrl *LikeController) FindAll(w http.ResponseWriter, r *http.Request) {
	var likes []models.Like
	if err := ctrl.DB.Find(&likes).Error; err != nil {
		http.Error(w, "Failed to retrieve likes", http.StatusInternalServerError)
		return
	}

	var response []models.LikeResponseDTO
	for _, n := range likes {
		response = append(response, models.LikeResponseDTOFromModel(n))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (ctrl *LikeController) FindByID(w http.ResponseWriter, r *http.Request) {
	//idStr :=  // Or use chi.URLParam if using Chi routing
	
	idStr := chi.URLParam(r, "id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}
	id := uint(idUint)

	var like models.Like
	if err := ctrl.DB.First(&like, "id = ?", id).Error; err != nil {
		http.Error(w, "Like not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.LikeResponseDTOFromModel(like))
}



func (ctrl *LikeController) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}
	id := uint(idUint)

	if err := ctrl.DB.Delete(&models.Like{}, "id = ?", id).Error; err != nil {
		http.Error(w, "Failed to delete like", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
