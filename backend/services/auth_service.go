package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/qoal/file-processor/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (as *AuthService) Register(user *models.User) error {
	// Check if user already exists
	var existingUser models.User
	err := as.db.Table("qoal_user").Where("email = ?", user.Email).First(&existingUser).Error
	if err == nil {
		return errors.New("user already exists")
	}
	if err != gorm.ErrRecordNotFound {
		fmt.Printf("User existence check error: %v\n", err)
		return err // Return other database errors
	}
	fmt.Printf("User existence check passed, no existing user found\n")

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Create user - omit ID to let PostgreSQL generate UUID automatically
	if err := as.db.Table("qoal_user").Select("email", "password", "name", "created_at", "updated_at").Create(user).Error; err != nil {
		return err
	}

	// Reload the user to get the generated UUID
	var createdUser models.User
	fmt.Printf("Attempting to reload user with email: %s\n", user.Email)
	if err := as.db.Table("qoal_user").Where("email = ?", user.Email).First(&createdUser).Error; err != nil {
		fmt.Printf("Failed to reload user: %v\n", err)
		return err
	}
	fmt.Printf("Successfully reloaded user with ID: %s\n", createdUser.ID)

	// Update the original user struct with the created data
	user.ID = createdUser.ID
	user.CreatedAt = createdUser.CreatedAt
	user.UpdatedAt = createdUser.UpdatedAt
	user.Password = "" // Clear password for security

	return nil
}

func (as *AuthService) Login(email, password string) (*models.User, error) {
	var user models.User
	if err := as.db.Table("qoal_user").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

func (as *AuthService) GenerateToken(user *models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not configured")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(secret))
}

func (as *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET not configured")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"].(string)
		var user models.User
		if err := as.db.Table("qoal_user").Where("id = ?", userID).First(&user).Error; err != nil {
			return nil, errors.New("user not found")
		}

		return &user, nil
	}

	return nil, errors.New("invalid token")
}
