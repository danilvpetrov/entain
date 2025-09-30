
PROTO_FILES += $(shell find api -type f -name '*.proto')
PROTO_GENERATED_FILES += $(foreach f,$(PROTO_FILES:.proto=.pb.go),$f)
PROTO_GENERATED_FILES += $(foreach f,$(PROTO_FILES:.proto=_grpc.pb.go),$f)
PROTO_GENERATED_FILES += $(foreach f,$(PROTO_FILES:.proto=.pb.gw.go),$f)
PROTO_GENERATED_FILES += $(foreach f,$(PROTO_FILES:.proto=.swagger.yaml),$f)

%.pb.go: %.proto
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--proto_path=. \
		--proto_path=google/api=$(shell pwd)/proto/google/api \
		$(@D)/*.proto

%_grpc.pb.go: %.proto
	protoc \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		--go-grpc_opt=require_unimplemented_servers=false \
		--proto_path=. \
		--proto_path=google/api=$(shell pwd)/proto/google/api \
		$(@D)/*.proto

%.pb.gw.go: %.proto
	protoc \
		--grpc-gateway_out=. \
		--grpc-gateway_opt=paths=source_relative \
		--proto_path=. \
		--proto_path=google/api=$(shell pwd)/proto/google/api \
		--experimental_allow_proto3_optional \
		$(@D)/*.proto

%.swagger.yaml: %.proto
	protoc \
		--openapiv2_out . \
		--openapiv2_opt=logtostderr=true,output_format=yaml \
		--proto_path=. \
		--proto_path=google/api=$(shell pwd)/proto/google/api \
		$(@D)/*.proto

.PHONY: test
test: $(PROTO_GENERATED_FILES)
	go test -v -race -count 1 ./...

.PHONY: generate
generate: $(PROTO_GENERATED_FILES)

.PHONY: precommit
precommit: test
	go mod tidy
	betteralign --apply -test_files ./...
	golangci-lint run --fix ./...

.PHONY: import-sports-events
import-sports-events:
	go run ./sports/testdata

.PHONY: run-gateway
run-gateway: artefacts/make/docker_jaeger.touch
	OTLP_TRACE_EXPORTER_ENDPOINT=localhost:4317 \
	go run ./cmd/gateway

.PHONY: run-racing
run-racing: artefacts/make/docker_jaeger.touch
	OTLP_TRACE_EXPORTER_ENDPOINT=localhost:4317 \
	RACING_DB_PATH="artefacts/db/racing.db" \
	go run ./cmd/racing

.PHONY: run-sports
run-sports: artefacts/make/docker_jaeger.touch
	OTLP_TRACE_EXPORTER_ENDPOINT=localhost:4317 \
	SPORTS_DB_PATH="artefacts/db/sports.db" \
	go run ./cmd/sports

artefacts/make/docker_jaeger.touch:
	@mkdir -p $(@D)
	docker run --name jaeger \
		--detach \
		--env COLLECTOR_OTLP_ENABLED=true \
		--publish "16686:16686/tcp" \
		--publish "4317:4317/tcp" \
		jaegertracing/all-in-one:latest
	@touch "$@"

.PHONY: clean
clean:
	rm -rf artefacts
	docker rm --force jaeger || true

