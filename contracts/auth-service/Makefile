PROTOC_GEN_GO       := $(shell go env GOPATH)/bin/protoc-gen-go.exe
PROTOC_GEN_GO_GRPC  := $(shell go env GOPATH)/bin/protoc-gen-go-grpc.exe

PROTO_FILES := auth.proto auth.proto

.PHONY: proto
proto: 
	protoc \
		--plugin=protoc-gen-go=$(PROTOC_GEN_GO) \
		--plugin=protoc-gen-go-grpc=$(PROTOC_GEN_GO_GRPC) \
		--proto_path=proto \
		--go_out=gen \
		--go_opt=paths=source_relative \
		--go-grpc_out=gen \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
