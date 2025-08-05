package configMinio

import (
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

)

var MinioClient *minio.Client

func InitMinio() {
	endpoint := os.Getenv("MINIO_ENDPOINT") + ":" + os.Getenv("MINIO_PORT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := false // Change to true if you're using HTTPS

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("❌ Failed to initialize MinIO: %v", err)
	}

	MinioClient = client
	log.Println("✅ MinIO client initialized successfully")
}
