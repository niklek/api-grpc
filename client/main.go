package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"os"
	"time"
	"context"
	"strconv"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"fmt"
	pb "api-grpc/places"
)

const (
	serverAddress = "localhost:50051"
	caCertFile = "cert/ca-cert.pem"
	defaultId = int64(42)
)

// Loads TLS credentials
func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	// connect to the server
	//creds, _ := credentials.NewClientTLSFromFile(certFile, "")
	creds, err := loadTLSCredentials()
    if err != nil {
        log.Fatal("cannot load TLS credentials: ", err)
    }

    conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to connect to the server: %v", err)
	}
	defer conn.Close()
	c := pb.NewPlacesClient(conn)

	// place id
	placeId := defaultId
	if len(os.Args) > 1 {
		placeId, _ = strconv.ParseInt(os.Args[1], 10, 64)
	}

	// get a place by id
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	p, err := c.GetById(ctx, &pb.PlaceIdRequest{Id: placeId})
	if err != nil {
		log.Printf("could not get place by id: %d error: %v", placeId, err)
	}
	log.Printf("Place id: %d name: %s", p.GetPlace().GetId(), p.GetPlace().GetName())
}
