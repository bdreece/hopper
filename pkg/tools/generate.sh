#!/usr/bin/env sh

rm -rf $PROJECT_ROOT/**/*.pb.go $PROJECT_ROOT/pkg/proto/grpc
mkdir -p $PROJECT_ROOT/pkg/proto/grpc

protoc --proto_path=/usr/include \
    --proto_path=$PROJECT_ROOT/pkg/proto \
    --go_out=$PROJECT_ROOT/pkg/proto \
    --go_opt="paths=source_relative" \
    --go-grpc_out=$PROJECT_ROOT/pkg/proto/grpc \
    --go-grpc_opt="paths=source_relative" \
    $PROJECT_ROOT/pkg/proto/*.proto

mv $PROJECT_ROOT/pkg/proto/services.pb.go \
    $PROJECT_ROOT/pkg/proto/grpc