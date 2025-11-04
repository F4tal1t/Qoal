package tests

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/qoal/file-processor/services"
	"github.com/stretchr/testify/assert"
)

func TestS3Operations(t *testing.T) {
	// Load .env file
	loadEnvFile()

	// Skip if no AWS credentials
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Skip("AWS credentials not configured, skipping S3 tests")
	}

	// Real S3 service
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))
	downloader := s3manager.NewDownloader(sess)
	uploader := s3manager.NewUploader(sess)
	s3Service := services.NewS3Service(downloader, uploader, os.Getenv("AWS_S3_BUCKET"))

	t.Run("TestS3Upload", func(t *testing.T) {
		// Create test file
		testFile := "test_upload.txt"
		testContent := "Hello S3!"
		err := ioutil.WriteFile(testFile, []byte(testContent), 0644)
		assert.NoError(t, err)
		defer os.Remove(testFile)

		// Upload to S3
		err = s3Service.UploadFile(testFile, "test/upload.txt")
		assert.NoError(t, err)
	})

	t.Run("TestS3Download", func(t *testing.T) {
		// Download from S3
		downloadFile := "test_download.txt"
		err := s3Service.DownloadFile("test/upload.txt", downloadFile)
		assert.NoError(t, err)
		defer os.Remove(downloadFile)

		// Verify content
		content, err := ioutil.ReadFile(downloadFile)
		assert.NoError(t, err)
		assert.Equal(t, "Hello S3!", string(content))
	})
}

func loadEnvFile() {
	// Try multiple paths
	paths := []string{".env", "../.env", "../../.env"}

	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				os.Setenv(parts[0], parts[1])
			}
		}
		return // Found and loaded .env file
	}
}
