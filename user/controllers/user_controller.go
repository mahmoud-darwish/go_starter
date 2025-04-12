package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"starter/auth"
	"starter/pkg/logger"
	"starter/pkg/utils"
	"starter/user/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB        *gorm.DB
	Validator *validator.Validate
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db, Validator: validator.New()}
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var reqDTO models.UserRegisterRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		log.Error().Err(err).Msg("Failed to decode register request body")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := c.Validator.Struct(reqDTO); err != nil {
		log.Error().Err(err).Msg("Validation failed")
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user := models.User{
		Email:     reqDTO.Email,
		Password:  reqDTO.Password,
		FirstName: reqDTO.FirstName,
		LastName:  reqDTO.LastName,
	}

	if err := c.DB.Create(&user).Error; err != nil {
		log.Error().Err(err).Str("email", user.Email).Msg("Failed to create user")
		// if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
		// 	utils.RespondWithError(w, http.StatusConflict, "Email already exists")
		// } else {
		// 	utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
		// }
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		log.Error().Err(err).Str("email", user.Email).Msg("Failed to generate token")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	respDTO := models.UserResponseDTO{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	response := map[string]interface{}{
		"user":  respDTO,
		"token": token,
	}
	utils.RespondWithJSON(w, http.StatusCreated, response)
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var reqDTO models.UserLoginRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		log.Error().Err(err).Msg("Failed to decode login request body")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := c.Validator.Struct(reqDTO); err != nil {
		log.Error().Err(err).Msg("Validation failed")
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if err := c.DB.Where("email = ?", reqDTO.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().Str("email", reqDTO.Email).Msg("User not found")
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		} else {
			log.Error().Err(err).Str("email", reqDTO.Email).Msg("Database error during user lookup")
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqDTO.Password)); err != nil {
		log.Warn().Err(err).Str("email", reqDTO.Email).Msg("Password verification failed")
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		log.Error().Err(err).Str("email", reqDTO.Email).Msg("Failed to generate JWT")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := map[string]string{
		"token": token,
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Err(err).Msg("Invalid user ID")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	if err := c.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().Int("id", id).Msg("User not found")
			utils.RespondWithError(w, http.StatusNotFound, "User not found")
		} else {
			log.Error().Err(err).Int("id", id).Msg("Database error")
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		}
		return
	}

	respDTO := models.UserResponseDTO{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	utils.RespondWithJSON(w, http.StatusOK, respDTO)
}

func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Err(err).Msg("Invalid user ID")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var reqDTO models.UserResponseDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		log.Error().Err(err).Msg("Failed to decode update request body")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var user models.User
	if err := c.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().Int("id", id).Msg("User not found")
			utils.RespondWithError(w, http.StatusNotFound, "User not found")
		} else {
			log.Error().Err(err).Int("id", id).Msg("Database error")
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		}
		return
	}

	user.Email = reqDTO.Email
	user.FirstName = reqDTO.FirstName
	user.LastName = reqDTO.LastName

	if err := c.DB.Save(&user).Error; err != nil {
		log.Error().Err(err).Int("id", id).Msg("Failed to update user")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	respDTO := models.UserResponseDTO{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	utils.RespondWithJSON(w, http.StatusOK, respDTO)
}

func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Err(err).Msg("Invalid user ID")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := c.DB.Delete(&models.User{}, id).Error; err != nil {
		log.Error().Err(err).Int("id", id).Msg("Failed to delete user")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
