PROTOC_GEN_GO := $(go env GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(go env GOPATH)/bin/protoc-gen-go-grpc

all: protos/keyvalue_grpc.pb.go protos/keyvalue.pb.go

$(PROTOC_GEN_GO):
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

$(PROTOC_GEN_GO_GRPC):
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

protos/keyvalue_grpc.pb.go protos/keyvalue.pb.go : protos/keyvalue.proto | $(PROTOC_GEN_GO) $(PROTOC_GEN_GO_GRPC)
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	$<

.PHONY: clean
clean :
	-rm protos/*.pb.go
