package main

import (
	"github.com/dgb9/db-account-server/internal/proc"
	"github.com/rs/zerolog/log"
)

func main() {
	err := proc.Start()

	if err != nil {
		log.Error().Err(err).Msg("Application exited with error")
	}
}
