package main

import (
	"os"

	"github.com/garliclabs/backup-pg-to-remote-storage/cmd/config"
	"github.com/garliclabs/backup-pg-to-remote-storage/cmd/storage"
	pg "github.com/habx/pg-commands"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type RemoteStorage interface {
	Upload(cfg config.Storage, file string)
	Delete(cfg config.Storage, file string)
	RetentionDelete(cfg config.Storage)
}

func main() {

	cfg := getConfig()
	err := config.Validate(cfg)
	if err != nil {
		log.Panicf("Configuration is invalid see error: %s", err.Error())
	}

	for _, database := range cfg.Databases {
		dbCfg := config.MapGlobalStorageToDbIfNotSet(database, cfg.GlobalStorageConfig)
		dumpFile := dumpDatabase(dbCfg)
		getStorageProvider(dbCfg.StorageConfig).Upload(dbCfg.StorageConfig, dumpFile)
		getStorageProvider(dbCfg.StorageConfig).RetentionDelete(dbCfg.StorageConfig)
		removeDumpFromFilesystem(dumpFile)
	}
}

func getStorageProvider(cfg config.Storage) RemoteStorage {
	if cfg.S3Config != (config.S3{}) {
		return storage.S3Storage{}
	} else {
		log.Panicf("No storage provider configured")
		return nil
	}
}

func getConfig() config.Config {
	log.Info("Reading config file")
	f, err := os.ReadFile("config.yml")

	if err != nil {
		log.Error(err)
	}

	var config config.Config
	err = yaml.Unmarshal(f, &config)

	if err != nil {
		log.Error(err)
	}

	return config
}

func dumpDatabase(databaseConfig config.Database) string {
	log.Infof("Starting dumping Database %s at %s", databaseConfig.Database, databaseConfig.Host)
	dumper, err := pg.NewDump(&pg.Postgres{
		Host:     databaseConfig.Host,
		Port:     databaseConfig.Port,
		DB:       databaseConfig.Database,
		Username: databaseConfig.Username,
		Password: databaseConfig.Password,
	})
	if err != nil {
		log.Error(err)
	}

	dump := dumper.Exec(pg.ExecOptions{StreamPrint: false})
	if dump.Error != nil {
		log.Error(dump.Error.Err)
		log.Error(dump.Output)

	} else {
		log.Infof("Dumping Database %s at %s success", databaseConfig.Database, databaseConfig.Host)
	}
	return dump.File
}

func removeDumpFromFilesystem(File string) {
	err := os.Remove(File)
	if err != nil {
		log.Errorln("Removing Dump-File failed: ", err)
	}
}