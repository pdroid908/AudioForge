package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	"price-comparator/handlers"
	"price-comparator/internal/config"
	"price-comparator/internal/services"
)

func cleanExpiredFiles(baseDir string) {
	for _, dirName := range []string{"uploads", "exports"} {
		dirPath := filepath.Join(baseDir, dirName)
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			fullPath := filepath.Join(dirPath, entry.Name())
			info, err := entry.Info()
			if err != nil {
				continue
			}
			if time.Since(info.ModTime()) > 10*time.Minute {
				_ = os.Remove(fullPath)
			}
		}
	}
}

func main() {
	cfg := config.Load()

	if err := os.MkdirAll(cfg.TempDir, 0o755); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(cfg.TempDir, "uploads"), 0o755); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(cfg.TempDir, "exports"), 0o755); err != nil {
		log.Fatal(err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	service := services.New(cfg)
	h := handlers.New(cfg, service)
	h.Register(r)

	go func() {
		for {
			cleanExpiredFiles(cfg.TempDir)
			time.Sleep(60 * time.Second)
		}
	}()

	log.Printf("server listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

