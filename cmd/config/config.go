package config

import (
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	GlobalStorageConfig StorageConfig    `yaml:"storage"`
	Databases           []DatabaseConfig `yaml:"databases" validate:"required"`
}

type DatabaseConfig struct {
	Host          string        `yaml:"host" validate:"required,fqdn"`
	Port          int           `yaml:"port" validate:"required"`
	Database      string        `yaml:"database" validate:"required"`
	Username      string        `yaml:"username" validate:"required"`
	Password      string        `yaml:"password" validate:"required"`
	Retention     int           `yaml:"retention" validate:"required"`
	StorageConfig StorageConfig `yaml:"storageConfig"`
}

type StorageConfig struct {
	S3Config S3Config `yaml:"s3"`
}

type S3Config struct {
	Endpoint  string `yaml:"endpoint" validate:"required"`
	AccessKey string `yaml:"accessKey" validate:"required"`
	SecretKey string `yaml:"secretKey" validate:"required"`
	Bucket    string `yaml:"bucket" validate:"required"`
}

// Validates given configuration with github.com/go-playground/validator/v10 library
func ValidateConfig(config Config) error {
	validate := validator.New()
	err := validate.Struct(config)
	if err != nil {
		log.Errorf("Validation failed")
		return err
	} else {
		log.Infof("Validation succeeded")
	}
	return nil
}

// TODO needs unit test
func PutGlobalStorageToDbIfNotSet(databaseConfig DatabaseConfig, GlobalS3Config StorageConfig) DatabaseConfig {
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
