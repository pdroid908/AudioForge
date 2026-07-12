package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

func BuildStoredName(original string) string {
	ext := strings.ToLower(filepath.Ext(original))
	if ext == "" {
		ext = ".mp3"
	}
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
}
