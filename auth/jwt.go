package auth

import (
	"strings"
	"time"

	"github.com/NeowayLabs/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/ffhenkes/kripto/fs"
	"github.com/ffhenkes/kripto/model"
)

var logJ = logger.Namespace("kripto.jwt")

const (
	dataRsa    = "/data/rsa"
	keyName    = "kripto"
	signMethod = "RS256"
)

type (
	// JwtAuth represent the json authorization token type
	JwtAuth struct {
		c *model.Credentials
	}
)

// NewJwtAuth returns a reference for JwtAuth type and embed Credentials
func NewJwtAuth(c *model.Credentials) *JwtAuth {
	return &JwtAuth{c}
}

// GenerateToken uses Credentials to generate a new Jwt (json authorization token)
func (jwta *JwtAuth) GenerateToken() (string, error) {

	sys := fs.NewFileSystem(dataRsa)

	privateKey, err := sys.ReadKey(keyName)
	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.GetSigningMethod(signMethod))
	token.Claims = &model.CustomClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwta.c.TokenExpiresIn).Unix(),
		},
		jwta.c.Username,
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", nil
	}

	return tokenString, nil
}

// ValidateToken checks the token integrity
func ValidateToken(authorization string) (bool, error) {

	sys := fs.NewFileSystem(dataRsa)

	pub, err := sys.ReadPublicKey(keyName)
	if err != nil {
		return false, err
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		return false, err
	}

	tokenString := strings.TrimSpace(authorization)
	token, err := jwt.ParseWithClaims(tokenString, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	expiresAt := time.Unix(token.Claims.(*model.CustomClaims).StandardClaims.ExpiresAt, 0)
	username := token.Claims.(*model.CustomClaims).Username

	if err != nil {
		logJ.Warn("User: %s Non valid token! Expires At: %v", username, expiresAt)
		return token.Valid, err
	}

	logJ.Info("User: %s Token Expires At: %v", username, expiresAt)
	return token.Valid, nil
}
