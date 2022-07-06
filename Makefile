export GO111MODULE=on
export GOPROXY=direct
export GOSUMDB=off
export CGO_ENABLED=0

-include .env

MIGRATIONS_DIR = ./sql/migrations/

start:
	APP_TYPE=server go run -v ./cmd/start

service:
	APP_TYPE=service go run -v ./cmd/start

build:
	go build -o ./build/server -v cmd/server/server.go

upgrade:
	GOWORK=off go-mod-upgrade
	go mod tidy

compose-start:
	docker-compose up -d

compose-stop:
	docker-compose stop

docker-build:
	docker build . -t sxwebdev/go-test-app -f Dockerfile
	docker image prune -f

docker-start:
	docker run --name sxwebdev/go-test-app -it -d  -p 5432:5432 -p 4222:4222 -p 1883:1883 -p 8080:8080 -p 9000:9000 --rm sxwebdev/go-test-app

gen-protos:
	rm -rf ./pb/*
	protoc \
	--go_out=:pb \
	--go-grpc_out=:pb \
	proto/*.proto

migrate:
	migrate -path "$(MIGRATIONS_DIR)" -database "$(DB_DSN)" $(filter-out $@,$(MAKECMDGOALS))

create-migration:
	migrate create -ext sql -dir "$(MIGRATIONS_DIR)" $(filter-out $@,$(MAKECMDGOALS))

%:
	@:
