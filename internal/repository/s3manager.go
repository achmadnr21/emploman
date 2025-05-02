package repository

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Repository struct {
	s3Client *s3.S3
	bucket   string
}

func NewS3Repository(s3Client *s3.S3, bucket string) *S3Repository {
	return &S3Repository{
		s3Client: s3Client,
		bucket:   bucket,
	}
}
func (r *S3Repository) UploadFile(key string, fileBytes []byte, contentType string) (string, error) {
	// Upload the file to S3
	_, err := r.s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(fileBytes),
		ACL:         aws.String("public-read"),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return key, nil
}
func (r *S3Repository) GetFileURL(key string) (string, error) {
	// Generate a pre-signed URL for the file
	req, _ := r.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed to generate pre-signed URL: %w", err)
	}

	return urlStr, nil
}
func (r *S3Repository) DeleteFile(key string) error {
	// Delete the file from S3
	_, err := r.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}
func (r *S3Repository) GetFile(key string) ([]byte, error) {
	// Get the file from S3
	result, err := r.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file from S3: %w", err)
	}
	defer result.Body.Close()

	// Read the file content
	fileBytes, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	return fileBytes, nil
}
func (r *S3Repository) GetFileMetadata(key string) (string, error) {
	// Get the file metadata from S3
	result, err := r.s3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get file metadata from S3: %w", err)
	}

	return *result.ContentType, nil
}
