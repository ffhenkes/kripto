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
	data_rsa         = "/data/rsa"
	private_key_name = "kripto.rsa"
	public_key_name  = "kripto.rsa.pub"
	sign_method      = "RS256"
	time_frame       = 24
)

type (
	JwtAuth struct {
		c *model.Credentials
	}
)

func NewJwtAuth(c *model.Credentials) *JwtAuth {
	return &JwtAuth{c}
}

func (jwta *JwtAuth) GenerateToken() (string, error) {

	sys := fs.NewFileSystem(data_rsa)

	privateKey, err := sys.ReadKey(private_key_name)
	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.GetSigningMethod(sign_method))
	token.Claims = &model.CustomClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time_frame).Unix(),
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

func ValidateToken(authorization string) (bool, error) {

	sys := fs.NewFileSystem(data_rsa)

	pub, err := sys.ReadKey(public_key_name)
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

	expires_at := time.Unix(token.Claims.(*model.CustomClaims).StandardClaims.ExpiresAt, 0)
	username := token.Claims.(*model.CustomClaims).Username

	if err != nil {
		logJ.Warn("User: %s Non valid token! Expires At: %v", username, expires_at)
		return token.Valid, err
	}

	logJ.Info("User: %s Token Expires At: %v", username, expires_at)
	return token.Valid, nil
}
