package main

import (
	"context"
	"time"

	"github.com/garliclabs/backup-pg-to-remote-storage/cmd/config"
	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/credentials"
	log "github.com/sirupsen/logrus"
)

type S3Storage struct {
}

func initializeS3Client(S3Config config.S3Config) *minio.Client {
	log.Info("Initializing the S3 Client.")
	useSSL := true
	options := &minio.Options{
		Creds:  credentials.NewStaticV4(S3Config.AccessKey, S3Config.SecretKey, ""),
		Secure: useSSL,
	}
	// Initialize minio client object.
	minioClient, err := minio.New(S3Config.Endpoint, options)
	if err != nil {
		log.Errorf("Connection to S3 Storage at %s failed.", S3Config.Endpoint)
		log.Error(err)
	} else {
		log.Infof("Connection to S3 Storage at %s succeeded.", S3Config.Endpoint)
	}
	return minioClient
}

func (s3 S3Storage) Upload(storageCfg config.StorageConfig, file string) {
	cfg := storageCfg.S3Config

}
func (s3 S3Storage) Delete(cfg config.StorageConfig, file string) {
	cfg := storageCfg.S3Config

}
func (s3 S3Storage) RetentionDelete(cfg config.StorageConfig) {
	cfg := storageCfg.S3Config

}

//////////////////////////////////////////////////////

func uploadtoS3(minioClient *minio.Client, databaseConfig config.DatabaseConfig, localFilePath string) {
	object := databaseConfig.Database + "-backup-" + time.Now().Format("01-02-2006") + ".sql.tar.gz"
	ctx := context.Background()
	contentType := "application/octet-stream"
	options := minio.PutObjectOptions{ContentType: contentType}

	log.Infof("Uploading file %s to %s.", object, minioClient.EndpointURL())
	_, err := minioClient.FPutObject(ctx, databaseConfig.StorageConfig.Bucket, object, localFilePath, options)
	if err != nil {
		log.Errorf("Uploading file %s to %s failed.", object, minioClient.EndpointURL())
		log.Error(err)
	}
	log.Infof("Successfully uploaded %s to %s", object, minioClient.EndpointURL())
}

func removeBackupsOlderThanRetention(minioClient *minio.Client, databaseConfig config.DatabaseConfig) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Infof("Getting objects of Storage at %s in bucket %s.", minioClient.EndpointURL(), databaseConfig.StorageConfig.Bucket)
	objects := minioClient.ListObjects(ctx, databaseConfig.StorageConfig.Bucket, minio.ListObjectsOptions{
		Prefix:    databaseConfig.Database,
		Recursive: true,
	})

	// Remove Backups that are older than the given Retention
	hoursInDay := 24
	retentionInHours := databaseConfig.Retention * hoursInDay

	log.Infof("Removing objects older than %d days in Bucket %s at S3 Storage %s", databaseConfig.Retention, databaseConfig.StorageConfig.Bucket, minioClient.EndpointURL())
	for object := range objects {
		age := time.Since(object.LastModified)
		ageInHours := age.Hours()
		is_age_bigger_than_retention := ageInHours > float64(retentionInHours)
		if is_age_bigger_than_retention {
			err := minioClient.RemoveObject(ctx, databaseConfig.StorageConfig.Bucket, object.Key, minio.RemoveObjectOptions{})
			if err != nil {
				log.Errorf("Removing object %s failed.", object.Key)
				log.Error(err)
			}
		}
	}
}
