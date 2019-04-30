#!/bin/bash

# 输出目录
GO_PUT_PATH='./'

protoc --go_out=paths=source_relative,plugins=grpc:$GO_PUT_PATH health.proto