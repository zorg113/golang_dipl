MODULE := github.com/zorg113/golang_dipl/atibruteforce
PROTO_DIR := ./internal/controller/grpcapi/proto

.PHONY: gen
gen:
	protoc \
			--proto_path=$(PROTO_DIR) \
			--go_out=. \
			--go_opt=module=$(MODULE) \
			--go-grpc_out=. \
			--go-grpc_opt=module=$(MODULE) \
			$(PROTO_DIR)/common/common.proto \
			$(PROTO_DIR)/authorization/authorization_service.proto \
			$(PROTO_DIR)/authorization/authorization_message.proto \
			$(PROTO_DIR)/bucket/bucket_message.proto \
			$(PROTO_DIR)/bucket/bucket_service.proto \
			$(PROTO_DIR)/blacklist/blacklist_message.proto \
			$(PROTO_DIR)/blacklist/blacklist_service.proto \
			$(PROTO_DIR)/whitelist/whitelist_message.proto \
			$(PROTO_DIR)/whitelist/whitelist_service.proto 
	
.PHONY: mock_gen
mock_gen:
	mockgen -source=model/service/blacklist.go -destination=store/adapters/mocks/mock_blacklist.go
	mockgen -source=model/service/whitelist.go -destination=store/adapters/mocks/mock_whitelist.go

.PHONY:clean
clean:
	rm -f internal/controller/grpcapi/blacklistpb/*
	rm -f internal/controller/grpcapi/whitelistpb/*
	rm -f internal/controller/grpcapi/authorizationpb/*
	rm -f internal/controller/grpcapi/bucketpb/*
	rm -f internal/controller/grpcapi/commonpb/*

.PHONY: build.bin

build.bin:
	go build -o ./build/antibriteforce_service ./cmd/service

.PHONY: lint
lint:
	golangci-lint run

.PHONY: build.docker
build.docker:
	docker build --tag  antibrf -- .

.PHONY: build
build:
	docker compose -f docker-compose.yml build

.PHONY: run
run: build
	docker compose -f docker-compose.yml up

.PHONY: stop
stop: 
	docker compose -f docker-compose.yml down

.PHONY: integration.test.run
integration.test.run:
	docker compose -f docker-compose-test.yml build
	docker compose -f docker-compose-test.yml up

.PHONY: integration.test.stop	
integration.test.stop:
	docker compose -f docker-compose-test.yml down

.PHONY:migrate
migrate:
	migrate -version $(version)

.PHONY: migrate.down
migrate.down:
	migrate -source file://migrations -database postgres://localhost:5433/antibruteforce-service-database?sslmode=disable down

.PHONY: migrate.up
migrate.up:
	migrate -source file://migrations -database postgres://localhost:5433/antibruteforce-service-database?sslmode=disable up

.PHONY: test
test:
	go test -race -count 100 ./...