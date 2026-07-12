package services

import (
	"os"
	"path/filepath"
	"testing"

	"price-comparator/internal/config"
)

func TestCleanupArtifactsRemovesUploadAndExport(t *testing.T) {
	tempDir := t.TempDir()
	cfg := config.Config{TempDir: tempDir, FFmpegPath: "ffmpeg"}
	service := New(cfg)

	uploadPath := filepath.Join(tempDir, "uploads", "sample.mp3")
	exportPath := filepath.Join(tempDir, "exports", "result.mp3")
	if err := os.MkdirAll(filepath.Dir(uploadPath), 0o755); err != nil {
		t.Fatalf("mkdir upload dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(exportPath), 0o755); err != nil {
		t.Fatalf("mkdir export dir: %v", err)
	}
	if err := os.WriteFile(uploadPath, []byte("audio"), 0o644); err != nil {
		t.Fatalf("write upload file: %v", err)
	}
	if err := os.WriteFile(exportPath, []byte("audio"), 0o644); err != nil {
		t.Fatalf("write export file: %v", err)
	}

	if err := service.CleanupArtifacts(uploadPath, "result.mp3"); err != nil {
		t.Fatalf("cleanup artifacts: %v", err)
	}

	if _, err := os.Stat(uploadPath); !os.IsNotExist(err) {
		t.Fatalf("expected upload file removed, stat err = %v", err)
	}
	if _, err := os.Stat(exportPath); !os.IsNotExist(err) {
		t.Fatalf("expected export file removed, stat err = %v", err)
	}
}
