package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	go_cache "github.com/sibeur/go-cache"
	"github.com/sibeur/gotaro/core/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApiClientRepository struct {
	db    *mongo.Database
	cache go_cache.Cache
}

func NewApiClientRepository(db *mongo.Database, cache go_cache.Cache) *ApiClientRepository {
	return &ApiClientRepository{db: db, cache: cache}
}

func (r *ApiClientRepository) FindByKey(key string) (*entity.APIClient, error) {
	ctx := context.TODO()

	var apiClient entity.APIClient
	err := r.db.Collection(apiClient.GetCollName()).FindOne(ctx, bson.M{"key": key}).Decode(&apiClient)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &apiClient, nil
}

func (r *ApiClientRepository) FindByScope(scope string) (*entity.APIClient, error) {
	ctx := context.TODO()
	var apiClient entity.APIClient

	err := r.db.Collection(apiClient.GetCollName()).FindOne(ctx, bson.M{"scopes": scope}).Decode(&apiClient)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &apiClient, nil
}

func (r *ApiClientRepository) FindByID(id string) (*entity.APIClient, error) {
	ctx := context.TODO()
	var apiClient entity.APIClient
	err := r.db.Collection(apiClient.GetCollName()).FindOne(ctx, bson.M{"_id": id}).Decode(&apiClient)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &apiClient, nil
}

func (r *ApiClientRepository) Create(apiClient *entity.APIClient) error {
	apiClient.CreatedAt = time.Now()
	apiClient.UpdatedAt = time.Now()
	apiClient.ID = uuid.NewString()
	ctx := context.TODO()
	_, err := r.db.Collection(entity.APIClient{}.GetCollName()).InsertOne(ctx, apiClient)
	if err != nil {
		return err
	}
	return nil
}
