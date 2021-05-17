package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/AntonAverchenkov/cards-http-service/internal/api"
	"github.com/AntonAverchenkov/cards-http-service/internal/state"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/hashicorp/go-multierror"
	"github.com/jessevdk/go-flags"
	"github.com/labstack/echo/v4"
)

type CommandLineOptions struct {
	Address             string `long:"address"                env:"ADDRESS"                description:"Listen to http traffic on this tcp address"      default:"localhost:8080"`
	SessionsPersistTo   string `long:"sessions-persist-to"    env:"SESSIONS_PERSIST_TO"    description:"Persist the sessions to this file on exit"       default:""`
	SessionsRestoreFrom string `long:"sessions-restore-from"  env:"SESSIONS_RESTORE_FROM"  description:"Restore the sessions from this file on startup"  default:""`
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
	/* */ log.Println("run(): cards-http-service begin")
	defer log.Println("run(): cards-http-service end")

	signalled := make(chan os.Signal, 1)
	signal.Notify(
		signalled,
		syscall.SIGHUP,  // sent to the process when its controlling terminal is closed
		syscall.SIGINT,  // sent to the process by its controlling terminal when a user wishes to interrupt the process
		syscall.SIGTERM, // sent to the process to request its termination
	)

	swagger, err := api.GetSwagger()
	if err != nil {
		return fmt.Errorf("could not load swagger spec: %w", err)
	}

	var (
		lock     sync.Mutex
		sessions *state.SessionManager
	)
	if cl.SessionsRestoreFrom != "" {
		log.Printf("run(): restoring sessions from %q\n", cl.SessionsRestoreFrom)

		sessions, err = state.Restore(cl.SessionsRestoreFrom)
		if err != nil {
			log.Printf("run(): could not restore sessions; starting new ones: %v\n", err)
			sessions = state.NewSessionManager()
		}
	} else {
		sessions = state.NewSessionManager()
	}

	handlers := handlers{
		lock:     &lock,
		sessions: sessions,
	}

	server := echo.New()
	server.Use(middleware.OapiRequestValidator(swagger))

	api.RegisterHandlers(server, &handlers)

	log.Printf("run(): starting to listen & serve on %q\n", cl.Address)

	done := make(chan bool)
	go func() {
		server.Start(cl.Address)
		close(done)
	}()

	// block until an interrupt signal or server exit
	select {
	case <-signalled:
		log.Println("run(): signalled, exiting")
	case <-done:
		log.Println("run(): the server has stopped")
	}

	if cl.SessionsPersistTo != "" {
		lock.Lock()
		defer lock.Unlock()

		log.Printf("run(): persisting sessions to %q\n", cl.SessionsPersistTo)

		if err := sessions.Persist(cl.SessionsPersistTo); err != nil {
			errs = multierror.Append(errs, fmt.Errorf("could not persist: %w", err))
		}
	}

	return nil
}
