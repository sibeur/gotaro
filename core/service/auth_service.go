package service

import (
	"errors"

	"github.com/sibeur/gotaro/core/common"
	"github.com/sibeur/gotaro/core/entity"
	"github.com/sibeur/gotaro/core/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.Repository
}

func NewAuthService(repo *repository.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (u *AuthService) Login(apiKey string, secretKey string) (*entity.Auth, error) {
	apiClient, err := u.repo.APIClient.FindByKey(apiKey)
	if err != nil {
		return nil, errors.New(common.ErrAuthenticationFailedMsg)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(apiClient.Secret), []byte(secretKey)); err != nil {
		return nil, errors.New(common.ErrAuthenticationFailedMsg)
	}
	auth, err := u.repo.Auth.CreateToken(apiClient.ID, apiClient.Scopes)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func (u *AuthService) RefreshToken(refreshToken string) (*entity.Auth, error) {
	token, err := u.repo.Auth.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	claims := token.Claims
	subject, err := claims.GetSubject()
	if err != nil {
		return nil, err
	}

	issuer, err := claims.GetIssuer()
	if err != nil {
		return nil, err
	}
	if issuer != common.JWTIssuerRefreshToken {
		return nil, errors.New(common.ErrJWTTokenInvalidMsg)
	}

	apiClient, err := u.repo.APIClient.FindByID(subject)
	if err != nil {
		return nil, err
	}

	if apiClient == nil {
		return nil, errors.New(common.ErrJWTTokenInvalidMsg)
	}

	auth, err := u.repo.Auth.CreateToken(subject, apiClient.Scopes)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func (u *AuthService) ValidateToken(token string) (*jwt.Token, error) {
	return u.repo.Auth.ValidateToken(token)
}
