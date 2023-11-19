cmake_minimum_required(VERSION 3.21)
project(grpc_example C CXX)

set(CMAKE_CXX_STANDARD 20)

set(protobuf_MODULE_COMPATIBLE TRUE)
find_package(Protobuf CONFIG REQUIRED)
message(STATUS "Using protobuf ${Protobuf_VERSION}")

set(_PROTOBUF_LIBPROTOBUF protobuf::libprotobuf)
set(_REFLECTION gRPC::grpc++_reflection)
set(_PROTOBUF_PROTOC $<TARGET_FILE:protobuf::protoc>)

find_package(gRPC CONFIG REQUIRED)
message(STATUS "Using gRPC ${gRPC_VERSION}")

set(_GRPC_GRPCPP gRPC::grpc++)
set(_GRPC_CPP_PLUGIN_EXECUTABLE $<TARGET_FILE:gRPC::grpc_cpp_plugin>)

get_filename_component(hw_proto "protos/hello.proto" ABSOLUTE)
get_filename_component(hw_proto_path "${hw_proto}" PATH)

# Generated sources
set(hw_proto_srcs "${CMAKE_CURRENT_BINARY_DIR}/hello.pb.cc")
set(hw_proto_hdrs "${CMAKE_CURRENT_BINARY_DIR}/hello.pb.h")
set(hw_grpc_srcs "${CMAKE_CURRENT_BINARY_DIR}/hello.grpc.pb.cc")
set(hw_grpc_hdrs "${CMAKE_CURRENT_BINARY_DIR}/hello.grpc.pb.h")

add_custom_command(
        OUTPUT "${hw_proto_srcs}" "${hw_proto_hdrs}" "${hw_grpc_srcs}" "${hw_grpc_hdrs}"
        COMMAND ${_PROTOBUF_PROTOC}
        ARGS --grpc_out "${CMAKE_CURRENT_BINARY_DIR}"
        --cpp_out "${CMAKE_CURRENT_BINARY_DIR}"
        -I "${hw_proto_path}"
        --plugin=protoc-gen-grpc="${_GRPC_CPP_PLUGIN_EXECUTABLE}"
        "${hw_proto}"
        DEPENDS "${hw_proto}")

include_directories("${CMAKE_CURRENT_BINARY_DIR}")

add_library(hw_grpc_proto
    ${hw_grpc_srcs}
    ${hw_grpc_hdrs}
    ${hw_proto_srcs}
    ${hw_proto_hdrs})

target_link_libraries(hw_grpc_proto
    ${_REFLECTION}
    ${_GRPC_GRPCPP}
    ${_PROTOBUF_LIBPROTOBUF})

foreach(_target hello_client hello_server)
    add_executable(${_target} "${_target}.cpp")
    target_link_libraries(${_target}
            hw_grpc_proto
            ${_PROTOBUF_LIBRARY}
            ${_REFLECTION}
            ${_GRPC_GRPCPP}
    )
endforeach()