package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"starter/Comments/models"
	//"strconv"
	"starter/auth"
	"github.com/go-chi/chi/v5"
)

type CommentsController struct {
	DB *gorm.DB
}

func NewCommentsController(db *gorm.DB) *CommentsController {
	return &CommentsController{DB: db}
}

func (ctrl *CommentsController) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uint)
	if !ok {
		http.Error(w, "Unauthorized: user ID not found", http.StatusUnauthorized)
		return
	}
	var input models.CommentCreateRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comment := models.Comment{
		//ID:        uuid.New(),
		UserID:    userID,
		VideoID: input.VideoID,
		Content:   input.Content,
	}

	if err := ctrl.DB.Create(&comment).Error; err != nil {
		http.Error(w, "Failed to create Comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.CommentResponseDTOFromModel(comment))
}

func (ctrl *CommentsController) FindAll(w http.ResponseWriter, r *http.Request) {
	var comments []models.Comment
	if err := ctrl.DB.Find(&comments).Error; err != nil {
		http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
		return
	}

	var response []models.CommentResponseDTO
	for _, n := range comments {
		response = append(response, models.CommentResponseDTOFromModel(n))
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (ctrl *CommentsController) FindByID(w http.ResponseWriter, r *http.Request) {
	//idStr :=  // Or use chi.URLParam if using Chi routing
	
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid comment ID from findbyid", http.StatusBadRequest)
		return
	}

	var comment models.Comment
	if err := ctrl.DB.First(&comment, "id = ?", id).Error; err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.CommentResponseDTOFromModel(comment))
}



