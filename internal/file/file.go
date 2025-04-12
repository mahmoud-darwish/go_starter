package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"starter/pkg/logger"
)

type UploadService struct {
	UploadDir string
}

func NewUploadService() *UploadService {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	return &UploadService{UploadDir: uploadDir}
}

func (s *UploadService) UploadFile(file io.Reader, filename string, prefix string) (string, error) {
	log := logger.GetLogger()

	storedFilename := fmt.Sprintf("%s_%s", prefix, filename)
	filepath := filepath.Join(s.UploadDir, storedFilename)

	dst, err := os.Create(filepath)
	if err != nil {
		log.Error().Err(err).Str("filename", storedFilename).Msg("Failed to create file")
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		log.Error().Err(err).Str("filename", storedFilename).Msg("Failed to write file")
		return "", err
	}

	log.Info().Str("filename", storedFilename).Msg("File uploaded successfully")
	return storedFilename, nil
}
