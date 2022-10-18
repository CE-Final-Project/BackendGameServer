
proto_kafka_messages:
	@echo Generating kafka messages proto
	cd proto/kafka_messages && protoc --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. kafka.proto

# proto auth
proto_auth:
	@echo Generating auth microservice proto
	cd authentication/proto && protoc --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. auth_service.proto \
	&& protoc --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. auth_message.proto \
	&& protoc --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. account_message.proto \
	&& protoc --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. role_message.proto


swagger_gateway:
	@echo Starting swagger generating
	cd gateway && swag init -g ./cmd/main.go


gen_rsa:
	@echo Generating RSA...
	cd rsa && openssl genrsa -out jwt 4096 && openssl rsa -in jwt -pubout -out jwt.pub