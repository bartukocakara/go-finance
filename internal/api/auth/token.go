package auth

import (
	"errors"
	"financial-app/internal/model"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key") // Todo change the key

var accessTokenDuration = time.Duration(30) * time.Minute

// var accessTokenDuration = time.Duration(60) * time.Second // TEST: we will get this token, try to get user, first time we should get one. we will wait for 1 miniute and try again. it should fail.
var refreshTokenDuration = time.Duration(30*24) * time.Hour

// var refreshTokenDuration = time.Duration(120) * time.Second // TEST: after access token fails we send this token and it should return new one. we will wait 1 miniute and we should get expried error

type Claims struct {
	UserID model.UserID `json:"userID"`
	jwt.StandardClaims
}

type Token struct {
	AccessToken           string `json:"accessToken,omitempty"`
	AccessTokenExpiresAt  int64  `json:"expiresAt,omitempty"`
	RefreshToken          string `json:"refreshToken,omitempty"`
	RefreshTokenExpiresAt int64  `json:"-"` // We will store this time in database with refresh token
}

func IssueToken(principal model.Principal) (*Token, error) {
	if principal.UserID == model.NilUserID {
		return nil, errors.New("Invalid principal")
	}

	accessToken, accessTokenExpiresAt, err := generateToken(principal, accessTokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshTokenExpiresAt, err := generateToken(principal, refreshTokenDuration)
	if err != nil {
		return nil, err
	}
	token := Token{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}
	return &token, nil
}

func generateToken(principal model.Principal, duration time.Duration) (string, int64, error) {
	now := time.Now()
	claims := &Claims{
		UserID: principal.UserID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(duration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", 0, err
	}

	return tokenString, claims.ExpiresAt, nil
}
