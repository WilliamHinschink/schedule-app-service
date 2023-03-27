package server

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
	"schedule-app-service/utils"
	"syscall"
)

type Instance struct {
	log *log.Logger
}

func (i Instance) BootServerWithGracefulShutdown(a *fiber.App) {
	// Create channel for idle connections.
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint,
			os.Interrupt,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGABRT,
			syscall.SIGQUIT,
		) // Catch OS signals.
		<-sigint

		// Received an interrupt signal, shutdown.
		if err := a.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			i.log.Printf("Oops... Server is not shutting down! Reason: %v", err)
		}
		i.log.Printf("Server Shutdown Gracefully")

		close(idleConnsClosed)
	}()

	// Run server.
	if err := a.Listen(utils.GetEnvDefault("SERVER_URL", ":8080")); err != nil {
		i.log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnsClosed
}

func (i Instance) BootServer(a *fiber.App) {
	// Run server.
	if err := a.Listen(utils.GetEnvDefault("SERVER_URL", ":8080")); err != nil {
		i.log.Printf("Oops... Server is not running! Reason: %v", err)
	}
}

func New(l *log.Logger) *Instance {
	return &Instance{l}
}
