package model

import "github.com/dgrijalva/jwt-go"

type (
	Credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	CustomClaims struct {
		*jwt.StandardClaims
		Username string
	}
)
