package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"price-comparator/internal/config"
	"price-comparator/internal/ffmpeg"
	"price-comparator/internal/models"
	"price-comparator/internal/utils"
)

type AudioService struct {
	cfg config.Config
	mu  sync.RWMutex

	exportSources map[string]string
}

func New(cfg config.Config) *AudioService {
	return &AudioService{cfg: cfg, exportSources: make(map[string]string)}
}

func (s *AudioService) CleanupArtifacts(uploadPath, exportName string) error {
	if uploadPath != "" {
		_ = os.Remove(uploadPath)
	}
	if exportName != "" {
		exportPath := filepath.Join(s.cfg.TempDir, "exports", exportName)
		_ = os.Remove(exportPath)
	}
	return nil
}

func (s *AudioService) CleanupDownload(exportName string) error {
	s.mu.RLock()
	uploadPath := s.exportSources[exportName]
	s.mu.RUnlock()

	s.mu.Lock()
	delete(s.exportSources, exportName)
	s.mu.Unlock()

	return s.CleanupArtifacts(uploadPath, exportName)
}

func (s *AudioService) StoreUpload(file multipart.File, header *multipart.FileHeader) (string, string, int64, error) {
	if header == nil {
		return "", "", 0, fmt.Errorf("file tidak ditemukan")
	}
	if !strings.EqualFold(filepath.Ext(header.Filename), ".mp3") {
		return "", "", 0, fmt.Errorf("format file harus mp3")
	}
	if header.Size > 50<<20 {
		return "", "", 0, fmt.Errorf("ukuran file maksimal 50MB")
	}
	if err := utils.EnsureDir(filepath.Join(s.cfg.TempDir, "uploads")); err != nil {
		return "", "", 0, err
	}

	storedName := utils.BuildStoredName(header.Filename)
	storedPath := filepath.Join(s.cfg.TempDir, "uploads", storedName)

	out, err := os.Create(storedPath)
	if err != nil {
		return "", "", 0, err
	}
	defer out.Close()

	written, err := io.Copy(out, file)
	if err != nil {
		return "", "", 0, err
	}

	return strings.TrimSuffix(storedName, filepath.Ext(storedName)), storedPath, written, nil
}

func (s *AudioService) ExportFile(inputPath string, req models.ExportRequest) (string, error) {
	if err := utils.EnsureDir(filepath.Join(s.cfg.TempDir, "exports")); err != nil {
		return "", err
	}

	outputName := fmt.Sprintf("%d.mp3", time.Now().UnixNano())
	outputPath := filepath.Join(s.cfg.TempDir, "exports", outputName)
	if err := ffmpeg.Exec(inputPath, outputPath, req, s.cfg.FFmpegPath); err != nil {
		return "", err
	}

	s.mu.Lock()
	s.exportSources[outputName] = inputPath
	s.mu.Unlock()
	return outputName, nil
}
