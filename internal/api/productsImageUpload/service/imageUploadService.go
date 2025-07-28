package imageUploadService

import (
	"context"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func init() {
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"
	port := os.Getenv("MINIO_PORT")
	endpoint := os.Getenv("MINIO_ENDPOINT") + ":" + port

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
	bucket := os.Getenv("MINIO_BUCKET")
	expiry := time.Duration(expireMins) * time.Minute

	uploadURL, err := MinioClient.PresignedPutObject(context.Background(), bucket, fileName, expiry)
	if err != nil {
		return "", "", err
	}

	fileURL, err := GetFileURL(fileName, expireMins)
	if err != nil {
		return "", "", err
	}

	return uploadURL.String(), fileURL, nil
}

func GetFileURL(fileName string, expireMins int) (string, error) {
	bucket := os.Getenv("MINIO_BUCKET")
	expiry := time.Duration(expireMins) * time.Minute

	reqParams := url.Values{}

	fileURL, err := MinioClient.PresignedGetObject(context.Background(), bucket, fileName, expiry, reqParams)
	if err != nil {
		return "", err
	}

	return fileURL.String(), nil
}

func FetchAllEnvVariables() map[string]string {
	envMap, err := godotenv.Read(".env")
	if err != nil {
		log.Println("Error reading .env file:", err)
		return map[string]string{
			"error": "Failed to load .env",
		}
	}
	return envMap
}
