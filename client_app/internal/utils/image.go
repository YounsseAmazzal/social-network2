package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

const (
	MaxUploadSize          = 20 << 20
	MaxPixelSize           = 8000
	UploadDirectorypost    = "static/socialup/posts"
	UploadDirectorycomment = "static/socialup/comment"
	UploadDirectoryauth    = "static/socialup/auth"

	FilePermissions = 0644
	DirPermissions  = 0755
)

var supportedExtaaions = map[string]string{
	"jpeg": ".jpg",
	"png":  ".png",
	"gif":  ".gif",
}

// check the safty of the image
func ValidateImageContent(imgData []byte) (string, error) {
	if len(imgData) == 0 {
		return "", errors.New("image data is empty")
	}

	config, format, err := image.DecodeConfig(bytes.NewReader(imgData))
	if err != nil {
		return "", fmt.Errorf("invalid image format: %w", err)
	}

	if config.Width > MaxPixelSize || config.Height > MaxPixelSize {
		return "", fmt.Errorf("dimensions %dx%d exceed limit of %dpx", config.Width, config.Height, MaxPixelSize)
	}

	ext, isSupported := supportedExtaaions[format]
	//fmt.Println("Imaaage format:",format)
	if !isSupported {
		return "", fmt.Errorf("unsupported format: %s (only PNG, JPEG, GIF allowed)", format)
	}

	//(Security Check) all the imaaage
	if _, _, err := image.Decode(bytes.NewReader(imgData)); err != nil {
		return "", fmt.Errorf("corrupted image data: %w", err)
	}

	return ext, nil
}

// save image with hash ,uploead and save the immage
func ProcessImageUpload(file multipart.File, header *multipart.FileHeader, authorpostsorcomment string) (string, error) {
	if file == nil || header == nil {
		return "", errors.New("no file was uploaded")
	}
	//preee cheeck
	if header.Size > MaxUploadSize {
		// //fmt.Println("IIIIIIIIIIIIMMMMMAAAAAAAAAAGEEE",header.Size)
		return "", fmt.Errorf("file size (%d) exceeds limit of 20MB", header.Size)
	}
	//reaaad on memoery
	fileBuffer, err := io.ReadAll(io.LimitReader(file, MaxUploadSize+1))
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	if int64(len(fileBuffer)) > MaxUploadSize {
		return "", errors.New("file size actual content exceeds 20MB")
	}
	fileExtension, err := ValidateImageContent(fileBuffer)
	if err != nil {
		return "", err
	}

	checksum := sha256.Sum256(fileBuffer)
	uniqueName := hex.EncodeToString(checksum[:]) + fileExtension
	// log.Fatal("authorposts", authorpostsorcomment)
	if authorpostsorcomment == "post" {
		if err := os.MkdirAll(UploadDirectorypost, DirPermissions); err != nil {
			return "", fmt.Errorf("failed to create posst directory: %w", err)
		}
		destinationPath := filepath.Join(UploadDirectorypost, uniqueName)
			if err := os.WriteFile(destinationPath, fileBuffer, FilePermissions); err != nil {
		return "", fmt.Errorf("failed to write file to disk: %w", err)
		}
	}
	if authorpostsorcomment == "auth" {
		if err := os.MkdirAll(UploadDirectoryauth, DirPermissions); err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
		destinationPath := filepath.Join(UploadDirectoryauth, uniqueName)
			if err := os.WriteFile(destinationPath, fileBuffer, FilePermissions); err != nil {
		return "", fmt.Errorf("failed to write file to disk: %w", err)
			}
	}
	if authorpostsorcomment == "comment" {
		if err := os.MkdirAll(UploadDirectorycomment, DirPermissions); err != nil {
			return "", fmt.Errorf("failed to create comment directory: %w", err)
		}
		destinationPath := filepath.Join(UploadDirectorycomment, uniqueName)
	if err := os.WriteFile(destinationPath, fileBuffer, FilePermissions); err != nil {
		return "", fmt.Errorf("failed to write file to disk: %w", err)
	}
	}

	// if err := os.MkdirAll(UploadDirectorypost, DirPermissions); err != nil {
	// 	return "", fmt.Errorf("failed to create directory: %w", err)
	// }

	// destinationPath := filepath.Join(UploadDirectorypost, uniqueName)
	// if err := os.WriteFile(destinationPath, fileBuffer, FilePermissions); err != nil {
	// 	return "", fmt.Errorf("failed to write file to disk: %w", err)
	// }

	return uniqueName, nil
}
