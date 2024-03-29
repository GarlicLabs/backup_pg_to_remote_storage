package main

import (
	"os"

	"github.com/garliclabs/backup-pg-to-remote-storage/cmd/config"
	"github.com/garliclabs/backup-pg-to-remote-storage/cmd/storage"
	pg "github.com/habx/pg-commands"
	log "github.com/sirupsen/logrus"
)

type RemoteStorage interface {
	Upload(cfg config.Storage, file string) error
	Delete(cfg config.Storage, file string) error
	RetentionDelete(dbConfig config.Database) error
}

func main() {

	cfg := config.GetConfig()
	err := config.Validate(cfg)
	if err != nil {
		log.Errorf("Configuration is invalid see error: %s", err.Error())
		os.Exit(1)
	}

	var dumpFailed bool
	for _, database := range cfg.Databases {
		dbCfg := config.MapGlobalStorageToDbIfNotSet(database, cfg.GlobalStorageConfig)
		dumpFile, err := dumpDatabase(dbCfg)
		if err != nil {
			dumpFailed = errorHandling(err, dumpFile)
			continue
		}
		err = getStorageProvider(dbCfg.StorageConfig).Upload(dbCfg.StorageConfig, dumpFile)
		if err != nil {
			dumpFailed = errorHandling(err, dumpFile)
			continue
		}
		err = getStorageProvider(dbCfg.StorageConfig).RetentionDelete(dbCfg)
		if err != nil {
			dumpFailed = errorHandling(err, dumpFile)
			continue
		}
		removeDumpFromFilesystem(dumpFile)
	}
	if dumpFailed {
		log.Error("At least one dump failed, please check logs for more information")
		os.Exit(1)
	} else {
		log.Info("Finished backuping all databases")
	}
}

func errorHandling(err error, dumpFile string) bool {
	removeDumpFromFilesystem(dumpFile)
	log.Error(err)
	return true
}

func getStorageProvider(cfg config.Storage) RemoteStorage {
	if cfg.S3Config != (config.S3{}) {
		return storage.S3Storage{}
	} else {
		log.Error("No storage provider configured")
		os.Exit(1)
		return nil
	}
}

func dumpDatabase(databaseConfig config.Database) (string, error) {
	log.Infof("Starting dumping Database %s at %s", databaseConfig.Database, databaseConfig.Host)
	dumper, err := pg.NewDump(&pg.Postgres{
		Host:     databaseConfig.Host,
		Port:     databaseConfig.Port,
		DB:       databaseConfig.Database,
		Username: databaseConfig.Username,
		Password: databaseConfig.Password,
	})
	if err != nil {
		return "", err
	}

	dump := dumper.Exec(pg.ExecOptions{StreamPrint: false})
	if dump.Error != nil {
		log.Error(dump.Output)
		return dump.File, dump.Error.Err
	} else {
		log.Infof("Dumping Database %s at %s success", databaseConfig.Database, databaseConfig.Host)
	}
	return dump.File, nil
}

func removeDumpFromFilesystem(File string) {
	log.Infof("Removing File %s", File)
	err := os.Remove(File)
	if err != nil {
		log.Errorln("Removing Dump-File failed: ", err)
	}
}
