package handlers

import (
	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
	"log"
	"schedule-app-service/config"
	"schedule-app-service/models"
	"schedule-app-service/repository"
)

type Instance struct {
	log *log.Logger
}

func (i Instance) PostCustomer(c *fiber.Ctx) error {
	customer := &models.Customer{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(customer); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	db, err := config.NewCassandraDatabase(i.log).InitCluster()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	response, err := repository.NewCustomerRepository(i.log).PostCustomer(db, customer)
	if err != nil {
		// Return, if books not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":     true,
			"msg":       err.Error(),
			"count":     0,
			"customers": nil,
		})
	}
	return c.JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"customer": response,
	})
}

func (i Instance) GetCustomers(c *fiber.Ctx) error {
	i.log.Println("Handle GET Instance")

	db, err := config.NewCassandraDatabase(i.log).InitCluster()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// id := c.Request().Header.Peek("")
	customers, state, err := repository.NewCustomerRepository(i.log).GetCustomers(db, "")
	if err != nil {
		// Return, if books not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":     true,
			"msg":       "customers were not found",
			"count":     0,
			"customers": nil,
		})
	}

	c.Set("X-State-Id", state)

	return c.JSON(customers)
}

func (i Instance) GetCustomerById(ctx *fiber.Ctx) error {
	i.log.Println("Handle GET Instance")

	db, err := config.NewCassandraDatabase(i.log).InitCluster()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Id is required",
		})
	}

	idCustomer, _ := gocql.ParseUUID(id)

	customer, err := repository.NewCustomerRepository(i.log).GetCustomerById(db, idCustomer)
	if err != nil {
		// Return, if books not found.
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":     true,
			"msg":       err.Error(),
			"customers": nil,
		})
	}

	if customer == nil {
		// Return, if books not found.
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":     true,
			"msg":       "Customer not found",
			"customers": nil,
		})
	}

	return ctx.JSON(customer)
}

func (i Instance) DeleteCustomer(ctx *fiber.Ctx) error {
	db, err := config.NewCassandraDatabase(i.log).InitCluster()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Id is required",
		})
	}

	idCustomer, _ := gocql.ParseUUID(id)

	err = repository.NewCustomerRepository(i.log).DeleteCustomer(db, idCustomer)
	if err != nil {
		// Return, if books not found.
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":     true,
			"msg":       "customers were not found",
			"count":     0,
			"customers": nil,
		})
	}

	return ctx.Status(fiber.StatusAccepted).Send(nil)
}

func NewCustomer(l *log.Logger) *Instance {
	return &Instance{l}
}
