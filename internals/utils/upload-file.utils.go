package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// SaveUploadedFile menyimpan file upload
func SaveUploadedFile(ctx *gin.Context, file *multipart.FileHeader, destDir string, filename string) (string, error) {
	// Validasi ekstensi
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", fmt.Errorf("invalid file type: only jpg, jpeg, png allowed")
	}

	// Validasi ukuran file (max 1 MB)
	const maxSize = 1 << 20 // 1 MB
	if file.Size > maxSize {
		return "", fmt.Errorf("file too large: maximum size is 1MB")
	}

	// Buat folder tujuan jika belum ada
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Nama file final
	finalName := fmt.Sprintf("%s%s", filename, ext)
	fullPath := filepath.Join(destDir, finalName)

	// Simpan file langsung
	if err := ctx.SaveUploadedFile(file, fullPath); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return finalName, nil
}
