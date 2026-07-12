package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"price-comparator/internal/config"
	"price-comparator/internal/models"
	"price-comparator/internal/services"
)

type Handler struct {
	service *services.AudioService
	cfg     config.Config
}

func New(cfg config.Config, service *services.AudioService) *Handler {
	return &Handler{service: service, cfg: cfg}
}

func (h *Handler) Register(r *gin.Engine) {
	r.GET("/", h.index)
	r.POST("/upload", h.upload)
	r.POST("/export", h.export)
	r.POST("/cleanup", h.cleanup)
	r.GET("/download/:name", h.download)
}

func (h *Handler) index(c *gin.Context) {
	c.File(filepath.Join(h.cfg.StaticDir, "index.html"))
}

func (h *Handler) upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "upload invalid"})
		return
	}
	defer file.Close()

	id, storedPath, size, err := h.service.StoreUpload(file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.UploadResponse{
		ID:         id,
		FileName:   header.Filename,
		StoredName: storedPath,
		Size:       size,
	})
}

func (h *Handler) export(c *gin.Context) {
	var req models.ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inputPath := req.InputPath
	if inputPath == "" {
		inputPath = c.PostForm("inputPath")
	}
	if inputPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "inputPath required"})
		return
	}

	outputName, err := h.service.ExportFile(inputPath, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.ExportResponse{Success: true, DownloadURL: fmt.Sprintf("/download/%s", outputName), Message: "Export selesai"})
}

func (h *Handler) download(c *gin.Context) {
	name := c.Param("name")
	if strings.Contains(name, "..") || strings.Contains(name, string(filepath.Separator)) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	path := filepath.Join(h.cfg.TempDir, "exports", name)
	if _, err := os.Stat(path); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.File(path)
	_ = h.service.CleanupDownload(name)
}

func (h *Handler) cleanup(c *gin.Context) {
	var req struct {
		UploadPath string `json:"uploadPath"`
		ExportName string `json:"exportName"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		req.UploadPath = c.Query("uploadPath")
		req.ExportName = c.Query("exportName")
	}
	_ = h.service.CleanupArtifacts(req.UploadPath, req.ExportName)
	c.JSON(http.StatusOK, gin.H{"success": true})
}
