package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/nitin/tigerhall/core/internal/config"
	"github.com/nitin/tigerhall/core/internal/model"
)

func GenerateSignedTokens(uniqueKey string, expireTimeMin int) (string, error) {
	//expiry time of token
	expirationTime := time.Now().Add(time.Duration(expireTimeMin) * time.Minute)

	claims := &model.Claim{
		StandardClaims: jwt.StandardClaims{
			Subject:   uniqueKey,
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.JwtKey))
}

func ParseToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Claim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JwtKey), nil
	})

	if err != nil {
		return "", err
	}

	claim, ok := token.Claims.(*model.Claim)

	if !ok {
		return "", err
	}

	return claim.StandardClaims.Subject, nil
}
