package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/qoal/file-processor/handlers"
	"github.com/qoal/file-processor/models"
	"github.com/qoal/file-processor/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAuthTest() (*gin.Engine, *services.AuthService, func()) {
	// Setup test database
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{})

	// Setup services
	authService := services.NewAuthService(db)
	authHandler := handlers.NewAuthHandler(authService)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Public routes
	router.POST("/api/v1/auth/register", authHandler.Register)
	router.POST("/api/v1/auth/login", authHandler.Login)

	return router, authService, func() {
		// Cleanup
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
}

func TestRegister(t *testing.T) {
	router, _, cleanup := setupAuthTest()
	defer cleanup()

	// Test registration
	payload := models.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	jsonPayload, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, "test@example.com", response.User.Email)
}

func TestLogin(t *testing.T) {
	router, _, cleanup := setupAuthTest()
	defer cleanup()

	// First register a user
	registerPayload := models.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	jsonRegister, _ := json.Marshal(registerPayload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonRegister))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Now test login
	loginPayload := models.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonLogin, _ := json.Marshal(loginPayload)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonLogin))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, "test@example.com", response.User.Email)
}

func TestLoginWithInvalidCredentials(t *testing.T) {
	router, _, cleanup := setupAuthTest()
	defer cleanup()

	loginPayload := models.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "wrongpassword",
	}
	jsonLogin, _ := json.Marshal(loginPayload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonLogin))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}