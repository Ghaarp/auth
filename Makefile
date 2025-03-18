LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


generate:
	make generate-note-api

generate-note-api:
	mkdir -p pkg/auth_v1
	protoc --proto_path api/auth_v1 \
	--go_out=pkg/auth_v1 --go_opt=paths=source_relative \
	--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
	api/auth_v1/auth.proto

build:
	GOOS=linux GOARCH=amd64 go build -o auth_service cmd/server/main.go

copy-to-server:
	scp -v auth_service root@31.128.50.199:~

create-migration:
	goose create migrations/name sql

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t cr.selcloud.ru/f32f3g423w23efg32/auth-server:v0.0.1 .
	docker login -u token -p CRgAAAAAKcaswcPkDDMcFxAShN9RJ-1I5cV1GzFN cr.selcloud.ru/f32f3g423w23efg32
	docker push cr.selcloud.ru/f32f3g423w23efg32/auth-server:v0.0.1

run-local:
	go run cmd/server/main.go -config-path="cloud.env"

run:
	docker-compose --env-file docker.env up --build

docker-run-local-noenv:
	docker-compose up --build
