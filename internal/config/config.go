package config

import (
	"flag"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DB         DBConfig
	Port       string        `yaml:"port"`
	Timeout    time.Duration `yaml:"timeout"`
	StorageURL string        `yaml:"storage_service_url"`
}

type DBConfig struct {
	Host     string
	DBName   string
	User     string
	Password string
}

func MustLoadConfig() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("Config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("Config file does not exist: " + path)
	}

	var cfg Config
	data, _ := os.ReadFile(path)

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic("Failed to unmarshal yaml")
	}

	dbConfig := getDBconfig()

	cfg.DB = dbConfig

	return &cfg
}

func fetchConfigPath() string {
	var result string

	// --config="path/to/config.yaml"
	flag.StringVar(&result, "config", "./configs/config.yaml", "path to config file")
	flag.Parse()

	if result == "" {
		result = os.Getenv("CONFIG_PATH")
	}

	return result
}

func getDBconfig() DBConfig {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env file")
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		panic("DB_HOST is not set")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		panic("DB_NAME is not set")
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		panic("DB_USER is not set")
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		panic("DB_PASSWORD is not set")
	}

	return DBConfig{
		Host:     host,
		DBName:   dbName,
		User:     user,
		Password: password,
	}
}
