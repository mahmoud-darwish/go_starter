package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	ChannelModels "starter/channel/models"
	"starter/internal/cache"
	"starter/internal/file"
	"starter/pkg/logger"
	"starter/pkg/utils"
	"starter/video/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type VideoController struct {
	DB        *gorm.DB
	Validator *validator.Validate
	UploadSvc *file.UploadService
}

func NewVideoController(db *gorm.DB) *VideoController {
	return &VideoController{
		DB:        db,
		Validator: validator.New(),
		UploadSvc: file.NewUploadService(),
	}
}

func (c *VideoController) Create(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var reqDTO models.VideoCreateRequestDTO

	// err := r.ParseMultipartForm(50 << 20)
	// if err != nil {
	// 	log.Error().Err(err).Msg("Failed to parse multipart form")
	// 	utils.RespondWithError(w, http.StatusBadRequest, "Failed to parse form")
	// 	return
	// }

	// title := r.FormValue("title")
	// description := r.FormValue("description")
	// genre := r.FormValue("genre")
	// channelID, err := strconv.ParseUint(r.FormValue("channel_id"), 10, 32)
	// if err != nil {
	// 	log.Error().Err(err).Msg("Invalid channel ID")
	// 	utils.RespondWithError(w, http.StatusBadRequest, "Invalid channel ID")
	// 	return
	// }

	// videoType, err := strconv.Atoi(r.FormValue("type"))
	// if err != nil {
	// 	log.Error().Err(err).Msg("Invalid video type")
	// 	utils.RespondWithError(w, http.StatusBadRequest, "Invalid video type")
	// 	return
	// }
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

	video := models.Video{
		ID:          55,
		ChannelID:   reqDTO.ChannelID,
		Type:        reqDTO.Type,
		Title:       reqDTO.Title,
		Description: reqDTO.Description,
		TimeStamp:   time.Now(),
		Genre:       reqDTO.Genre,
	}

	if err := c.DB.Create(&video).Error; err != nil {
		log.Error().Err(err).Str("title", video.Title).Msg("Failed to create video")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create video")
		return
	}

	respDTO := models.VideoResponseDTO{
		ID:          video.ID,
		ChannelID:   video.ChannelID,
		Type:        video.Type,
		Path:        video.Path,
		Title:       video.Title,
		Description: video.Description,
		TimeStamp:   video.TimeStamp,
		Genre:       video.Genre,
	}

	utils.RespondWithJSON(w, http.StatusCreated, respDTO)
}

func (c *VideoController) ProcessVideoUpload(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	videoID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Err(err).Msg("Invalid video ID")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid video ID")
		return
	}

	var video models.Video
	if err := c.DB.First(&video, videoID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().Int("id", videoID).Msg("Video not found")
			utils.RespondWithError(w, http.StatusNotFound, "Video not found")
		} else {
			log.Error().Err(err).Int("id", videoID).Msg("Database error")
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		}
		return
	}

	err = r.ParseMultipartForm(500 << 20) // 500 MB limit
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse multipart form")
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	file, handler, err := r.FormFile("video")
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve file")
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to retrieve file")
		return
	}
	defer file.Close()

	// Upload the video with progress tracking
	path, err := c.UploadSvc.UploadFile(file, handler.Filename, fmt.Sprintf("video_%d", videoID))
	if err != nil {
		log.Error().Err(err).Int("id", videoID).Msg("Failed to upload video")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to upload video")
		return
	}

	// Update video path in database
	video.Path = path
	if err := c.DB.Save(&video).Error; err != nil {
		log.Error().Err(err).Int("id", videoID).Msg("Failed to update video")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update video")
		return
	}

	respDTO := models.VideoResponseDTO{
		ID:          video.ID,
		ChannelID:   video.ChannelID,
		Type:        video.Type,
		Path:        video.Path,
		Title:       video.Title,
		Description: video.Description,
		TimeStamp:   video.TimeStamp,
		Genre:       video.Genre,
	}

	utils.RespondWithJSON(w, http.StatusOK, respDTO)
}

