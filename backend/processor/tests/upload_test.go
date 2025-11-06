package tests

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/qoal/file-processor/handlers"
	"github.com/stretchr/testify/assert"
)

// 🎓 **Go Tutorial: Testing in Go**
// Go has built-in testing support with the testing package
func TestUploadHandler(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test upload handler
	uploadHandler := handlers.NewUploadHandler()

	// 🎓 **Go Tutorial: Table-Driven Tests**
	// This is a common pattern in Go for testing multiple scenarios
	tests := []struct {
		name           string
		fileContent    []byte
		fileName       string
		expectedStatus int
		expectedError  bool
	}{
		{
			name:           "Valid PDF upload",
			fileContent:    []byte("test pdf content"),
			fileName:       "test.pdf",
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "Valid image upload",
			fileContent:    []byte("fake image content"),
			fileName:       "test.jpg",
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "Invalid file type",
			fileContent:    []byte("fake exe content"),
			fileName:       "test.exe",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a buffer to write our multipart form
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			// 🎓 **Go Tutorial: Working with Multipart Forms**
			// Create form file
			part, err := writer.CreateFormFile("file", tt.fileName)
			assert.NoError(t, err)

			// Write file content
			_, err = part.Write(tt.fileContent)
			assert.NoError(t, err)

			// Close the writer to finalize the form
			err = writer.Close()
			assert.NoError(t, err)

			// Create test request
			req := httptest.NewRequest("POST", "/api/v1/upload", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			// Create response recorder
			w := httptest.NewRecorder()

			// Create Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the handler
			uploadHandler.HandleUpload(c)

			// Check response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// 🎓 **Go Tutorial: Benchmarking**
// Benchmark tests measure performance
func BenchmarkUploadHandler(b *testing.B) {
	gin.SetMode(gin.TestMode)
	uploadHandler := handlers.NewUploadHandler()

	// Create test data
	fileContent := bytes.Repeat([]byte("test content"), 1000) // 12KB file

	b.ResetTimer() // Start timing

	for i := 0; i < b.N; i++ {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, _ := writer.CreateFormFile("file", "benchmark.pdf")
		part.Write(fileContent)
		writer.Close()

		req := httptest.NewRequest("POST", "/api/v1/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		uploadHandler.HandleUpload(c)
	}
}

// 🎓 **Go Tutorial: Cleanup Tests**
// TestMain can be used for setup/teardown
func TestMain(m *testing.M) {
	// Setup: Create test upload directory
	testDir := "test_uploads"
	os.MkdirAll(testDir, 0755)

	// Run tests
	code := m.Run()

	// Cleanup: Remove test directory
	os.RemoveAll(testDir)

	os.Exit(code)
}
