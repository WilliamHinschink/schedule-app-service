package config

import (
	"github.com/gofiber/fiber/v2"
	"schedule-app-service/utils"
	"strconv"
	"time"
)

func FiberConfig() fiber.Config {

	readTimeout, _ := strconv.Atoi(utils.GetEnvDefault("SERVER_READ_TIMEOUT", "30"))

	return fiber.Config{
		ReadTimeout: time.Second * time.Duration(readTimeout),
	}
}
