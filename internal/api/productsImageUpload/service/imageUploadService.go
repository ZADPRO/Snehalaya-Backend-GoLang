package imageUploadService

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"time"

	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/productsImageUpload/config"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func init() {
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"
	endpoint := os.Getenv("MINIO_ENDPOINT") + ":" + "443"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ACCESS_KEY"), os.Getenv("MINIO_SECRET_KEY"), ""),
		Secure: useSSL,
	})
	if err != nil {
		panic(err)
	}

	MinioClient = client
}

func CreateUploadURL(fileName string, expireMins int) (string, string, error) {
	log.Printf("Creating presigned PUT URL | fileName: %s | expireMins: %d", fileName, expireMins)

	bucket := os.Getenv("MINIO_BUCKET")
	expiry := time.Duration(expireMins) * time.Minute

	uploadURL, err := MinioClient.PresignedPutObject(context.Background(), bucket, fileName, expiry)
	if err != nil {
		log.Printf("Failed to generate presigned PUT URL | Error: %v", err)
		return "", "", err
	}

	fileURL, err := GetFileURL(fileName, expireMins)
	if err != nil {
		log.Printf("Failed to generate file GET URL after PUT | Error: %v", err)
		return "", "", err
	}

	return uploadURL.String(), fileURL, nil
}

func GetFileURL(fileName string, expireMins int) (string, error) {
	log.Printf("Generating presigned GET URL | fileName: %s | expireMins: %d", fileName, expireMins)

	bucket := os.Getenv("MINIO_BUCKET")
	expiry := time.Duration(expireMins) * time.Minute
	reqParams := url.Values{}

	fileURL, err := MinioClient.PresignedGetObject(context.Background(), bucket, fileName, expiry, reqParams)
	if err != nil {
		log.Printf("Error generating presigned GET URL | Error: %v", err)
		return "", err
	}

	return fileURL.String(), nil
}

func FetchAllEnvVariables() map[string]string {
	log.Println("Reading .env variables from file")
	envMap, err := godotenv.Read(".env")
	if err != nil {
		log.Printf("Error reading .env file: %v", err)
		return map[string]string{
			"error": "Failed to load .env",
		}
	}

	log.Printf("Loaded %d environment variables", len(envMap))
	return envMap
}

func GeneratePresignedURL(extension string) (string, string, error) {
	if config.MinioClient == nil {
		log.Println("❌ MinIO client is nil")
		return "", "", fmt.Errorf("MinIO client not initialized")
	}

	bucket := "zadroit-dev"

	timestamp := time.Now().Unix()
	randomPart := rand.Intn(10000)
	filename := fmt.Sprintf("IMG-%d-%d.%s", timestamp, randomPart, extension)
	objectName := "uploads/" + filename

	log.Println("🔄 Generating pre-signed URL for:", objectName)

	presignedURL, err := config.MinioClient.PresignedPutObject(
		context.Background(),
		bucket,
		objectName,
		15*time.Minute,
	)
	if err != nil {
		log.Printf("❌ Failed to generate pre-signed URL: %v\n", err)
		return "", "", err
	}

	return presignedURL.String(), filename, nil
}
