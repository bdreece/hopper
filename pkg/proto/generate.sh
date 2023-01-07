#!/bin/bash

rm -rf ./**/*.pb.go grpc
mkdir -p grpc
protoc --proto_path=. \
    --go_out=. \
    --go_opt="paths=source_relative" \
    --go-grpc_out=./grpc \
    --go-grpc_opt="paths=source_relative" \
    ./*.proto

mv services.pb.go grpc