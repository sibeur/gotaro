package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	go_cache "github.com/sibeur/go-cache"
	"github.com/sibeur/gotaro/core/common"
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
	key := fmt.Sprintf(common.CacheGetMediaKey, ruleSlug, fileAliasName)
	mediaCache, _ := u.cache.Get(key)
	if mediaCache != "" {
		err := media.FromJSONString(mediaCache)
		if err != nil {
			log.Printf("Error decoding media: %v", err)
		}
		return &media, nil
	}
	filter := bson.M{"rule_slug": ruleSlug, "file_alias_name": fileAliasName, "deleted_at": nil}
	err := u.db.Collection(media.GetCollName()).FindOne(context.TODO(), filter).Decode(&media)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	go func(media *entity.Media) {
		mediaCache, err := media.ToJSONString()
		if err != nil {
			log.Printf("Error marshal media to cache: %v", err)
		}
		err = u.cache.SetWithExpire(fmt.Sprintf(common.CacheGetMediaKey, media.RuleSlug, media.FileAliasName), mediaCache, common.DefaultGetMediaCacheTTL)
		if err != nil {
			log.Println("Error set media to cache: ", err)
		}
	}(&media)
	return &media, nil
}

func (u *MediaRepository) SetSignedUrl(ruleSlug, fileAliasName, signedUrl string) error {
	filter := bson.M{"rule_slug": ruleSlug, "file_alias_name": fileAliasName, "deleted_at": nil}
	data := bson.M{"$set": bson.M{"file_path": signedUrl}}
	_, err := u.db.Collection(entity.Media{}.GetCollName()).UpdateOne(context.TODO(), filter, data)
	if err != nil {
		return err
	}
	return nil
}

func (u *MediaRepository) GetCachedSignedUrl(ruleSlug, fileAliasName string) (string, error) {
	key := fmt.Sprintf(common.CacheMediaSignedUrlKey, ruleSlug, fileAliasName)
	signedUrl, err := u.cache.Get(key)
	if err != nil {
		return "", err
	}
	return signedUrl, nil
}

func (u *MediaRepository) SetCachedSignedUrl(ruleSlug, fileAliasName, signedUrl string) {
	key := fmt.Sprintf(common.CacheMediaSignedUrlKey, ruleSlug, fileAliasName)
	err := u.cache.SetWithExpire(key, signedUrl, common.DefaultGetMediaCacheTTL)
	if err != nil {
		log.Printf("Error set media signed url to cache: %v", err)
	}
}
