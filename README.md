# api grpc

Test implementation of a server and a client communicating via gRPC.

* Simple client and server written in Go
* gRPC APIs use Protocol Buffers version 3 (proto3)
* Mutual TLS encryption
* * Note: the certificates are not encrypted
* Interceptors on both client and server log timing
* Uses a dummy entity: places and single API to get a place by id

## Usage

1. Generate `.pb.go` files
```sh
make gen
```

2. Generate keys and certs for mutual TLS
```sh
make cert
```

3. Run server (TODO build)
```sh
PORT=50051 CA_CERT_FILENAME=cert/ca-cert.pem SERVER_CERT_FILENAME=cert/server-cert.pem SERVER_KEY_FILENAME=cert/server-key.pem go run server/main.go
```

4. Run client (TODO build)
```sh
# with a param
go run cmd/client/main.go 43
2020/11/08 18:43:51 Place id: 43 name: Dummy name
# using default param 42
go run cmd/client/main.go
2020/11/08 18:43:52 Place id: 42 name: Dummy name
```

## Todos

* Add authentication Middleware
* Graceful shutdown
* ldflags
* Replace standard logger
* Add monitoring
* Add docker file
* Add tests
* Implement a storage
* Add GetAll, Create methods

