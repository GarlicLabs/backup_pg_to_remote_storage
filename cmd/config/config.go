package config

import (
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	GlobalStorageConfig Storage    `yaml:"storage" validate:"omitempty"`
	Databases           []Database `yaml:"databases" validate:"required,min=1,dive,customDatabaseValidation"`
}

type Database struct {
	Host          string  `yaml:"host" validate:"required"`
	Port          int     `yaml:"port" validate:"required,min=1"`
	Database      string  `yaml:"database" validate:"required,min=1"`
	Username      string  `yaml:"username" validate:"required,min=1"`
	Password      string  `yaml:"password" validate:"required,min=1"`
	Retention     int     `yaml:"retention" validate:"required,min=1"`
	StorageConfig Storage `yaml:"storageConfig" validate:"omitempty"`
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

func customDatabaseValidation(fl validator.FieldLevel) bool {
	config, ok := fl.Parent().Interface().(Config)
	if !ok {
		return false
	}
	if (config.GlobalStorageConfig != Storage{}) {
		return true
	}
	validate := validator.New()

	for _, database := range config.Databases {
		err := validate.Struct(database.StorageConfig.S3Config)
		if err != nil {
			return false
		}
	}
	return true
}

// Validates given configuration with github.com/go-playground/validator/v10 library
func Validate(cfg Config) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("customDatabaseValidation", customDatabaseValidation)
	err := validate.Struct(cfg)

	if err != nil {
		return err
	} else {
		log.Debugf("Validation succeeded")
	}
	return nil
}

// TODO needs unit test
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
