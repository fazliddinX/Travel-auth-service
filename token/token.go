package token

import (
	"Auth-service/models"
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

const (
	AccessTokenKey  = "key is not easy"
	RefreshTokenKey = "key is really hard"
)

type Claims struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func NewAccessToken(user models.LoginUser) (string, error) {
	claims := Claims{
		Id:    user.Id,
		Name:  user.Name,
		Age:   user.Age,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 3).Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := accessToken.SignedString([]byte(AccessTokenKey))

	if err != nil {
		return "", err
	}

	return token, err
}

func NewRefreshToken(user models.LoginUser) (string, error) {
	claims := Claims{
		Id:    user.Id,
		Name:  user.Name,
		Age:   user.Age,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := refreshToken.SignedString([]byte(RefreshTokenKey))

	return token, err
}

func RenewalAccessToken(claims *Claims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := refreshToken.SignedString([]byte(AccessTokenKey))

	return token, err
}

func ExtractClaimAcces(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(AccessTokenKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func ExtractClaimRefresh(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(RefreshTokenKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
