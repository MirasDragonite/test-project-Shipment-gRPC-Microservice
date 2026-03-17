.PHONY: proto up down logs restart rebuild
proto:
	mkdir -p api/proto/gen
	protoc --proto_path=api/proto \
           --go_out=api/proto/gen --go_opt=paths=source_relative \
           --go-grpc_out=api/proto/gen --go-grpc_opt=paths=source_relative \
           api/proto/shipment.proto

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f app

rebuild:
	docker-compose up -d --build app

restart: down up
