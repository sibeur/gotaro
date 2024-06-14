package main

import (
	"context"

	core_db "github.com/sibeur/gotaro/core/db"
	core_repository "github.com/sibeur/gotaro/core/repository"
	core_service "github.com/sibeur/gotaro/core/service"

	go_cache "github.com/sibeur/go-cache"
	app_pubsub "github.com/sibeur/gotaro/apps/pubsub"
	"github.com/sibeur/gotaro/core/common/driver"

	"github.com/joho/godotenv"
)

func main() {
	// This is the entry point of the application
	// dotenv load
	godotenv.Load()

	// load mongodb
	mongoDB, err := core_db.NewMongoDBConnection()
	if err != nil {
		panic(err)
	}
	defer mongoDB.Client().Disconnect(context.Background())

	//load cache
	cache := go_cache.NewCache()

	// load driver manager
	driverManager := driver.NewDriverManager()

	// load reapository
	repo := core_repository.NewRepository(mongoDB, cache)

	// load service
	service := core_service.NewService(repo, driverManager)

	// start pubsub app
	app_pubsub.NewApp(service).Run()
}
