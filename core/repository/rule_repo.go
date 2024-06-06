package repository

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	core_cache "github.com/sibeur/gotaro/core/common/cache"
	"github.com/sibeur/gotaro/core/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RuleRepository struct {
	db    *mongo.Database
	cache core_cache.CacheUseCase
}

func NewRuleRepository(db *mongo.Database, cache core_cache.CacheUseCase) *RuleRepository {
	return &RuleRepository{db: db, cache: cache}
}

func (u *RuleRepository) FindAll() ([]*entity.Rule, error) {
	ctx := context.TODO()
	var rules []*entity.Rule
	cur, err := u.db.Collection(entity.Rule{}.GetCollName()).Find(ctx, bson.M{"deleted_at": nil})
	if err != nil {
		log.Printf("Error finding rules: %v", err)
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var rule entity.Rule
		err := cur.Decode(&rule)
		if err != nil {
			log.Printf("Error decoding rule: %v", err)
			return nil, err
		}
		rules = append(rules, &rule)
	}
	return rules, nil
}

func (u *RuleRepository) Create(rule *entity.Rule) error {
	rule.ID = uuid.NewString()
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()
	_, err := u.db.Collection(entity.Rule{}.GetCollName()).InsertOne(context.TODO(), rule)
	if err != nil {
		return err
	}
	return nil
}

func (u *RuleRepository) Update(rule *entity.Rule) error {
	rule.UpdatedAt = time.Now()
	filter := bson.M{"slug": rule.Slug, "deleted_at": nil}
	_, err := u.db.Collection(entity.Rule{}.GetCollName()).UpdateOne(context.TODO(), filter, bson.M{"$set": rule})
	if err != nil {
		return err
	}
	return nil
}

func (u *RuleRepository) Delete(slug string) error {
	filter := bson.M{"slug": slug, "deleted_at": nil}
	data := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	_, err := u.db.Collection(entity.Rule{}.GetCollName()).UpdateOne(context.TODO(), filter, data)
	if err != nil {
		return err
	}
	return nil
}

func (u *RuleRepository) FindBySlug(slug string) (*entity.Rule, error) {
	var rule entity.Rule
	err := u.db.Collection(entity.Rule{}.GetCollName()).FindOne(context.TODO(), bson.M{"slug": slug, "deleted_at": nil}).Decode(&rule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &rule, nil
}
