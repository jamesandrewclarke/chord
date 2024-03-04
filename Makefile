PROTOC_GEN_GO := $(go env GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(go env GOPATH)/bin/protoc-gen-go-grpc

all: protos/chord.pb.go

$(PROTOC_GEN_GO):
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

$(PROTOC_GEN_GO_GRPC):
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

protos/chord_grpc.pb.go protos/chord.pb.go : protos/chord.proto | $(PROTOC_GEN_GO) $(PROTOC_GEN_GO_GRPC)
	protoc \
		--experimental_allow_proto3_optional \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	$<

.PHONY: clean
clean :
	-rm protos/*.pb.go
