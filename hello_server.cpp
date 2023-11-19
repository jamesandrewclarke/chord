#include <iostream>
#include <memory>
#include <string>

#include <grpcpp/ext/proto_server_reflection_plugin.h>
#include <grpcpp/grpcpp.h>
#include <grpcpp/health_check_service_interface.h>

#include "hello.grpc.pb.h"

using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::Status;

using Hello::Greeter;
using Hello::HelloRequest;
using Hello::HelloResponse;

class GreeterServiceImpl final : public Greeter::Service {
    Status Hello(ServerContext* context, const HelloRequest* request,
        HelloResponse* reply) override {
        std::string prefix ("Hello ");
        reply->set_response(prefix + request->name());
        return Status::OK;
    }

    Status SetValue(ServerContext* context, const Hello::SetKeyRequest* request,
        Hello::KeyResponse* response) override {
        const auto& key = request->key();
        const auto& value = request->value();
        if (!key.empty() && !value.empty())
            keys[key] = value;
        else
            return Status::CANCELLED;

        std::cout << "set : " << key << ": " << value << std::endl;
        return Status::OK;
    }

    Status GetValue(grpc::ServerContext* context, const Hello::KeyRequest* request,
        Hello::KeyResponse* response) override {
        if (const auto& key = request->key(); keys.contains(key)) {
            response->set_value(keys[key]);
        } else
            return Status::CANCELLED;

        return Status::OK;
    }

private:
    std::unordered_map<std::string, std::string> keys;
};

void RunServer(uint16_t port) {
    const std::string server_address = "0.0.0.0:50051";
    GreeterServiceImpl service;

    grpc::EnableDefaultHealthCheckService(true);
    grpc::reflection::InitProtoReflectionServerBuilderPlugin();

    ServerBuilder builder;

    builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());

    builder.RegisterService(&service);

    const std::unique_ptr<Server> server(builder.BuildAndStart());
    std::cout << "Server listening on " << server_address << std::endl;

    server->Wait();
}

int main(int argc, char** argv) {
    RunServer(8080);
}