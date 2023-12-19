package proc

import (
	"context"
	"database/sql"
	"github.com/dgb9/db-account-server/internal/config"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Start() error {
	log.Info().Msg("starting...")
	log.Info().Msg("loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	db, err := config.GetConnectionPool(cfg)
	if err != nil {
		return err
	}

	defer close(db)

	// create the router
	router := Router(cfg, db)

	// then, start the service
	err = startServer(cfg, router)

	return err
}

func startServer(config config.Config, router *mux.Router) error {
	srv := &http.Server{
		Addr: config.Port(),

		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Info().Msg(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	_ = srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Info().Msg("shutting down")
	os.Exit(0)

	return nil
}

func close(db *sql.DB) {
	if db != nil {
		_ = db.Close()
	}
}
