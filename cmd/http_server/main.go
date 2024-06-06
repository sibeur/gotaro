package main

import (
	"context"

	app_http "github.com/sibeur/gotaro/apps/http"
	"github.com/sibeur/gotaro/core/common"
	core_cache "github.com/sibeur/gotaro/core/common/cache"
	"github.com/sibeur/gotaro/core/common/driver"
	core_db "github.com/sibeur/gotaro/core/db"
	core_repository "github.com/sibeur/gotaro/core/repository"
	core_service "github.com/sibeur/gotaro/core/service"

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
	cache := core_cache.NewCache()

	// load driver manager
	driverManager := driver.NewDriverManager()

	// load reapository
	repo := core_repository.NewRepository(mongoDB, cache)

	// load service
	service := core_service.NewService(repo, driverManager)

	if err := service.Driver.LoadDriverManager(); err != nil {
		panic(err)
	}

	// check temporary folder
	temporaryFolder := common.TemporaryFolder
	if err := common.CreateFolder(temporaryFolder); err != nil {
		panic(err)
	}

	// start http app
	app_http.NewFiberApp(service).Run()

}
