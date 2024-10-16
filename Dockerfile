FROM golang:1.23.1
RUN apt-get update && apt-get install -y build-essential make default-mysql-client

# go のパッケージを install
RUN go install github.com/rubenv/sql-migrate/...@v1.5.2

WORKDIR /server
COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify
