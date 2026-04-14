# Optional convenience targets. The project builds with plain `go` commands.

PROTOC_GEN_GO_VERSION ?= v1.34.2
PROTOC_GEN_GO_GRPC_VERSION ?= v1.5.1

# go install puts protoc-gen-go / protoc-gen-go-grpc into GOPATH/bin; protoc looks up plugins by PATH.
GOPATH_BIN := $(shell go env GOPATH)/bin

.PHONY: proto proto-tools test build run-monolith run-pricing run-order clean

proto-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)

proto: proto-tools
	mkdir -p pkg/pricingpb
	PATH="$(GOPATH_BIN):$$PATH" protoc -I proto \
		--go_out=pkg/pricingpb --go_opt=paths=source_relative \
		--go-grpc_out=pkg/pricingpb --go-grpc_opt=paths=source_relative \
		proto/pricing.proto

test:
	go test ./...

build: clean
	mkdir -p bin
	go build -o bin/monolith ./monolith
	go build -o bin/order-service ./services/order
	go build -o bin/pricing-service ./services/pricing

run-pricing:
	go run ./services/pricing

run-order:
	PRICING_GRPC_ADDR=localhost:50051 PORT=8081 go run ./services/order

run-monolith:
	PORT=8080 go run ./monolith

clean:
	rm -rf bin
