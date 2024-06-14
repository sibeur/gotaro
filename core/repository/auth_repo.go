package repository

import (
	"errors"
	"os"
	"strconv"
	"time"

	go_cache "github.com/sibeur/go-cache"
	"github.com/sibeur/gotaro/core/common"
	"github.com/sibeur/gotaro/core/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthRepository struct {
	cache go_cache.Cache
}

func NewAuthRepository(cache go_cache.Cache) *AuthRepository {
	return &AuthRepository{cache: cache}
}

func (u *AuthRepository) CreateToken(subject string, audiences []string) (*entity.Auth, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New(common.ErrJWTSecretNotFoundMsg)
	}

	accessToken, err := u.createAccessToken(subject, audiences, secret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.createRefreshToken(subject, audiences, secret)
	if err != nil {
		return nil, err
	}
	return &entity.Auth{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *AuthRepository) ValidateToken(token string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New(common.ErrJWTSecretNotFoundMsg)
	}

	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New(common.ErrJWTTokenInvalidMsg)
		}
		return nil, err
	}
	return t, nil
}

func (u *AuthRepository) createAccessToken(subject string, audiences []string, secret string) (*entity.AccessToken, error) {
	expMinutes := 5
	if os.Getenv("ACCESS_TOKEN_EXPIRY_MINUTES") != "" {
		expMinutes, _ = strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY_MINUTES"))
	}
	now := time.Now()
	iat := jwt.NewNumericDate(now)
	// Set the expiration time to 5 minutes from now
	exp := jwt.NewNumericDate(now.Add(time.Minute * time.Duration(expMinutes)))
	claims := &jwt.RegisteredClaims{
		Subject:   subject,
		Issuer:    common.JWTIssuerAccessToken,
		ExpiresAt: exp,
		Audience:  audiences,
		IssuedAt:  iat,
		ID:        uuid.NewString(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	return &entity.AccessToken{
		Token:     ss,
		IssuedAt:  iat.Unix(),
		ExpiredAt: exp.Unix(),
	}, nil
}

func (u *AuthRepository) createRefreshToken(subject string, audiences []string, secret string) (*entity.RefreshToken, error) {
	// expMinutes in 7 days
	expMinutes := 60 * 24 * 7
	if os.Getenv("REFRESH_TOKEN_EXPIRY_MINUTES") != "" {
		expMinutes, _ = strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY_MINUTES"))
	}
	now := time.Now()
	iat := jwt.NewNumericDate(now)
	// Set the expiration time to 7 days from now
	exp := jwt.NewNumericDate(now.Add(time.Minute * time.Duration(expMinutes)))
	claims := &jwt.RegisteredClaims{
		Subject:   subject,
		Issuer:    common.JWTIssuerRefreshToken,
		ExpiresAt: exp,
		Audience:  audiences,
		IssuedAt:  iat,
		ID:        uuid.NewString(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	return &entity.RefreshToken{
		Token:     ss,
		IssuedAt:  iat.Unix(),
		ExpiredAt: exp.Unix(),
	}, nil
}
