package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AntonAverchenkov/cards-http-service/internal/api"
	"github.com/AntonAverchenkov/cards-http-service/internal/game"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/jessevdk/go-flags"
	"github.com/labstack/echo/v4"
)

type CommandLineOptions struct {
	Address string `long:"address"  env:"ADDRESS"  description:"Listen to http traffic on this tcp address"   default:"localhost:8080"`
}

func main() {
	var cl CommandLineOptions

	_, err := flags.Parse(&cl)
	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			os.Exit(0)
		}
		log.Fatalf("Error :: command-line argument parsing failed: %v\n", err)
	}

	if err := run(cl); err != nil {
		log.Fatalf("Error :: %v\n", err)
	}
}

func run(cl CommandLineOptions) (errs error) {
	/* */ log.Printf("run(): cards-http-service begin")
	defer log.Printf("run(): cards-http-service end")

	swagger, err := api.GetSwagger()
	if err != nil {
		return fmt.Errorf("could not load swagger spec: %w", err)
	}

	handlers := handlers{
		deck: game.NewDeck(),
	}

	server := echo.New()
	server.Use(middleware.OapiRequestValidator(swagger))

	api.RegisterHandlers(server, &handlers)

	log.Printf("run(): starting to listen & serve on %q\n", cl.Address)

	if err := server.Start(cl.Address); err != nil {
		return fmt.Errorf("could not start the server @ %q: %w", cl.Address, err)
	}

	return nil
}
