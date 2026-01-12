package aws

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type UploadService struct {
	S3Client *s3.Client
	bucket   string
	region   string
}

func NewUploadService(s3Client *s3.Client) *UploadService {
	return &UploadService{
		S3Client: s3Client,
		bucket:   os.Getenv("AWS_S3_BUCKET"),
		region:   os.Getenv("AWS_REGION"),
	}
}

func (u *UploadService) UploadImage(ctx context.Context, file multipart.File, header *multipart.FileHeader, bucketName string) (*UploadResponse, error) {
	if !isValidFileType(file, header) {
		return nil, errors.New("only image (PNG, JPEG, WebP) and PDF files are allowed")
	}

	if _, err := file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to reset file pointer: %v", err)
	}

	key := fmt.Sprintf(
		"%s/%d%s",
		bucketName,
		time.Now().UnixNano(),
		filepath.Ext(header.Filename),
	)

	contentType := header.Header.Get("Content-Type")

	_, err := u.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &u.bucket,
		Key:         &key,
		Body:        file,
		ContentType: &contentType,
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to upload Image: %v", err)
	}

	imageURL := fmt.Sprintf(
		"https://%s.s3.%s.amazonaws.com/%s",
		u.bucket,
		u.region,
		key,
	)

	response := &UploadResponse{
		URL: imageURL,
	}
	return response, nil
}

func isValidFileType(file multipart.File, header *multipart.FileHeader) bool {
	clientType := header.Header.Get("Content-Type")

	buff := make([]byte, 512)
	_, err := file.Read(buff)
	if err != nil && err != io.EOF {
		return false
	}
	detectedType := http.DetectContentType(buff)

	// Define allowed types
	allowed := map[string]bool{
		"image/png":       true,
		"image/jpeg":      true,
		"image/webp":      true,
		"application/pdf": true, // PDF Support added
	}

	return allowed[clientType] || allowed[detectedType]
}
