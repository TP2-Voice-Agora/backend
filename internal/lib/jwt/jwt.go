package jwt

import (
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
