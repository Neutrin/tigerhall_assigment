package model

import "github.com/golang-jwt/jwt"

type Claim struct {
	jwt.StandardClaims
}
