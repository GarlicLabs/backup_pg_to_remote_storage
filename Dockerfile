FROM debian:stable

WORKDIR /app

COPY . /app

ENTRYPOINT ["./backup-pg-to-remote-storage"]
