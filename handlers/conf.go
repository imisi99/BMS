package handlers

import (
	"net/http"
	"os"
	"strings"
	"time"

	"BMS/core"
	"BMS/model"
	"BMS/schemas"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	JWTSecret  = os.Getenv("JWTSecret")
	JWTExpiry  = time.Hour * 24
	JWTRExpiry = time.Minute * 5
)

type JWTClaims struct {
	UserID  string
	Purpose string
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashed), err
}

func VerifyPassword(password, hashPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}

func GenerateJWT(userID, purpose string, expiry time.Duration) (string, error) {
	claims := JWTClaims{
		UserID:  userID,
		Purpose: purpose,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}

func VerifyJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &core.CustomError{Code: http.StatusBadRequest, Message: "Invalid token!"}
		}
		return []byte(JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, &core.CustomError{Code: http.StatusBadRequest, Message: "Invalid token!"}
	}

	payload, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, &core.CustomError{Code: http.StatusBadRequest, Message: "Invalid token!"}
	}

	return payload, nil
}

func (s *Service) Authorization(payload schemas.VerifyUser) (string, error) {
	var user model.User
	if err := s.DB.VerifyUser(&user, payload.UserName); err != nil {
		return "", err
	}

	if err := VerifyPassword(payload.Password, user.Password); err != nil {
		return "", &core.CustomError{Code: http.StatusNotFound, Message: "Invalid username or password."}
	}

	token, err := GenerateJWT(user.ID, "login", JWTExpiry)
	if err != nil {
		return "", &core.CustomError{Code: http.StatusInternalServerError, Message: "Failed to generate token."}
	}

	return token, nil
}

func (s *Service) Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Authorization header missing"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Bearer token not in auth header."})
			return
		}

		payload, err := VerifyJWT(parts[1])
		if err != nil {
			cm := err.(*core.CustomError)
			ctx.AbortWithStatusJSON(cm.Code, gin.H{"msg": cm.Message})
			return
		}

		ctx.Set("userID", payload.UserID)
		ctx.Set("purpose", payload.Purpose)

		ctx.Next()
	}
}
