package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/achmadnr21/emploman/config"
	"github.com/achmadnr21/emploman/infrastructure/objectstorage"
	pgsql "github.com/achmadnr21/emploman/infrastructure/rdbms"
	"github.com/achmadnr21/emploman/internal/utils"
)

func main() {
	fmt.Println("========> Welcome to Emploman API")

	sc, err := serviceInit()
	if err != nil {
		fmt.Println("[Error] initializing service:", err)
		panic("Service initialization failed")
	}

	s3client := objectstorage.GetS3()
	if s3client == nil {
		panic("[Error] S3 client is nil")
	}

	// === File to upload
	var objectKey string = "defaultprofile.jpg"
	var filePath string = fmt.Sprintf("cmd/s3test/%s", objectKey)
	var folderPath string = "pictureprofile"
	fullObjectKey := fmt.Sprintf("%s/%s", folderPath, objectKey)

	// === Baca file
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("[Error] reading file: %v", err))
	}

	// === Deteksi konten
	contentType := http.DetectContentType(fileBytes)
	if contentType != "image/jpeg" && contentType != "image/png" {
		panic(fmt.Sprintf("[Error] Unsupported file type: %s", contentType))
	}

	// === Upload ke S3
	_, err = s3client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(sc.S3bucket),
		Key:         aws.String(fullObjectKey),
		Body:        bytes.NewReader(fileBytes),
		ACL:         aws.String("public-read"),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		panic(fmt.Sprintf("[Error] uploading file to S3: %v", err))
	}
	fmt.Println("[Success] file uploaded to S3:", fullObjectKey)

	// // === List objek di bucket
	// listAndPrintObjects(s3client, sc.S3bucket)

	// // === Delete dan tunggu sampai object tidak ada
	// deleteAndWait(s3client, sc.S3bucket, fullObjectKey)

	// // === Cek akhir
	listAndPrintObjects(s3client, sc.S3bucket)
	headCheck(s3client, sc.S3bucket, fullObjectKey)
}

// ======================= FUNCTION ==========================

func deleteAndWait(s3client *s3.S3, bucket, key string) {
	_, err := s3client.DeleteObject(&s3.DeleteObjectInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(key),
		VersionId: nil,
	})

	// _, err := s3client.DeleteObjects(&s3.DeleteObjectsInput{
	// 	Bucket: aws.String(bucket),
	// 	Delete: &s3.Delete{
	// 		Objects: []*s3.ObjectIdentifier{
	// 			{
	// 				Key: aws.String(key),
	// 			},
	// 		},
	// 		Quiet: aws.Bool(false),
	// 	},
	// })
	if err != nil {
		panic(fmt.Sprintf("[Error] deleting object: %v", err))
	}
	fmt.Println("[Info] delete request sent, waiting for confirmation...")

	// // WaitUntilObjectNotExists
	// err = s3client.WaitUntilObjectNotExistsWithContext(context.Background(), &s3.HeadObjectInput{
	// 	Bucket: aws.String(bucket),
	// 	Key:    aws.String(key),
	// })
	// if err != nil {
	// 	panic(fmt.Sprintf("[Error] object still exists after delete: %v", err))
	// }
	// fmt.Println("[Success] object confirmed deleted:", key)
}

func listAndPrintObjects(s3client *s3.S3, bucket string) {
	result, err := s3client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		panic(fmt.Sprintf("[Error] listing objects: %v", err))
	}
	fmt.Println("[Success] listing objects in bucket:")
	for _, item := range result.Contents {
		fmt.Printf(" - %s : %d bytes\n", *item.Key, *item.Size)
	}
}

func headCheck(s3client *s3.S3, bucket, key string) {
	_, err := s3client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		fmt.Println("[Success] object sudah tidak ada di bucket.")
	} else {
		fmt.Println("[Warning] object masih terdeteksi di bucket.")
	}
}

func serviceInit() (config.Config, error) {
	var envload config.Config
	envload.LoadConfig()

	if err := pgsql.InitPG(envload.DbHost, int32(envload.DbPort), envload.DbUser, envload.DbPassword, envload.DbName, envload.DbSsl); err != nil {
		return config.Config{}, fmt.Errorf("[Error] PostgreSQL init: %v", err)
	}

	if err := objectstorage.InitS3(envload.S3accesskey, envload.S3secretkey, envload.S3endpoint, envload.S3bucket, envload.S3region); err != nil {
		return config.Config{}, fmt.Errorf("[Error] S3 init: %v", err)
	}

	utils.JwtInit(envload.JwtSecret, envload.RefreshSecret)
	return envload, nil
}
