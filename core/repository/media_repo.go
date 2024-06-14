package repository

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	go_cache "github.com/sibeur/go-cache"
	"github.com/sibeur/gotaro/core/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MediaRepository struct {
	db    *mongo.Database
	cache go_cache.Cache
}

func NewMediaRepository(db *mongo.Database, cache go_cache.Cache) *MediaRepository {
	return &MediaRepository{db: db, cache: cache}
}

func (u *MediaRepository) FindAll() ([]*entity.Media, error) {
	ctx := context.TODO()
	var medias []*entity.Media
	cur, err := u.db.Collection(entity.Media{}.GetCollName()).Find(ctx, bson.M{"deleted_at": nil})
	if err != nil {
		log.Printf("Error finding medias: %v", err)
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var media entity.Media
		err := cur.Decode(&media)
		if err != nil {
			log.Printf("Error decoding media: %v", err)
			return nil, err
		}
		medias = append(medias, &media)
	}
	return medias, nil
}

func (u *MediaRepository) Create(media *entity.Media) error {
	media.ID = uuid.NewString()
	media.CreatedAt = time.Now()
	media.UpdatedAt = time.Now()
	_, err := u.db.Collection(entity.Media{}.GetCollName()).InsertOne(context.TODO(), media)
	if err != nil {
		return err
	}
	return nil
}

func (u *MediaRepository) Delete(ruleSlug, fileAliasName string) error {
	filter := bson.M{"rule_slug": ruleSlug, "file_alias_name": fileAliasName, "deleted_at": nil}
	data := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	_, err := u.db.Collection(entity.Media{}.GetCollName()).UpdateOne(context.TODO(), filter, data)
	if err != nil {
		return err
	}
	return nil
}

func (u *MediaRepository) FindMedia(ruleSlug, fileAliasName string) (*entity.Media, error) {
	var media entity.Media
	filter := bson.M{"rule_slug": ruleSlug, "file_alias_name": fileAliasName, "deleted_at": nil}
	err := u.db.Collection(media.GetCollName()).FindOne(context.TODO(), filter).Decode(&media)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &media, nil
}
