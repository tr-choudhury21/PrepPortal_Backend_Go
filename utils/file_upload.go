package utils

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/tr-choudhury21/prepportal_backend/config"
)

func UploadFile(file multipart.File, fileName string) (string, error) {

	if config.CLD == nil {
		return "", errors.New("cloudinary is not initialized")
	}

	ctx := context.Background()

	// Upload the file to Cloudinary
	uploadResult, err := config.CLD.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: fileName,
		Folder:   "documents",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	return uploadResult.SecureURL, nil
}

// UploadImage uploads an image file to Cloudinary and returns the URL
func UploadImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {

	if config.CLD == nil {
		return "", errors.New("cloudinary is not initialized")
	}

	ctx := context.Background()

	// Upload to Cloudinary
	resp, err := config.CLD.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "blog_images", // Store images in a separate folder
	})
	if err != nil {
		return "", fmt.Errorf("error uploading image: %v", err)
	}

	return resp.SecureURL, nil
}
