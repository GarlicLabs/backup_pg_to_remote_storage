package config

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestValidConfigWithGlobalStorageConfig(t *testing.T) {
	var dbCfg []Database
	dbCfg = append(dbCfg, validDbConfigWithoutStorage0)
	dbCfg = append(dbCfg, validDbConfigWithoutStorage1)

	cfg := Config{validStorageConfig, dbCfg}

	err := Validate(cfg)
	if err != nil {
		log.Error(err)
		t.Fatalf("Validation should not fail as config is valid")
	}
}

func TestValidConfigWithDatabaseStorageConfig(t *testing.T) {
	var dbCfg []Database
	dbCfg = append(dbCfg, validDbConfigWithStorage0)
	dbCfg = append(dbCfg, validDbConfigWithStorage1)

	cfg := Config{Storage{}, dbCfg}

	err := Validate(cfg)
	if err != nil {
		log.Error(err)
		t.Fatalf("Validation should not fail as config is valid")
	}
}

func TestInvalidGlobalStorageConfig(t *testing.T) {
	var dbCfg []Database
	dbCfg = append(dbCfg, validDbConfigWithoutStorage0)

	cfg := Config{invalidStorageConfig, dbCfg}

	err := Validate(cfg)
	if err == nil {
		t.Fatalf("Validation should fail, storage config is invalid")
	}
}

func TestInvalidDatabaseConfig(t *testing.T) {
	var dbCfg []Database
	dbCfg = append(dbCfg, invalidDbConfigWithoutStorage0)

	cfg := Config{validStorageConfig, dbCfg}

	err := Validate(cfg)
	if err == nil {
		t.Fatalf("Validation should fail, database config is invalid")
	}
}

func TestNoConfigIsGiven(t *testing.T) {
	cfg := Config{Storage{}, []Database{}}
	err := Validate(cfg)
	if err == nil {
		t.Fatalf("Validation should fail, as no config is given")
	}
}

func TestInvalidDatabaseConfigAndValidStorageConfig(t *testing.T) {
	var dbCfg []Database
	dbCfg = append(dbCfg, invalidDbConfigWithoutStorage0)

	cfg := Config{validStorageConfig, dbCfg}
	err := Validate(cfg)
	if err == nil {
		t.Fatalf("Validation should fail, as database config is invalid")
	}
}

func TestInvalidGlobalStorageAndNoStorageConfig(t *testing.T) {
	var dbCfg []Database
	dbCfg = append(dbCfg, validDbConfigWithoutStorage0)

	cfg := Config{invalidStorageConfig, dbCfg}
	err := Validate(cfg)
	if err == nil {
		t.Fatalf("Validation should fail, as no valid storage config and invalid global storage config")
	}
}

func TestNoGlobalStorageConfigAndNoDbStorageConfig(t *testing.T) {
	var dbCfg []Database
	dbCfg = append(dbCfg, validDbConfigWithoutStorage0)

	cfg := Config{Storage{}, dbCfg}
	err := Validate(cfg)
	if err == nil {
		t.Fatalf("Validation should fail, as no StorageConfig was provied")
	}
}

func TestNoGlobalStorageConfigAndNoS3ConfigOnDb(t *testing.T) {
	var dbCfg []Database
	validDbConfigWithoutStorage0.StorageConfig = invalidStorageConfigNoConfig
	dbCfg = append(dbCfg, validDbConfigWithoutStorage0)

	cfg := Config{Storage{}, dbCfg}
	err := Validate(cfg)
	if err == nil {
		t.Fatalf("Validation should fail, as no StorageConfig was provied")
	}
}

var validS3Config = S3{
	Endpoint:  "localhost",
	AccessKey: "accessKey",
	SecretKey: "secretKey",
	Bucket:    "bucket",
}

var validStorageConfig = Storage{
	S3Config: validS3Config,
}

var invalidStorageConfigNoConfig = Storage{}

var invalidS3Config = S3{
	AccessKey: "accessKey",
	Bucket:    "bucket",
}

var invalidStorageConfig = Storage{
	S3Config: invalidS3Config,
}

var invalidDbConfigWithoutStorage0 = Database{
	Host:      "xxx",
	Port:      5432,
	Password:  "password",
	Retention: 30,
}

var validDbConfigWithoutStorage0 = Database{
	Host:      "postgres.example.com",
	Port:      5432,
	Database:  "Database",
	Username:  "username",
	Password:  "password",
	Retention: 30,
}

var validDbConfigWithoutStorage1 = Database{
	Host:      "127.0.0.1",
	Port:      5432,
	Database:  "Database",
	Username:  "username",
	Password:  "password",
	Retention: 30,
}

var validDbConfigWithStorage0 = Database{
	Host:          "postgres.example.com",
	Port:          5432,
	Database:      "Database",
	Username:      "username",
	Password:      "password",
	Retention:     30,
	StorageConfig: validStorageConfig,
}

var validDbConfigWithStorage1 = Database{
	Host:          "127.0.0.1",
	Port:          5432,
	Database:      "Database",
	Username:      "username",
	Password:      "password",
	Retention:     30,
	StorageConfig: validStorageConfig,
}
