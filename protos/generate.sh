#!/bin/bash

protoc --proto_path=./proto \
  --go_out=. \
  --go_opt=Mauth.proto=./gen/authGrpc \
  --go-grpc_out=. \
  --go-grpc_opt=Mauth.proto=./gen/authGrpc \
  auth.proto