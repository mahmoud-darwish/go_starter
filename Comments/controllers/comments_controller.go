package controllers

import (
	"encoding/json"
	"net/http"


	"gorm.io/gorm"
	"starter/Comments/models"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CommentsController struct {
	DB *gorm.DB
}

func NewCommentsController(db *gorm.DB) *CommentsController {
	return &CommentsController{DB: db}
}

func (ctrl *CommentsController) Create(w http.ResponseWriter, r *http.Request) {
	var input models.CommentCreateRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comment := models.Comment{
		//ID:        uuid.New(),
		UserID:    input.UserID,
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
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}
	id := uint(idUint)

	var comment models.Comment
	if err := ctrl.DB.First(&comment, "id = ?", id).Error; err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.CommentResponseDTOFromModel(comment))
}



