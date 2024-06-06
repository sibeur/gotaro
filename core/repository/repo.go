package repository

import (
	core_cache "github.com/sibeur/gotaro/core/common/cache"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	Driver *DriverRepository
	Rule   *RuleRepository
	Media  *MediaRepository
}

func NewRepository(mongoDB *mongo.Database, cache core_cache.CacheUseCase) *Repository {
	return &Repository{
		Driver: NewDriverRepository(mongoDB, cache),
		Rule:   NewRuleRepository(mongoDB, cache),
		Media:  NewMediaRepository(mongoDB, cache),
	}
}
