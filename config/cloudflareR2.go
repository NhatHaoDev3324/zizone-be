package config

import (
	"context"
	"os"
	"strings"

	"github.com/NhatHaoDev3324/zizone-be/pkg/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var R2Client *s3.Client

func InitR2() *s3.Client {
	endpoint := strings.TrimSpace(os.Getenv("R2_ENDPOINT"))
	accessKeyID := strings.TrimSpace(os.Getenv("R2_ACCESS_KEY_ID"))
	secretAccessKey := strings.TrimSpace(os.Getenv("R2_SECRET_ACCESS_KEY"))
	bucketName := strings.TrimSpace(os.Getenv("R2_BUCKET_NAME"))

	if endpoint == "" || accessKeyID == "" || secretAccessKey == "" || bucketName == "" {
		log.LogError("Missing required R2 environment variables")
		return nil
	}

	cfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion("auto"),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
	)
	if err != nil {
		log.LogError("Failed to load AWS SDK config: " + err.Error())
		return nil
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
		o.Region = "auto"
	})

	R2Client = client
	log.LogSuccess("Connected to Cloudflare R2 successfully")
	return client
}

func GetR2BucketName() string {
	return strings.TrimSpace(os.Getenv("R2_BUCKET_NAME"))
}

func GetR2PublicURL() string {
	return strings.TrimSpace(os.Getenv("R2_PUBLIC_URL"))
}
