package bulkImageUploadService

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	logger "github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/helper/Logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func init() {
	log := logger.InitLogger()
	useSSL := true

	endpoint := "test-zad.brightoncloudtech.com:443"
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Errorf("❌ Failed to initialize MinIO client: %v", err)
	}

	MinioClient = client
	log.Infof("✅ MinIO client initialized for Bulk Image Upload")
}

func CreatePresignedURLs(fileName string, expireMins int) (string, string, error) {
	bucket := os.Getenv("MINIO_BUCKET")
	expiry := time.Duration(expireMins) * time.Minute

	objectName := "bulk-images/" + strings.ToUpper(fileName)

	uploadURL, err := MinioClient.PresignedPutObject(context.Background(), bucket, objectName, expiry)
	if err != nil {
		log.Printf("Failed to create presigned PUT URL for %s: %v", objectName, err)
		return "", "", err
	}

	viewURL, err := GetImageViewURL(objectName, expireMins)
	if err != nil {
		log.Printf("Failed to create presigned GET URL for %s: %v", objectName, err)
		return "", "", err
	}

	return uploadURL.String(), viewURL, nil
}

func GetImageViewURL(fileName string, expireMins int) (string, error) {
	bucket := os.Getenv("MINIO_BUCKET")
	expiry := time.Duration(expireMins) * time.Minute
	reqParams := url.Values{}

	url, err := MinioClient.PresignedGetObject(context.Background(), bucket, fileName, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("Failed to generate view URL: %w", err)
	}
	return url.String(), nil
}
