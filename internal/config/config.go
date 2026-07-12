package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Port       string
	TempDir    string
	StaticDir  string
	FFmpegPath string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	workingDir, _ := os.Getwd()
	tempDir := os.Getenv("TEMP_DIR")
	if tempDir == "" {
		tempDir = filepath.Join(workingDir, "temp")
	} else if !filepath.IsAbs(tempDir) {
		tempDir = filepath.Join(workingDir, tempDir)
	}

	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = workingDir
	}
	if !filepath.IsAbs(staticDir) {
		staticDir = filepath.Join(workingDir, staticDir)
	}

	ffmpegPath := os.Getenv("FFMPEG_PATH")
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}

	return Config{
		Port:       port,
		TempDir:    tempDir,
		StaticDir:  staticDir,
		FFmpegPath: ffmpegPath,
	}
}
