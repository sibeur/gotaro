package service

import (
	"github.com/sibeur/gotaro/core/common/driver"
	"github.com/sibeur/gotaro/core/repository"
)

type Service struct {
	Rule      *RuleService
	Driver    *DriverService
	Media     *MediaService
	ApiClient *ApiClientService
	Auth      *AuthService
}

func NewService(repo *repository.Repository, driverManager *driver.DriverManager) *Service {
	return &Service{
		Rule:      NewRuleService(repo),
		Driver:    NewDriverService(repo, driverManager),
		Media:     NewMediaService(repo, driverManager),
		ApiClient: NewApiClientService(repo),
		Auth:      NewAuthService(repo),
	}
}
