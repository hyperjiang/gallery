package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// Config - the config structure
type Config struct {
	Admin    AdminConfig
	Database DatabaseConfig
	Server   ServerConfig
}

var configDir string

// config file names
var (
	adminConf    = "admin"
	databaseConf = "database"
	serverConf   = "server"
)

// Dir - get the config dir
func Dir() string {
	if configDir == "" {
		if configDir = os.Getenv("CONFIG_DIR"); configDir == "" {
			configDir = "/runtime/config"
		}
	}
	return configDir
}

// Load configs from file
func Load(filename string, config interface{}) error {

	configFile := filepath.Join(Dir(), filename+".yml")

	log.Println(configFile)

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println(err)
		return err
	}

	return yaml.Unmarshal(data, config)
}

// LoadAll - load all configs
func LoadAll() *Config {

	var admin AdminConfig
	if err := Load(adminConf, &admin); err != nil {
		log.Println("Cannot load admin config file")
		return nil
	}

	var db DatabaseConfig
	if err := Load(databaseConf, &db); err != nil {
		log.Println("Cannot load database config file")
		return nil
	}

	var server ServerConfig
	if err := Load(serverConf, &server); err != nil {
		log.Println("Cannot load server config file")
		return nil
	}

	return &Config{
		Admin:    admin,
		Database: db,
		Server:   server,
	}
}
