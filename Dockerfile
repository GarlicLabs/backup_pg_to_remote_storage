FROM debian:stable

WORKDIR /app

RUN apt-get update && apt-get install postgresql-client -y

COPY . .

ENTRYPOINT ["./backup_to_remote_storage"]
