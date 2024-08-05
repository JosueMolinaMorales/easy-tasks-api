package config

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JosueMolinaMorales/EasyTasksAPI/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func NewClaims(userID string) *Claims {
	return &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "EasyTasksAPI",
		},
	}
}

func (c *Claims) GetExpirationTime() (*jwt.NumericDate, error) {
	return c.ExpiresAt, nil
}

func (c *Claims) GetIssuedAt() (*jwt.NumericDate, error) {
	return c.IssuedAt, nil
}

func (c *Claims) GetNotBefore() (*jwt.NumericDate, error) {
	return c.NotBefore, nil
}

func (c *Claims) GetIssuer() (string, error) {
	return c.Issuer, nil
}

func (c *Claims) GetSubject() (string, error) {
	return c.Subject, nil
}

func (c *Claims) GetAudience() (jwt.ClaimStrings, error) {
	return c.Audience, nil
}

func (c *Claims) SignToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	signed, err := token.SignedString([]byte("secret")) // TODO: Replace secret
	if err != nil {
		return "", err
	}

	return signed, nil
}

func VerifyToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		fmt.Println("could not parse claims")
	}

	return claims, nil
}

func ExtractIDFromToken(ctx *gin.Context) (string, *errors.RequestError) {
	tokenAny, ok := ctx.Get("Token")
	if !ok {
		return "", errors.NewRequestError(http.StatusInternalServerError, "No token found")
	}
	token := tokenAny.(*Claims)
	author := token.UserID

	return author, nil
}
