.PHONY:
.SILENT:

build:
    go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/segmenter ./cmd/segmenter/main.go 

run: build
    docker compose up segmenter

test:
    go test -v ./...

swag:
    swag init -g internal/app/app.go