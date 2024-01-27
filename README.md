# backup_pg_to_remote_storage

backup_pg_to_remote_storage is a small programm to dump postgres databases and store them into a remote storage.  
! Currently just s3 is supported as a remote storage, extending the program to support more remote storage can be easily implemented.  

## Usage

The system on that you want to execute the program needs to have a postgres client installed, the container image ofcourse has it already installed.  

The container image is weekly updated with the stable debian image.
If you are using the container image make sure to pull regulary, the tag is always `latest`.

The program will iterate over all configured databases, if one backup should fail it will be continued with the next database. The program will then end with the exit code 1.  

You are able to set a global storage config or a storage config on each database you want to backup. If you have a global storage config and a storage config on a database, the config on the database will be used.  

You can specify the config location with these enviroment variable: `BACKUP_PG_CONFIG_PATH`.  

### Usage with docker

```bash
docker run ghcr.io/garliclabs/backup-pg-to-remote-storage:latest
```

### Usage with binary

```bash
go build
BACKUP_PG_CONFIG_PATH=./config.yml ./backup-pg-to-remote-storage
```

### Usage on k8s

You can find a bunch of kubernetes manifest files in the `./k8s/` folder, create your own secret file and apply everything.  

### Kubernetes

We provide some kubernetes manifest files (see `./k8s/`), for an easier kubernetes deployment.  

## Development

### Pre requriments

* docker & docker-compose
* golang > 1.20
* For convinece Make
* Postgres client

### Download dependecies and build

1. You need to download all dependecies

```bash
#If you have make
make update
#Else
go mod download
```

2. Build application

```bash
#If you have make
make build
#Else
go build -o backup_to_remote_storage cmd/main.go
```

### Setup test enviroment

A local development enviroment can be created with the `docker-compose.yml` file in the `./test/` folder.  
A minio server will be started with a bucket called `test` and a postgres database will be started populated with a few dummy data (can be found at `./test/init.sql`).  

You need a config file, as template you can use `./secret.example.yml` the values inside this configuration are already valid for the testing enviroment started with docker-compose, so you are ready to start.  

A makefile is provided with a few convince mappings around needed commands.  

### Testing

Run unit tests:  

```bash
#If you have make
make unittests
#Else
go test ./cmd/...
```
