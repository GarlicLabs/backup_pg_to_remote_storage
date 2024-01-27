package config

import (
	"testing"
)

func ShouldMapGlobalToEmptyDbStorage(t *testing.T) {
	cfg := MapGlobalStorageToDbIfNotSet(validDbConfigWithoutStorage0, validStorageConfig)
	if cfg.StorageConfig.S3Config.AccessKey != validStorageConfig.S3Config.AccessKey {
		t.Fatalf("Unit test failed case access key is different")
	}
	if cfg.StorageConfig.S3Config.SecretKey != validStorageConfig.S3Config.SecretKey {
		t.Fatalf("Unit test failed case secret key is different")
	}

	if cfg.StorageConfig.S3Config.Bucket != validStorageConfig.S3Config.Bucket {
		t.Fatalf("Unit test failed case bucket is different")
	}

	if cfg.StorageConfig.S3Config.Endpoint != validStorageConfig.S3Config.Endpoint {
		t.Fatalf("Unit test failed case endpoint is different")
	}
}

func ShouldNotMapGlobalBecauseDbHasStorage(t *testing.T) {
	cfg := MapGlobalStorageToDbIfNotSet(validDbConfigWithStorage0, validStorageConfig)
	if cfg.StorageConfig.S3Config.AccessKey == validStorageConfig.S3Config.AccessKey {
		t.Fatalf("Unit test failed case access key is the same")
	}
	if cfg.StorageConfig.S3Config.SecretKey == validStorageConfig.S3Config.SecretKey {
		t.Fatalf("Unit test failed case secret key is the same")
	}

	if cfg.StorageConfig.S3Config.Bucket == validStorageConfig.S3Config.Bucket {
		t.Fatalf("Unit test failed case bucket is the same")
	}

	if cfg.StorageConfig.S3Config.Endpoint == validStorageConfig.S3Config.Endpoint {
		t.Fatalf("Unit test failed case endpoint is the same")
	}
}

var diffrentValidS3Config = S3{
	Endpoint:  "diffrentLocalhost",
	AccessKey: "diffrentAccessKey",
	SecretKey: "diffrentSecretKey",
	Bucket:    "diffrentBucket",
}

var diffrentValidStorageConfig = Storage{
	S3Config: validS3Config,
}
