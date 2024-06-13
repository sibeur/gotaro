package service

import (
	"errors"

	"github.com/sibeur/gotaro/core/common"
	"github.com/sibeur/gotaro/core/entity"
	"github.com/sibeur/gotaro/core/repository"
	"golang.org/x/crypto/bcrypt"
)

type ApiClientService struct {
	repo *repository.Repository
}

func NewApiClientService(repo *repository.Repository) *ApiClientService {
	return &ApiClientService{repo: repo}
}

func (u *ApiClientService) GenerateFirstSuperAdmin() (string, string, error) {
	client, err := u.repo.APIClient.FindByScope(common.APIClientSuperAdminScope)
	if err != nil {
		return "", "", err
	}
	if client != nil {
		return client.Key, "", errors.New(common.ErrAPIClientAlreadyExistMsg)
	}

	apikey := common.RandomString(16)
	secretKey := common.RandomString(32)

	// hashedSecretKey with bcrypt
	hashedSecretKey, err := bcrypt.GenerateFromPassword([]byte(secretKey), bcrypt.DefaultCost)

	if err != nil {
		return "", "", err
	}

	client = &entity.APIClient{
		Key:    apikey,
		Secret: string(hashedSecretKey),
		Scopes: []string{common.APIClientSuperAdminScope},
	}

	err = u.repo.APIClient.Create(client)
	if err != nil {
		return "", "", err
	}

	return apikey, secretKey, nil
}

func (u *ApiClientService) FindByID(id string) (*entity.APIClient, error) {
	return u.repo.APIClient.FindByID(id)
}

func (u *ApiClientService) FindByKey(key string) (*entity.APIClient, error) {
	return u.repo.APIClient.FindByKey(key)
}
