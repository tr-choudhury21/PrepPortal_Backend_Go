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
		return "", errors.New("Cloudinary is not initialized")
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
