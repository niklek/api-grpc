package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"fmt"
	"net"
	"context"
	pb "api-grpc/places"
)

const (
	port = ":50051"
	caCertFile = "cert/ca-cert.pem"
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
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs: certPool,
	}

    return credentials.NewTLS(config), nil
}

// Returns a place by id
func (s *server) GetById(ctx context.Context, in *pb.PlaceIdRequest) (*pb.PlaceResponse, error) {

	dummyPlace := &pb.Place{
		Id: in.GetId(),
		Name: "Dummy name",
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
        log.Fatal("cannot load TLS credentials: ", err)
    }

    s := grpc.NewServer(grpc.Creds(tlsCredentials))
	
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	pb.RegisterPlacesServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

