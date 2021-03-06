package model

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type (
	// Credentials represents the authentication model containing username and password
	// This model will be embed into the Login and Jwt types
	Credentials struct {
		Username       string        `json:"username"`
		Password       string        `json:"password"`
		TokenExpiresIn time.Duration `json:"token_expires_in"`
	}

	// CustomClaims represents the customizable model embed into the Jwt
	// This one is used to carry on custom data within the token
	CustomClaims struct {
		*jwt.StandardClaims
		Username string
	}
)
