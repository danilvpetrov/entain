module github.com/danilvpetrov/entain

go 1.25.0

tool (
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	google.golang.org/genproto/googleapis/api
	google.golang.org/grpc/cmd/protoc-gen-go-grpc
	google.golang.org/protobuf/cmd/protoc-gen-go
)

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.2
	github.com/mattn/go-sqlite3 v1.14.32
	google.golang.org/genproto/googleapis/api v0.0.0-20250922171735-9219d122eba9
	google.golang.org/grpc v1.75.1
	google.golang.org/protobuf v1.36.9
	syreclabs.com/go/faker v1.2.3
)

require (
	github.com/kr/text v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250908214217-97024824d090 // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.5.1 // indirect
)
