package routes

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"schedule-app-service/handlers"
)

type Instance struct {
	log *log.Logger
}

func (i Instance) PublicRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	customers := handlers.NewCustomer(i.log)
	route.Get("/customers", customers.GetCustomers)
	route.Post("/customers", customers.PostCustomer)
	route.Get("/customers/:id", customers.GetCustomerById)
	route.Delete("/customers/:id", customers.DeleteCustomer)
}

func New(l *log.Logger) *Instance {
	return &Instance{l}
}