func (c *VideoController) Get(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Err(err).Msg("Invalid video ID")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid video ID")
		return
	}

	cacheKey := fmt.Sprintf("video:%d", id)
	ctx := context.Background()

	cached, err := cache.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var respDTO models.VideoResponseDTO
		if err := json.Unmarshal([]byte(cached), &respDTO); err != nil {
			log.Error().Err(err).Int("id", id).Msg("Failed to unmarshal cached video")
		} else {
			log.Info().Int("id", id).Msg("Video retrieved from cache")
			utils.RespondWithJSON(w, http.StatusOK, respDTO)
			return
		}
	}

	var video models.Video
	if err := c.DB.First(&video, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().Int("id", id).Msg("Video not found")
			utils.RespondWithError(w, http.StatusNotFound, "Video not found")
		} else {
			log.Error().Err(err).Int("id", id).Msg("Database error")
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		}
		return
	}

	respDTO := models.VideoResponseDTO{
		ID:          video.ID,
		ChannelID:   video.ChannelID,
		Type:        video.Type,
		Path:        video.Path,
		Title:       video.Title,
		Description: video.Description,
		TimeStamp:   video.TimeStamp,
		Genre:       video.Genre,
	}

	jsonData, err := json.Marshal(respDTO)
	if err == nil {
		err = cache.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute).Err()
		if err != nil {
			log.Warn().Err(err).Int("id", id).Msg("Failed to cache video")
		}
	}

	utils.RespondWithJSON(w, http.StatusOK, respDTO)
}

func (c *VideoController) GetAll(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()

	// Parse pagination parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 20 // Default limit
	}

	offset := (page - 1) * limit

	// Check if we can get from cache
	cacheKey := fmt.Sprintf("videos:page:%d:limit:%d", page, limit)
	ctx := context.Background()

	cached, err := cache.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var respDTO models.VideoListResponseDTO
		if err := json.Unmarshal([]byte(cached), &respDTO); err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal cached videos")
		} else {
			log.Info().Int("page", page).Int("limit", limit).Msg("Videos retrieved from cache")
			utils.RespondWithJSON(w, http.StatusOK, respDTO)
			return
		}
	}

	var videos []models.Video
	var total int64

	// Count total videos
	if err := c.DB.Model(&models.Video{}).Count(&total).Error; err != nil {
		log.Error().Err(err).Msg("Failed to count videos")
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Get paginated videos
	if err := c.DB.Limit(limit).Offset(offset).Order("created_at DESC").Find(&videos).Error; err != nil {
		log.Error().Err(err).Msg("Failed to retrieve videos")
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Convert to response DTOs
	var videoDTOs []models.VideoResponseDTO
	for _, video := range videos {
		videoDTOs = append(videoDTOs, models.VideoResponseDTO{
			ID:          video.ID,
			ChannelID:   video.ChannelID,
			Type:        video.Type,
			Path:        video.Path,
			Title:       video.Title,
			Description: video.Description,
			TimeStamp:   video.TimeStamp,
			Genre:       video.Genre,
		})
	}

	respDTO := models.VideoListResponseDTO{
		Videos: videoDTOs,
		Total:  total,
	}

	// Cache the results
	jsonData, err := json.Marshal(respDTO)
	if err == nil {
		err = cache.RedisClient.Set(ctx, cacheKey, jsonData, 5*time.Minute).Err()
		if err != nil {
			log.Warn().Err(err).Msg("Failed to cache videos list")
		}
	}

	utils.RespondWithJSON(w, http.StatusOK, respDTO)
}

func (c *VideoController) Update(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Err(err).Msg("Invalid video ID")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid video ID")
		return
	}

	var reqDTO models.VideoUpdateRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		log.Error().Err(err).Msg("Failed to decode update video request body")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := c.Validator.Struct(reqDTO); err != nil {
		log.Error().Err(err).Msg("Validation failed")
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var video models.Video
	if err := c.DB.First(&video, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().Int("id", id).Msg("Video not found")
			utils.RespondWithError(w, http.StatusNotFound, "Video not found")
		} else {
			log.Error().Err(err).Int("id", id).Msg("Database error")
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
		}
		return
	}

	// Update fields if provided
	if reqDTO.Title != "" {
		video.Title = reqDTO.Title
	}
	if reqDTO.Description != "" {
		video.Description = reqDTO.Description
	}
	if reqDTO.Genre != "" {
		video.Genre = reqDTO.Genre
	}
	if reqDTO.Type != 0 {
		video.Type = reqDTO.Type
	}

	if err := c.DB.Save(&video).Error; err != nil {
		log.Error().Err(err).Int("id", id).Msg("Failed to update video")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update video")
		return
	}

	// Invalidate cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("video:%d", id)
	cache.RedisClient.Del(ctx, cacheKey)

	respDTO := models.VideoResponseDTO{
		ID:          video.ID,
		ChannelID:   video.ChannelID,
		Type:        video.Type,
		Path:        video.Path,
		Title:       video.Title,
		Description: video.Description,
		TimeStamp:   video.TimeStamp,
		Genre:       video.Genre,
	}

	utils.RespondWithJSON(w, http.StatusOK, respDTO)
}

