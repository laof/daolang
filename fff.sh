#!/bin/bash

export _LAMBDA_SERVER_PORT=9006
export AWS_LAMBDA_RUNTIME_API=127.0.0.1:9006

# 启动你的 Go 应用
go run netlify/functions/hello/main.go