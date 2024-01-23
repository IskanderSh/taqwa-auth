package jwt

import (
	"time"

	"github.com/IskanderSh/taqwa-auth/internal/config"
	"github.com/IskanderSh/taqwa-auth/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user *models.User, configToken *config.Token) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(configToken.TTL).Unix()

	tokenString, err := token.SignedString(configToken.Secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
