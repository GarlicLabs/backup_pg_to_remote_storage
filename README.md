# backup_pg_to_remote_storage

backup_pg_to_remote_storage is a small programm to dump postgres databases and store them into a remote storage.  
! Currently just s3 is supported as a remote storage, extending the program to support more remote storage can be easily implemented.  

## TODO Usage

!!The program will not fail if a database is not dumpable or if the storage has an issue it will continue with the next one.  

The container image is weekly updated with the stable debian image.
If you are using the container image make sure to pull regulary, the tag is always `latest`.

If you prefere to use the binary, you need to have a postgres-client install on your operating system.  

### Setup config

// TODO: Put in readme/docs, that the global s3 config has to be set in any case, but can be overridden for any database. You can just override the bucket and leave the default endpoint and keys as is, but if these need to be changed, you have to set all the values on the database if needed.

### Usage with docker

```bash
#TODO define config & mount volume with config into container
docker run ghcr.io/garliclabs/backup-pg-to-remote-storage:latest
```

### Just use the program

```bash
#The config.yml file inside the directory of the binary will be used, if there is no the process will end with an error
go build
./backup-pg-to-remote-storage
```

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
