package config

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
)

type Config interface {
	Port() string
	Context() string
	Database() (username string, password string, machine string, port int, database string)
	Access() (url string, right string)
	SecondsToken() int
}

type config struct {
	DbUser            string `json:"dbUsername"`
	DbPassword        string `json:"dbPassword"`
	DbMachine         string `json:"dbMachine"`
	DbPort            int    `json:"dbPort"`
	DbDatabase        string `json:"dbDatabase"`
	ServerPort        string `json:"port"`
	ServerContext     string `json:"context"`
	AccessRight       string `json:"accessRight"`
	AuthURL           string `json:"authURL"`
	SecondsTokenValid int    `json:"secondsTokenValid"`
}

func (c *config) SecondsToken() int {
	return c.SecondsTokenValid
}

func (c *config) Port() string {
	return c.ServerPort
}

func (c *config) Context() string {
	return c.ServerContext
}

func (c *config) Database() (string, string, string, int, string) {
	return c.DbUser, c.DbPassword, c.DbMachine, c.DbPort, c.DbDatabase
}

func (c *config) Access() (url string, right string) {
	return c.AuthURL, c.AccessRight
}

func LoadConfig() (Config, error) {
	// get environment variable
	fileName := os.Getenv("CONFIG_FILE")

	// read file
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Err(err)
		return nil, err
	}

	// parse to config
	var cfg config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Err(err)
		return nil, err
	}

	log.Info().Msg("config loaded successfully...")
	log.Info().Msg(fmt.Sprintf("config: DbUser: %s", cfg.DbUser))
	log.Info().Msg(fmt.Sprintf("config: DbMachine: %s", cfg.DbMachine))
	log.Info().Msg(fmt.Sprintf("config: DbPort: %d", cfg.DbPort))
	log.Info().Msg(fmt.Sprintf("config: DbDatabase: %s", cfg.DbDatabase))
	log.Info().Msg(fmt.Sprintf("config: ServerPort: %s", cfg.ServerPort))
	log.Info().Msg(fmt.Sprintf("config: ServerContext: %s", cfg.ServerContext))
	log.Info().Msg(fmt.Sprintf("config: AccessRight: %s", cfg.AccessRight))
	log.Info().Msg(fmt.Sprintf("config: AuthURL: %s", cfg.AuthURL))
	log.Info().Msg(fmt.Sprintf("config: SecondsTokenValid: %d", cfg.SecondsTokenValid))

	return &cfg, nil
}
