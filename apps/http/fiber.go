package http

import (
	"os"

	"github.com/sibeur/gotaro/apps/http/handler"
	"github.com/sibeur/gotaro/core/common"
	"github.com/sibeur/gotaro/core/service"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

// FiberApp represents a Fiber application.
type FiberApp struct {
	Instance      *fiber.App
	Svc           *service.Service
	ruleHandler   *handler.RuleHandler
	driverHandler *handler.DriverHandler
	mediaHandler  *handler.MediaHandler
	authHandlerV1 *handler.AuthHandlerV1
}

// NewFiberApp creates a new instance of FiberApp.
func NewFiberApp(service *service.Service) *FiberApp {
	instance := fiber.New(fiber.Config{
		ErrorHandler: common.FiberDefaultErrorHandler,
	})
	return &FiberApp{
		Instance:      instance,
		Svc:           service,
		ruleHandler:   handler.NewRuleHandler(instance, service),
		driverHandler: handler.NewDriverHandler(instance, service),
		mediaHandler:  handler.NewMediaHandler(instance, service),
		authHandlerV1: handler.NewAuthHandlerV1(instance, service),
	}
}

// beforeMiddlewares sets up the middlewares to be executed before the main request handler.
func (f *FiberApp) beforeMiddlewares() {
	appEnv := os.Getenv("APP_ENV")

	// Create a zap logger
	logger, _ := zap.NewDevelopment()
	if appEnv == "prod" {
		logger, _ = zap.NewProduction()
	}

	// Add fiberzap middleware
	f.Instance.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))

	// Add recover middleware
	f.Instance.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
}

// afterMiddlewares sets up the middlewares to be executed after the main request handler.
func (f *FiberApp) afterMiddlewares() {

}

// Run starts the Fiber application and listens for incoming requests.
func (f *FiberApp) Run() {
	f.beforeMiddlewares()
	f.Instance.Get("/", func(c *fiber.Ctx) error {
		panic("Error")
		return c.SendString("Hello, World!")
	})
	f.authHandlerV1.Router()
	f.ruleHandler.Router()
	f.driverHandler.Router()
	f.mediaHandler.Router()
	f.afterMiddlewares()
	if err := f.Instance.Listen(":3000"); err != nil {
		panic(err)
	}
}
