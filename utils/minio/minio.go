package minio

import (
	"bytes"
	"context"
	"fmt"
	"imageboard/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func UploadImage(imageData []byte, sizeType config.ImageSizeType, fileName string, contentType string) error {
	ctx := context.Background()

	minioClient, err := minio.New(config.S3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.S3.AccessKey, config.S3.SecretAccessKey, ""),
		Secure: config.S3.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize MinIO client: %v", err)
	}

	bucketExists, err := minioClient.BucketExists(ctx, config.S3.BucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %v", err)
	}

	if !bucketExists {
		err = minioClient.MakeBucket(ctx, config.S3.BucketName, minio.MakeBucketOptions{
			Region: config.S3.Region,
		})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}
	}

	objectPath := fmt.Sprintf("%s/%s", sizeType, fileName)

	_, err = minioClient.PutObject(ctx, config.S3.BucketName, objectPath, bytes.NewReader(imageData), int64(len(imageData)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload image: %v", err)
	}

	return nil
}
