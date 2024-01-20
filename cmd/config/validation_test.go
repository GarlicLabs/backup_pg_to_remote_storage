package config

import (
	"testing"
)

//TODO do validation unittests
//TODO validate that StorageConfig has at least one filled object!
// Test Databases is empty

func TestValidConfigWithGlobalStorageConfig(t *testing.T) {
	var dbCfg []DatabaseConfig
	dbCfg = append(dbCfg, validDbConfigWithoutStorage0)
	dbCfg = append(dbCfg, validDbConfigWithoutStorage1)

	cfg := Config{validStorageConfig, dbCfg}

	err := ValidateConfig(cfg)
	if err != nil {
		t.Fatalf("Validation should not fail as config is valid")
	}
}

func TestValidConfigWithDatabaseStorageConfig(t *testing.T) {
	var dbCfg []DatabaseConfig
	dbCfg = append(dbCfg, validDbConfigWithoutStorage0)
	dbCfg = append(dbCfg, validDbConfigWithoutStorage1)

	cfg := Config{StorageConfig{}, dbCfg}

	err := ValidateConfig(cfg)
	if err != nil {
		t.Fatalf("Validation should not fail as config is valid")
	}
}

func TestInvalidGlobalStorageConfig(t *testing.T) {
	var dbCfg []DatabaseConfig
	dbCfg = append(dbCfg, validDbConfigWithoutStorage0)

	cfg := Config{invalidStorageConfig, dbCfg}

	err := ValidateConfig(cfg)
	if err == nil {
		t.Fatalf("Validation should fail, storage config is invalid")
	}
}

func TestInvalidDatabaseConfig(t *testing.T) {
	var dbCfg []DatabaseConfig
	dbCfg = append(dbCfg, invalidDbConfigWithoutStorage0)

	cfg := Config{validStorageConfig, dbCfg}

	err := ValidateConfig(cfg)
	if err == nil {
		t.Fatalf("Validation should fail, database config is invalid")
	}
}

func TestNoConfigIsGiven(t *testing.T) {
	cfg := Config{StorageConfig{}, []DatabaseConfig{}}
	err := ValidateConfig(cfg)
	if err == nil {
		t.Fatalf("Validation should fail, as no config is given")
	}
}

func TestInvalidDatabaseConfigAndValidStorageConfig(t *testing.T) {
	var dbCfg []DatabaseConfig
	dbCfg = append(dbCfg, invalidDbConfigWithoutStorage0)

	cfg := Config{validStorageConfig, dbCfg}
	err := ValidateConfig(cfg)
	if err == nil {
		t.Fatalf("Validation should fail, as database config is invalid")
	}
}

func TestInvalidGlobalStorageAndNoStorageConfig(t *testing.T) {
	var dbCfg []DatabaseConfig
	dbCfg = append(dbCfg, validDbConfigWithoutStorage0)

	cfg := Config{invalidStorageConfig, dbCfg}
	err := ValidateConfig(cfg)
	if err == nil {
		t.Fatalf("Validation should fail, as no valid storage config and invalid global storage config")
	}
}

func TestNoGlobalStorageConfigAndNoDbStorageConfig(t *testing.T) {
	var dbCfg []DatabaseConfig
	dbCfg = append(dbCfg, validDbConfigWithoutStorage0)

	cfg := Config{StorageConfig{}, dbCfg}
	err := ValidateConfig(cfg)
	if err == nil {
		t.Fatalf("Validation should fail, as no StorageConfig was provied")
	}
}

func TestNoGlobalStorageConfigAndNoS3ConfigOnDb(t *testing.T) {
	var dbCfg []DatabaseConfig
	validDbConfigWithoutStorage0.StorageConfig = invalidStorageConfigNoConfig
	dbCfg = append(dbCfg, validDbConfigWithoutStorage0)

	cfg := Config{StorageConfig{}, dbCfg}
	err := ValidateConfig(cfg)
	if err == nil {
		t.Fatalf("Validation should fail, as no StorageConfig was provied")
	}
}

var validS3Config = S3Config{
	Endpoint:  "localhost",
	AccessKey: "accessKey",
	SecretKey: "secretKey",
	Bucket:    "bucket",
}

var validStorageConfig = StorageConfig{
	S3Config: validS3Config,
}

var invalidStorageConfigNoConfig = StorageConfig{}

var invalidS3Config = S3Config{
	AccessKey: "accessKey",
	Bucket:    "bucket",
}

var invalidStorageConfig = StorageConfig{
	S3Config: invalidS3Config,
}

var invalidDbConfigWithoutStorage0 = DatabaseConfig{
	Host:      "xxx",
	Port:      5432,
	Password:  "password",
	Retention: 30,
}

var validDbConfigWithoutStorage0 = DatabaseConfig{
	Host:      "postgres.example.com",
	Port:      5432,
	Database:  "Database",
	Username:  "username",
	Password:  "password",
	Retention: 30,
}

var validDbConfigWithoutStorage1 = DatabaseConfig{
	Host:      "127.0.0.1",
	Port:      5432,
	Database:  "Database",
	Username:  "username",
	Password:  "password",
	Retention: 30,
}

var validDbConfigWithStorage0 = DatabaseConfig{
	Host:          "postgres.example.com",
	Port:          5432,
	Database:      "Database",
	Username:      "username",
	Password:      "password",
	Retention:     30,
	StorageConfig: validStorageConfig,
}

var validDbConfigWithStorage1 = DatabaseConfig{
	Host:          "127.0.0.1",
	Port:          5432,
	Database:      "Database",
	Username:      "username",
	Password:      "password",
	Retention:     30,
	StorageConfig: validStorageConfig,
}
