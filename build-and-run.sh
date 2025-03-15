#!/bin/bash
# Script to build WASM and start the server

cd ebiten

echo "Building WebAssembly..."
GOOS=js GOARCH=wasm go build -o ../site/main.wasm

# Check if build was successful
if [ $? -eq 0 ]; then
  echo "WebAssembly build successful"
else
  echo "WebAssembly build failed"
  exit 1
fi

cd ../site

echo "Starting server..."
go run main.go