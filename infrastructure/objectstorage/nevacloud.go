package objectstorage

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var s3Client *s3.S3

func createSession(access_key string, secret_key string, endpoint string, region string) (*session.Session, error) {
	// Create a new session with the provided credentials and endpoint
	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(access_key, secret_key, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		S3ForcePathStyle: aws.Bool(true), // Force path-style URLs
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	return sess, nil
}

func InitS3(access_key string, secret_key string, endpoint string, bucket string, region string) error {
	sess, err := createSession(access_key, secret_key, endpoint, region)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	s3Client = s3.New(sess)
	// Check if the S3 client is working by listing buckets and print it
	buckets, err := s3Client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("failed to list buckets: %w", err)
	}
	for _, bucket := range buckets.Buckets {
		fmt.Printf("Bucket: %s\n", *bucket.Name)
	}
	return nil
}
func GetS3() *s3.S3 {
	return s3Client
}
