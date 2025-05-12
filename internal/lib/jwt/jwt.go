package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/ictisagora/backend/internal/models"
	"time"
)

// TODO: replace provision of jwt token
func NewToken(user models.User, duration time.Duration, jwtSecret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":   user.UID,
		"email": user.Email,
		"exp":   time.Now().Add(duration).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(jwtSecret))

	return tokenString
}

func ParseToken(tokenString string, jwtSecret string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid signing method")
		}

		return []byte(jwtSecret), nil
	})

	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", "", fmt.Errorf("Invalid token or claims")
	}

	uid, ok := claims["uid"].(string)
	if !ok {
		return "", "", fmt.Errorf("Invalid uid")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", "", fmt.Errorf("Invalid email")
	}

	return uid, email, nil
}
