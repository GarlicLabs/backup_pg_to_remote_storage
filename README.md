# backup-pg-to-remote-storage

----

* Config Unittests
* S3
* Integration tests
* Deployment
* Write README

----

// TODO: Put in readme/docs, that the global s3 config has to be set in any case, but can be overridden for any database. You can just override the bucket and leave the default endpoint and keys as is, but if these need to be changed, you have to set all the values on the database if needed.


## Description

backup-pg-to-remote-storage is a small programm to dumo postgres databases and store them into a s3 storage.  

## Usage

The container image is weekly updated with the stable debian image.
If you are using the container image make sure to pull regulary, the tag is always `latest`.


Usage with docker:  

```bash
#TODO define config & mount volume with config into container
docker run ghcr.io/garliclabs/backup-pg-to-remote-storage:latest
```

Just use the program:  

```bash
#TODO define config
go build
./backup-pg-to-remote-storage
```

## Development

### Testing

You need docker compose

First start integration testing infrastructure:  

```bash
cd test/
docker-compose up -d
```


