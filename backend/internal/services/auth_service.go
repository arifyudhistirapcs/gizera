package services

import (
	"errors"
	"time"

	"github.com/erp-sppg/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("kredensial tidak valid")
	ErrUserNotFound       = errors.New("pengguna tidak ditemukan")
	ErrUserInactive       = errors.New("akun tidak aktif")
	ErrInvalidToken       = errors.New("token tidak valid")
)

// JWTClaims represents the claims in the JWT token
type JWTClaims struct {
	UserID    uint   `json:"user_id"`
	Role      string `json:"role"`
	SPPGID    *uint  `json:"sppg_id,omitempty"`
	YayasanID *uint  `json:"yayasan_id,omitempty"`
	jwt.RegisteredClaims
}

// AuthService handles authentication operations
type AuthService struct {
	db        *gorm.DB
	jwtSecret []byte
}

// NewAuthService creates a new authentication service
func NewAuthService(db *gorm.DB, jwtSecret string) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: []byte(jwtSecret),
	}
}

// Login authenticates a user with NIK/Email and password
func (s *AuthService) Login(identifier, password string) (*models.User, string, error) {
	var user models.User

	// Find user by NIK or Email, preload SPPG and Yayasan relations
	result := s.db.Preload("SPPG").Preload("Yayasan").Where("nik = ? OR email = ?", identifier, identifier).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", result.Error
	}

	// Check if user is active
	if !user.IsActive {
		return nil, "", ErrUserInactive
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Generate JWT token with tenant info
	token, err := s.GenerateToken(user.ID, user.Role, user.SPPGID, user.YayasanID)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

// GenerateToken creates a new JWT token for a user
func (s *AuthService) GenerateToken(userID uint, role string, sppgID *uint, yayasanID *uint) (string, error) {
	claims := JWTClaims{
		UserID:    userID,
		Role:      role,
		SPPGID:    sppgID,
		YayasanID: yayasanID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// RefreshToken generates a new token for an existing valid token
func (s *AuthService) RefreshToken(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Verify user still exists and is active
	var user models.User
	result := s.db.First(&user, claims.UserID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", ErrUserNotFound
		}
		return "", result.Error
	}

	if !user.IsActive {
		return "", ErrUserInactive
	}

	// Generate new token preserving tenant claims
	return s.GenerateToken(user.ID, user.Role, claims.SPPGID, claims.YayasanID)
}

// HashPassword hashes a password using bcrypt
func (s *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	result := s.db.Preload("SPPG").Preload("Yayasan").First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}
