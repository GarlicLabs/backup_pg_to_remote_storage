package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	GlobalStorageConfig Storage    `yaml:"storage" validate:"omitempty"`
	Databases           []Database `yaml:"databases" validate:"required,min=1,dive,validateStorageIsSet"`
}

type Database struct {
	Host          string  `yaml:"host" validate:"required"`
	Port          int     `yaml:"port" validate:"required,min=1"`
	Database      string  `yaml:"database" validate:"required,min=1"`
	Username      string  `yaml:"username" validate:"required,min=1"`
	Password      string  `yaml:"password" validate:"required,min=1"`
	Retention     int     `yaml:"retention" validate:"required,min=1"`
	StorageConfig Storage `yaml:"storage" validate:"omitempty"`
}

type Storage struct {
	//Change validate required to required_if, when we support for more storage providers
	S3Config S3 `yaml:"s3" validate:"required"`
}

type S3 struct {
	Endpoint  string `yaml:"endpoint" validate:"required"`
	AccessKey string `yaml:"accessKey" validate:"required"`
	SecretKey string `yaml:"secretKey" validate:"required"`
	Bucket    string `yaml:"bucket" validate:"required"`
}

func GetConfig() Config {
	log.Info("Reading config file")
	configFile := getConfigPath()
	f, err := os.ReadFile(configFile)

	if err != nil {
		log.Panic(err)
	}

	fmt.Print(f)

	var config Config
	err = yaml.Unmarshal(f, &config)

	if err != nil {
		log.Panic(err)
	}

	return config
}

func getConfigPath() string {
	configPath := os.Getenv("BACKUP_PG_CONFIG_PATH")
	return configPath
}

func validateStorageIsSet(fl validator.FieldLevel) bool {
	config, ok := fl.Parent().Interface().(Config)
	if !ok {
		log.Errorf("Not ok on validateStorageIsSet")
		return false
	}
	if (config.GlobalStorageConfig != Storage{}) {
		log.Debugf("GlobalStorageConfig is set")
		return true
	}
	validate := validator.New()
	for _, database := range config.Databases {
		err := validate.Struct(database.StorageConfig.S3Config)
		if err != nil {
			log.Errorf("%+v\n", err)
			return false
		}
	}
	return true
}

// Validates given configuration with github.com/go-playground/validator/v10 library
func Validate(cfg Config) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("validateStorageIsSet", validateStorageIsSet)
	err := validate.Struct(cfg)

	if err != nil {
		return err
	} else {
		log.Debugf("Validation succeeded")
	}
	return nil
}

func MapGlobalStorageToDbIfNotSet(databaseConfig Database, GlobalS3Config Storage) Database {
	builtConfig := databaseConfig
	if databaseConfig.StorageConfig.S3Config.Endpoint == "" {
		builtConfig.StorageConfig.S3Config.Endpoint = GlobalS3Config.S3Config.Endpoint
	}
	if databaseConfig.StorageConfig.S3Config.AccessKey == "" {
		builtConfig.StorageConfig.S3Config.AccessKey = GlobalS3Config.S3Config.AccessKey
	}
	if databaseConfig.StorageConfig.S3Config.SecretKey == "" {
		builtConfig.StorageConfig.S3Config.SecretKey = GlobalS3Config.S3Config.SecretKey
	}
	if databaseConfig.StorageConfig.S3Config.Bucket == "" {
		builtConfig.StorageConfig.S3Config.Bucket = GlobalS3Config.S3Config.Bucket
	}
	return builtConfig
}
