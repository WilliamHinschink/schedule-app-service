package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"schedule-app-service/config"
	"schedule-app-service/middleware"
	"schedule-app-service/routes"
	"schedule-app-service/server"
)

func main() {
	l := log.New(os.Stdout, "schedule-api-service-", log.LstdFlags)

	configs := config.FiberConfig()

	app := fiber.New(configs)

	middleware.FiberMiddleware(app)

	routes.New(l).PublicRoutes(app) // Register a public routes for app.

	server.New(l).BootServerWithGracefulShutdown(app)
}
