package main

import (
	pb "api-grpc/places"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
	"time"
)

const (
	port           = ":50051"
	caCertFile     = "cert/ca-cert.pem"
	serverCertFile = "cert/server-cert.pem"
	serverKeyFile  = "cert/server-key.pem"
)

type server struct {
	pb.UnimplementedPlacesServer
}

// Loads TLS credentials
func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed client's certificate
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}

// Unary interceptor to handle logging and auth
func unaryInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// TODO auth

	h, err := handler(ctx, req)

	// logging
	log.Printf("Method:%s\tResponseTime:%s\tError:%v\n", info.FullMethod, time.Since(start), err)

	return h, err
}

// Returns a place by id
func (s *server) GetById(ctx context.Context, in *pb.PlaceIdRequest) (*pb.PlaceResponse, error) {

	// TODO request storage
	dummyPlace := &pb.Place{
		Id:       in.GetId(),
		Name:     "Dummy name",
		Location: "Dummy location",
	}

	// test output
	log.Printf("Received: get place by id: %v", in.GetId())

	return &pb.PlaceResponse{Place: dummyPlace}, nil
}

func main() {
	//creds, _ := credentials.NewServerTLSFromFile(certFile, keyFile)
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("Cannot load TLS credentials: ", err)
	}
	s := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.UnaryInterceptor(unaryInterceptor),
	)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	pb.RegisterPlacesServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
