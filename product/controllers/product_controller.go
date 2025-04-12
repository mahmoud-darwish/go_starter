package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"starter/internal/cache"
	"starter/internal/file"
	"starter/pkg/logger"
	"starter/pkg/utils"
	"starter/product/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ProductController struct {
	DB        *gorm.DB
	Validator *validator.Validate
	UploadSvc *file.UploadService
}

func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{
		DB:        db,
		Validator: validator.New(),
		UploadSvc: file.NewUploadService(),
	}
}

func (c *ProductController) Create(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var reqDTO models.ProductCreateRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		log.Error().Err(err).Msg("Failed to decode create product request body")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := c.Validator.Struct(reqDTO); err != nil {
		log.Error().Err(err).Msg("Validation failed")
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	product := models.Product{
		Name:        reqDTO.Name,
		Description: reqDTO.Description,
		Price:       reqDTO.Price,
	}

	if err := c.DB.Create(&product).Error; err != nil {
		log.Error().Err(err).Str("name", product.Name).Msg("Failed to create product")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	respDTO := models.ProductResponseDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	utils.RespondWithJSON(w, http.StatusCreated, respDTO)
}

func (c *ProductController) Get(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Err(err).Msg("Invalid product ID")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	cacheKey := fmt.Sprintf("product:%d", id)
	ctx := context.Background()

	cached, err := cache.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var respDTO models.ProductResponseDTO
		if err := json.Unmarshal([]byte(cached), &respDTO); err != nil {
			log.Error().Err(err).Int("id", id).Msg("Failed to unmarshal cached product")
		} else {
			log.Info().Int("id", id).Msg("Product retrieved from cache")
			utils.RespondWithJSON(w, http.StatusOK, respDTO)
			return
		}
	}

	var product models.Product
	if err := c.DB.First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().Int("id", id).Msg("Product not found")
			utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		} else {
			log.Error().Err(err).Int("id", id).Msg("Database error")
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		}
		return
	}

	respDTO := models.ProductResponseDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Image:       product.Image,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	jsonData, err := json.Marshal(respDTO)
	if err == nil {
		err = cache.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute).Err()
		if err != nil {
			log.Warn().Err(err).Int("id", id).Msg("Failed to cache product")
		}
	}

	utils.RespondWithJSON(w, http.StatusOK, respDTO)
}

func (c *ProductController) UploadImage(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Err(err).Msg("Invalid product ID")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
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

	var product models.Product
	if err := c.DB.First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().Int("id", id).Msg("Product not found")
			utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		} else {
			log.Error().Err(err).Int("id", id).Msg("Database error")
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		}
		return
	}

	filename, err := c.UploadSvc.UploadFile(file, handler.Filename, fmt.Sprintf("product_%d", id))
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("Failed to upload file")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to upload file")
		return
	}

	product.Image = filename
	if err := c.DB.Save(&product).Error; err != nil {
		log.Error().Err(err).Int("id", id).Msg("Failed to update product")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	respDTO := models.ProductResponseDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Image:       product.Image,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	utils.RespondWithJSON(w, http.StatusOK, respDTO)
}

func (c *ProductController) Update(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Err(err).Msg("Invalid product ID")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var reqDTO models.ProductUpdateRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		log.Error().Err(err).Msg("Failed to decode update product request body")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := c.Validator.Struct(reqDTO); err != nil {
		log.Error().Err(err).Msg("Validation failed")
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var product models.Product
	if err := c.DB.First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().Int("id", id).Msg("Product not found")
			utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		} else {
			log.Error().Err(err).Int("id", id).Msg("Database error")
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		}
		return
	}

	if reqDTO.Name != "" {
		product.Name = reqDTO.Name
	}
	if reqDTO.Description != "" {
		product.Description = reqDTO.Description
	}
	if reqDTO.Price > 0 {
		product.Price = reqDTO.Price
	}

	if err := c.DB.Save(&product).Error; err != nil {
		log.Error().Err(err).Int("id", id).Msg("Failed to update product")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	respDTO := models.ProductResponseDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	utils.RespondWithJSON(w, http.StatusOK, respDTO)
}

func (c *ProductController) Delete(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Err(err).Msg("Invalid product ID")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	if err := c.DB.Delete(&models.Product{}, id).Error; err != nil {
		log.Error().Err(err).Int("id", id).Msg("Failed to delete product")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
