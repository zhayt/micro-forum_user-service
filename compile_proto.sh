#!/bin/bash

PROTO_FILES_PATH=/home/zhayt/google-dev/golang/micro-forum/micro-forum-proto
OUTPUT_PATH=/home/zhayt/google-dev/golang/micro-forum/user-service

mkdir -p $OUTPUT_PATH/proto

protoc \
    --proto_path=$PROTO_FILES_PATH \
    --go_out=$OUTPUT_PATH/proto \
    --go_opt=paths=source_relative \
    --go-grpc_out=$OUTPUT_PATH/proto \
    --go-grpc_opt=paths=source_relative \
    $PROTO_FILES_PATH/user.proto