package storage

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/garliclabs/backup-pg-to-remote-storage/cmd/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
)

type S3Storage struct {
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
	Bucket    string `yaml:"bucket"`
}

func (s3 S3Storage) initClient(s3Config config.S3) (*minio.Client, error) {
	log.Info("Initializing the S3 Client.")
	useSSL := true

	creds := credentials.NewStaticV4(s3Config.AccessKey, s3Config.SecretKey, "")

	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = true

	var transport http.RoundTripper = &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	options := &minio.Options{
		Creds:     creds,
		Secure:    useSSL,
		Transport: transport,
	}
	minioClient, err := minio.New(s3Config.Endpoint, options)
	if err != nil {
		log.Errorf("Connection to S3 Storage at %s failed: %s", s3Config.Endpoint, err)
		return nil, err
	} else {
		log.Infof("Connection to S3 Storage at %s succeeded.", s3Config.Endpoint)
	}
	return minioClient, nil
}

func (s3 S3Storage) Upload(storageCfg config.Storage, file string) error {
	cfg := storageCfg.S3Config
	ctx := context.Background()

	minioClient, err := s3.initClient(cfg)

	if err != nil {
		return err
	}

	contentType := "application/octet-stream"
	options := minio.PutObjectOptions{ContentType: contentType}

	log.Infof("Uploading file %s to %s.", file, minioClient.EndpointURL())
	_, err = minioClient.FPutObject(ctx, cfg.Bucket, file, file, options)
	if err != nil {
		log.Errorf("Uploading file %s to %s failed.", file, minioClient.EndpointURL())
		log.Error(err)
		return err
	} else {
		log.Infof("Successfully uploaded %s to %s", file, minioClient.EndpointURL())
	}
	return nil
}

func (s3 S3Storage) Delete(storageCfg config.Storage, file string) error {
	cfg := storageCfg.S3Config
	ctx := context.Background()

	minioClient, err := s3.initClient(cfg)
	if err != nil {
		return err
	}
	err = minioClient.RemoveObject(ctx, cfg.Bucket, file, minio.RemoveObjectOptions{})
	if err != nil {
		log.Errorf("Removing object %s failed.", file)
		log.Error(err)
		return err
	}
	return nil
}

func (s3 S3Storage) RetentionDelete(dbConfig config.Database) error {
	cfg := dbConfig.StorageConfig.S3Config
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	minioClient, err := s3.initClient(cfg)
	if err != nil {
		return err
	}

	log.Infof("Getting objects of Storage at %s in bucket %s.", minioClient.EndpointURL(), cfg.Bucket)
	objects := minioClient.ListObjects(ctx, cfg.Bucket, minio.ListObjectsOptions{
		Prefix:    dbConfig.Database,
		Recursive: true,
	})

	// Remove Backups that are older than the given Retention
	hoursInDay := 24
	retentionInHours := dbConfig.Retention * hoursInDay

	log.Infof("Removing objects older than %d days in Bucket %s at S3 Storage %s", dbConfig.Retention, cfg.Bucket, minioClient.EndpointURL())
	for object := range objects {
		age := time.Since(object.LastModified)
		ageInHours := age.Hours()
		is_age_bigger_than_retention := ageInHours > float64(retentionInHours)

		if is_age_bigger_than_retention {
			err := s3.Delete(dbConfig.StorageConfig, object.Key)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
