# gRPC Python to Go MLX Example

## Setup Requirements

```
# Install proto compiler for your OS
$ brew install protobuf

# Python dependencies
$ python3 -m pip install grpcio
$ python3 -m pip install grpcio-tools
$ python3 -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. TextGeneration.proto

# Go dependencies
# https://grpc.io/docs/languages/go/quickstart/
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative TextGeneration.proto
```

I manually moved some of the files into the `pkg/eproto` folder after generating them with the `protoc` command and updated the code imports to work with those. This will be fixed once this example is fully implemented into the main application.