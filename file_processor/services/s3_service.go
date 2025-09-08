package services

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Service struct {
	downloader *s3manager.Downloader
	uploader   *s3manager.Uploader
	bucket     string
}

func NewS3Service(downloader *s3manager.Downloader, uploader *s3manager.Uploader, bucket string) *S3Service {
	return &S3Service{
		downloader: downloader,
		uploader:   uploader,
		bucket:     bucket,
	}
}

func (s *S3Service) DownloadFile(sourcePath, destinationPath string) error {
	file, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", destinationPath, err)
	}
	defer file.Close()

	_, err = s.downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(sourcePath),
	})
	if err != nil {
		return fmt.Errorf("failed to download file from S3: %w", err)
	}

	return nil
}

func (s *S3Service) UploadFile(sourcePath, destinationPath string) error {
	file, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", sourcePath, err)
	}
	defer file.Close()

	_, err = s.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(destinationPath),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return nil
}
