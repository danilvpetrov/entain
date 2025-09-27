
PROTO_FILES += $(shell find api -type f -name '*.proto')

.PHONY: test
test:
	go test -v -race -count 1 ./...

.PHONY: generate
generate: $(PROTO_FILES)
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		--go-grpc_opt=require_unimplemented_servers=false \
		--grpc-gateway_out=. \
		--grpc-gateway_opt=paths=source_relative \
		--experimental_allow_proto3_optional \
		--openapiv2_out . \
		--openapiv2_opt=logtostderr=true,output_format=yaml \
		--proto_path=. \
		--proto_path=google/api=$(shell pwd)/proto/google/api \
		 $(PROTO_FILES)


.PHONY: run-gateway
run-gateway:
	go run ./cmd/gateway


.PHONY: run-racing
run-racing:
	go run ./cmd/racing


.PHONY: precommit
precommit: generate test
	go mod tidy
	betteralign --apply -test_files ./...
	golangci-lint run --fix ./...
