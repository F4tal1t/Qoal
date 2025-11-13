package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

type S3Storage struct {
	client   *s3.S3
	uploader *s3manager.Uploader
	bucket   string
}

func NewS3Storage(region, bucket, accessKey, secretKey string) (*S3Storage, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return &S3Storage{
		client:   s3.New(sess),
		uploader: s3manager.NewUploader(sess),
		bucket:   bucket,
	}, nil
}

func (s *S3Storage) SaveFile(file io.Reader, filename string, fileSize int64) (string, error) {
	ext := filepath.Ext(filename)
	baseName := strings.TrimSuffix(filename, ext)
	cleanBaseName := strings.ReplaceAll(baseName, " ", "_")
	cleanBaseName = strings.ReplaceAll(cleanBaseName, "..", "_")

	uniqueID := uuid.New().String()
	key := fmt.Sprintf("uploads/%s/%s_%s%s", time.Now().Format("2006/01/02"), cleanBaseName, uniqueID[:8], ext)

	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	_, err = s.uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(GetMimeType(filename)),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	return key, nil
}

func (s *S3Storage) GetFile(filePath string) (io.ReadCloser, error) {
	result, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file from S3: %w", err)
	}

	return result.Body, nil
}

func (s *S3Storage) DeleteFile(filePath string) error {
	_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

func (s *S3Storage) GetOutputPath(jobID string, targetFormat string) string {
	return fmt.Sprintf("processed/%s/%s.%s", time.Now().Format("2006/01/02"), jobID, targetFormat)
}

func (s *S3Storage) SaveProcessedFile(file io.Reader, jobID, targetFormat string) (string, error) {
	key := s.GetOutputPath(jobID, targetFormat)

	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	_, err = s.uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(GetMimeType("." + targetFormat)),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	return key, nil
}

func (s *S3Storage) GetPresignedURL(filePath string, expiration time.Duration) (string, error) {
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
	})

	url, err := req.Presign(expiration)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url, nil
}

func (s *S3Storage) DownloadToFile(ctx context.Context, key string, writer io.WriterAt) error {
	downloader := s3manager.NewDownloaderWithClient(s.client)
	_, err := downloader.DownloadWithContext(ctx, writer, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}
