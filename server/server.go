package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Config structure for the server
type Config struct {
	Port     int    // port number
	Host     string // hostname
	CertFile string // TLS Certificate file name (.pem) (optional)
	KeyFile  string // TLS Certificate key file name (.key) (optional)
}

// Start the server
// @title           Challenge Bravo - Bravo-Server API
// @version         1.0
// @description     This is the challenge bravo sample server
// @contact.name    André Gustavo de Andrade
// @contact.url     https://github.com/aandrade1234/challenge-bravo
// @contact.email   aandrade1234@gmail.com
// @license.name    Apache 2.0
// @license.url     https://www.apache.org/licenses/LICENSE-2.0.html
// @host            localhost:8080
// @BasePath        /api/v1
func Start(config Config) {

	// Create the server
	app := fiber.New()

	// Middleware used to recover from a panic event
	app.Use(recover.New())

	// Enable CORS Middleware to accept any origin
	app.Use(cors.New())

	// Create server routes
	createRoutes(app)

	// Start the server on a new go goroutine
	go func() {
		var err error
		addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
		if len(config.CertFile) > 0 && len(config.KeyFile) > 0 {
			err = app.ListenTLS(addr, config.CertFile, config.KeyFile)
		} else {
			err = app.Listen(addr)
		}
		if err != nil {
			log.Panicln(err)
		}
	}()

	// create a channel and wait for TERM/INT signal from the OS
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	_ = <-c

	// Stop the server gracefully
	_ = app.Shutdown()
}
