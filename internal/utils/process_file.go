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

	"github.com/google/uuid"
)

var allowExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

var allowMimeTypes = map[string]bool{
	"image/jpeg": true, // .jpeg & .jpg
	"image/png":  true,
}

const maxSize = 5 << 20

func ValidateAndSaveFile(FileHeader *multipart.FileHeader, uploadDir string) (string, error) {

	// Check extension in filename
	ext := strings.ToLower(filepath.Ext(FileHeader.Filename))
	if !allowExts[ext] {
		return "", errors.New("unsupported file extension")
	}

	// Check the file size
	if FileHeader.Size > maxSize {
		return "", fmt.Errorf("file too large (max %dMB)", maxSize/(1<<20))
	}

	// Check the file type
	file, err := FileHeader.Open()
	if err != nil {
		return "", errors.New("cannot open file")
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", errors.New("cannot read file")
	}

	mimeType := http.DetectContentType(buffer)
	if !allowMimeTypes[mimeType] {
		return "", fmt.Errorf("invalid MIME type: %s", mimeType)
	}

	// Change filename (abc.jpg)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Create folder if not exist
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", errors.New("cannot create upload folder")
	}

	// uploadDir "./upload" + filename ("abc.jpg")
	savePath := filepath.Join(uploadDir, filename)
	if err := saveFile(FileHeader, savePath); err != nil {
		return "", err
	}

	return filename, nil
}

func saveFile(FileHeader *multipart.FileHeader, destination string) error {
	// Open the current file
	src, err := FileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create destination file
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	// Move current file -> destination file
	_, err = io.Copy(out, src)

	return err
}