func (c *VideoController) Delete(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Err(err).Msg("Invalid video ID")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid video ID")
		return
	}

	// Get video first to get the file path
	var video models.Video
	if err := c.DB.First(&video, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().Int("id", id).Msg("Video not found")
			utils.RespondWithError(w, http.StatusNotFound, "Video not found")
			return
		} else {
			log.Error().Err(err).Int("id", id).Msg("Database error")
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
			return
		}
	}

	// Delete the video file if it exists
	if video.Path != "" {
		filePath := filepath.Join(c.UploadSvc.UploadDir, video.Path)
		if _, err := os.Stat(filePath); err == nil {
			if err := os.Remove(filePath); err != nil {
				log.Warn().Err(err).Str("path", filePath).Msg("Failed to delete video file")
			}
		}
	}

	// Delete from database
	if err := c.DB.Delete(&models.Video{}, id).Error; err != nil {
		log.Error().Err(err).Int("id", id).Msg("Failed to delete video")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete video")
		return
	}

	// Invalidate cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("video:%d", id)
	cache.RedisClient.Del(ctx, cacheKey)

	// Invalidate list cache keys (simple approach)
	cache.RedisClient.Del(ctx, "videos:page:*")

	w.WriteHeader(http.StatusNoContent)
}

// ...existing code...

// Response struct combining video and channel info
type VideoWithChannelResponse struct {
    Video   models.Video              `json:"video"`
    Channel ChannelModels.ChannelResponseDTO `json:"channel"`
}

// Handler to get video with channel info
func (c *VideoController) GetWithChannel(w http.ResponseWriter, r *http.Request) {
    log := logger.GetLogger()
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        log.Error().Err(err).Msg("Invalid video ID")
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid video ID")
        return
    }

    // Fetch video
    var video models.Video
    if err := c.DB.First(&video, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            utils.RespondWithError(w, http.StatusNotFound, "Video not found")
        } else {
            utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
        }
        return
    }

    // Fetch channel
    var channel ChannelModels.Channel
    if err := c.DB.First(&channel, video.ChannelID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            utils.RespondWithError(w, http.StatusNotFound, "Channel not found")
        } else {
            utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
        }
        return
    }

    resp := VideoWithChannelResponse{
        Video:   video,
        Channel: ChannelModels.ChannelResponseDTOFromModel(channel),
    }

    utils.RespondWithJSON(w, http.StatusOK, resp)
}

// ...existing code...



// Handler to get all videos with their channel info
func (c *VideoController) GetAllWithChannel(w http.ResponseWriter, r *http.Request) {
    log := logger.GetLogger()

    var videos []models.Video
    if err := c.DB.Find(&videos).Error; err != nil {
        log.Error().Err(err).Msg("Failed to retrieve videos")
        utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
        return
    }

    var result []VideoWithChannelResponse
    for _, video := range videos {
        var channel ChannelModels.Channel
        if err := c.DB.First(&channel, video.ChannelID).Error; err != nil {
            // If channel not found, skip this video or handle as needed
            continue
        }
        result = append(result, VideoWithChannelResponse{
            Video:   video,
            Channel: ChannelModels.ChannelResponseDTOFromModel(channel),
        })
    }

    utils.RespondWithJSON(w, http.StatusOK, result)
}

func (c *VideoController) GetAllByChannelID(w http.ResponseWriter, r *http.Request) {
    log := logger.GetLogger()
    channelIDStr := chi.URLParam(r, "channelID")
    channelID, err := strconv.Atoi(channelIDStr)
    if err != nil {
        log.Error().Err(err).Msg("Invalid channel ID")
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid channel ID")
        return
    }

    var videos []models.Video
    if err := c.DB.Where("channel_id = ?", channelID).Find(&videos).Error; err != nil {
        log.Error().Err(err).Msg("Failed to retrieve videos")
        utils.RespondWithError(w, http.StatusInternalServerError, "Database error")
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, videos)
}