//
// Created by james on 19/11/2023.
//

#include <iostream>
#include <memory>
#include <string>

#include <grpcpp/grpcpp.h>

#include "hello.grpc.pb.h"

using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;

using Hello::Greeter;
using Hello::HelloRequest;
using Hello::HelloResponse;

class GreeterClient {
public:
    GreeterClient(std::shared_ptr<Channel> channel)
        : stub_(Greeter::NewStub(channel)) {}

    std::string SayHello(const std::string& user) {
        HelloRequest request;
        request.set_name(user);

        HelloResponse reply;
        ClientContext context;

        Status status = stub_->Hello(&context, request, &reply);

        if (status.ok()) {
            return reply.response();
        } else {
            std::cout << status.error_code() << ": " << status.error_message()
                << std::endl;
            return "RPC failed";
        }
    }

    void SetKey(const std::string& key, const std::string& value) {
        Hello::SetKeyRequest request;
        request.set_key(key);
        request.set_value(value);

        ClientContext context;
        Hello::KeyResponse reply;

        stub_->SetValue(&context, request, &reply);
    }

    std::string GetKey(const std::string& key) {
        Hello::KeyRequest request;
        request.set_key(key);

        ClientContext context;
        Hello::KeyResponse reply;

        Status status = stub_->GetValue(&context, request, &reply);
        if (status.ok()) {
            return reply.value();
        }

        return "not found";
    }



private:
        std::unique_ptr<Greeter::Stub> stub_;
};

int main(int argc, char** argv) {
    GreeterClient greeter(
        grpc::CreateChannel("localhost:50051", grpc::InsecureChannelCredentials()));

    const std::string user("world");
    const std::string reply = greeter.SayHello(user);

    std::cout << "Greeter received: " << reply << std::endl;

    greeter.SetKey("1", "1");
    greeter.SetKey("2", "2");

    std::cout << "got keys: "
        << "1: " << greeter.GetKey("1") << std::endl
        << "2: " << greeter.GetKey("2") << std::endl;
}
