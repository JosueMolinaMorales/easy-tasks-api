package config

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    string           `json:"userId"`
	ExpiresAt *jwt.NumericDate `json:"exp"`
	IssuedAt  *jwt.NumericDate `json:"iat"`
	NotBefore *jwt.NumericDate `json:"nbf,omitempty"`
	Issuer    string           `json:"iss"`
	Subject   string           `json:"sub,omitempty"`
	Audience  jwt.ClaimStrings `json:"aud,omitempty"`
}

func NewClaims(userID string) *Claims {
	return &Claims{
		UserID:    userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "EasyTasksAPI",
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
