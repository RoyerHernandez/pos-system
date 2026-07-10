package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const maxImageSize = 2 * 1024 * 1024 // 2MB

var allowedImageTypes = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

var allowedMIMETypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

func SaveImage(file multipart.File, header *multipart.FileHeader, basePath string, subDir string) (string, error) {
	if header.Size > maxImageSize {
		return "", errors.New("image size exceeds 2MB limit")
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedImageTypes[ext] {
		return "", errors.New("only jpg, jpeg, and png files are allowed")
	}

	// Validate actual file content type
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file header: %w", err)
	}
	contentType := http.DetectContentType(buf[:n])
	if !allowedMIMETypes[contentType] {
		return "", errors.New("file content is not a valid JPEG or PNG image")
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("failed to reset file reader: %w", err)
	}

	baseDir := "../pos/views/img"
	destDir := filepath.Join(baseDir, basePath, subDir)

	// Path traversal protection
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve base directory: %w", err)
	}
	absDest, err := filepath.Abs(destDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve destination directory: %w", err)
	}
	if !strings.HasPrefix(absDest, absBase) {
		return "", errors.New("invalid upload path")
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	destPath := filepath.Join(destDir, filename)

	// Verify final path is still within base directory
	absPath, err := filepath.Abs(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve file path: %w", err)
	}
	if !strings.HasPrefix(absPath, absBase) {
		return "", errors.New("invalid file path")
	}

	dst, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return filename, nil
}
