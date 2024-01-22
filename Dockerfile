FROM debian:stable

WORKDIR /app

COPY . .

ENTRYPOINT ["./backup_pg_to_remote_storage"]
