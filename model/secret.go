package model

import "github.com/dgrijalva/jwt-go"

type (
	Auth struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	Secret struct {
		App  string            `json:"app"`
		Vars map[string]string `json:"vars"`
	}

	CustomClaims struct {
		*jwt.StandardClaims
		Username string
	}
)
