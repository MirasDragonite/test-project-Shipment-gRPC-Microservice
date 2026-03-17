.PHONY: proto
proto:
	mkdir -p api/proto/gen
	protoc --proto_path=api/proto \
           --go_out=api/proto/gen --go_opt=paths=source_relative \
           --go-grpc_out=api/proto/gen --go-grpc_opt=paths=source_relative \
           api/proto/shipment.proto