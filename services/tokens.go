package services

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/pkg/errors"
)

type TokensService interface {
	GenerateSignedToken(req *UserClaims) (string, *Claims, error)
	ParseSignedToken(tokenStr string) (*UserClaims, error)
}

type UserClaims struct {
	Id   string
	Role string
}

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type tokensService struct {
	validFor   time.Duration
	signingKey []byte
}

func NewTokensService(validFor time.Duration, signingKey []byte) TokensService {
	return &tokensService{
		validFor:   validFor,
		signingKey: signingKey,
	}
}

type Claims struct {
	jwt.StandardClaims
	Role string `json:"role"`
}

func (us *tokensService) GenerateSignedToken(req *UserClaims) (string, *Claims, error) {
	claims := &Claims{
		Role: req.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.NewTime(float64(time.Now().Add(us.validFor).Unix())),
			Subject:   req.Id,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	strToken, err := token.SignedString(us.signingKey)
	if err != nil {
		return "", nil, errors.Wrap(err, "could not sign token")
	}

	return strToken, claims, nil
}

func (us *tokensService) ParseSignedToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the expected algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return us.signingKey, nil
	})

	if err != nil {
		if _, ok := err.(*jwt.TokenExpiredError); ok {
			return nil, ErrExpiredToken
		}

		return nil, errors.Wrap(err, "could not parse claims from string token")
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return &UserClaims{
			Id:   claims.Subject,
			Role: claims.Role,
		}, nil
	}

	return nil, ErrInvalidToken
}
