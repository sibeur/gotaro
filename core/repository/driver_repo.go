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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DriverRepository struct {
	db    *mongo.Database
	cache core_cache.CacheUseCase
}

func NewDriverRepository(db *mongo.Database, cache core_cache.CacheUseCase) *DriverRepository {
	return &DriverRepository{db: db, cache: cache}
}

func (u *DriverRepository) FindAll() ([]*entity.Driver, error) {
	ctx := context.TODO()
	var drivers []*entity.Driver
	cur, err := u.db.Collection(entity.Driver{}.GetCollName()).Find(ctx, bson.M{"deleted_at": nil})
	if err != nil {
		log.Printf("Error finding drivers: %v", err)
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var driver entity.Driver
		err := cur.Decode(&driver)
		if err != nil {
			log.Printf("Error decoding driver: %v", err)
			return nil, err
		}
		drivers = append(drivers, &driver)
	}
	return drivers, nil
}

func (u *DriverRepository) FindAllSimple() ([]*entity.Driver, error) {
	ctx := context.TODO()
	var drivers []*entity.Driver
	projection := bson.M{"_id": 1, "slug": 1, "name": 1, "type": 1}
	cur, err := u.db.Collection(entity.Driver{}.GetCollName()).Find(ctx, bson.M{"deleted_at": nil}, &options.FindOptions{
		Projection: projection,
	})
	if err != nil {
		log.Printf("Error finding drivers: %v", err)
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var driver entity.Driver
		err := cur.Decode(&driver)
		if err != nil {
			log.Printf("Error decoding driver: %v", err)
			return nil, err
		}
		drivers = append(drivers, &driver)
	}
	return drivers, nil
}

func (u *DriverRepository) Create(driver *entity.Driver) error {
	driver.ID = uuid.NewString()
	driver.CreatedAt = time.Now()
	driver.UpdatedAt = time.Now()
	_, err := u.db.Collection(entity.Driver{}.GetCollName()).InsertOne(context.TODO(), driver)
	if err != nil {
		return err
	}
	return nil
}

func (u *DriverRepository) Update(driver *entity.Driver) error {
	driver.UpdatedAt = time.Now()
	filter := bson.M{"slug": driver.Slug, "deleted_at": nil}
	_, err := u.db.Collection(entity.Driver{}.GetCollName()).UpdateOne(context.TODO(), filter, bson.M{"$set": driver})
	if err != nil {
		return err
	}
	return nil
}

func (u *DriverRepository) Delete(slug string) error {
	filter := bson.M{"slug": slug, "deleted_at": nil}
	data := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	_, err := u.db.Collection(entity.Driver{}.GetCollName()).UpdateOne(context.TODO(), filter, data)
	if err != nil {
		return err
	}
	return nil
}

func (u *DriverRepository) FindBySlug(slug string) (*entity.Driver, error) {
	var driver entity.Driver
	err := u.db.Collection(driver.GetCollName()).FindOne(context.TODO(), bson.M{"slug": slug, "deleted_at": nil}).Decode(&driver)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &driver, nil
}

func (u *DriverRepository) FindByID(id string) (*entity.Driver, error) {
	var driver entity.Driver
	err := u.db.Collection(driver.GetCollName()).FindOne(context.TODO(), bson.M{"_id": id, "deleted_at": nil}).Decode(&driver)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &driver, nil
}
