.PHONY: gen
gen:
	protoc --proto_path=internal/controller/grpcapi/proto/blacklist internal/controller/grpcapi/proto/blacklist/*.proto  --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import
	protoc --proto_path=internal/controller/grpcapi/proto/whitelist internal/controller/grpcapi/proto/whitelist/*.proto  --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import
	protoc --proto_path=internal/controller/grpcapi/proto/bucket internal/controller/grpcapi/proto/bucket/*.proto  --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import
	protoc --proto_path=internal/controller/grpcapi/proto/authorization internal/controller/grpcapi/proto/authorization/*.proto  --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import

.PHONY: mok_gen
mok_gen:
	echo "empty"

.PHONY:clean
clean:
	rm -f internal/controller/grpcapi/blacklistpb/*
	rm -f internal/controller/grpcapi/whitelistpb/*
	rm -f internal/controller/grpcapi/authorizationpb/*
	rm -f internal/controller/grpcapi/bucketpb/*



