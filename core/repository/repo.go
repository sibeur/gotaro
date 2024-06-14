package repository

import (
	go_cache "github.com/sibeur/go-cache"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	Driver    *DriverRepository
	Rule      *RuleRepository
	Media     *MediaRepository
	APIClient *ApiClientRepository
	Auth      *AuthRepository
}

func NewRepository(mongoDB *mongo.Database, cache go_cache.Cache) *Repository {
	return &Repository{
		Driver:    NewDriverRepository(mongoDB, cache),
		Rule:      NewRuleRepository(mongoDB, cache),
		Media:     NewMediaRepository(mongoDB, cache),
		APIClient: NewApiClientRepository(mongoDB, cache),
		Auth:      NewAuthRepository(cache),
	}
}
