package main

import (
	pb "api-grpc/places"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type server struct {
	pb.UnimplementedPlacesServer
}

// Loads TLS credentials
func loadTLSCredentials(caCertFile, serverCertFile, serverKeyFile string) (credentials.TransportCredentials, error) {
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

// Returns a place by id
func (s *server) GetById(ctx context.Context, in *pb.PlaceIdRequest) (*pb.PlaceResponse, error) {

	// TODO request storage
	dummyPlace := &pb.Place{
		Id:       in.GetId(),
		Name:     "Dummy name",
		Location: "Dummy location",
	}

	// test output
	fmt.Println("Received: get place by id:", in.GetId())

	return &pb.PlaceResponse{Place: dummyPlace}, nil
}

func main() {
	logz, _ := zap.NewProduction()
	defer logz.Sync()
	logz.Info("Starting server",
		zap.String("logger", "zap"),
	)

	// Load configuration params from env, no default values
	port := os.Getenv("PORT")
	if port == "" {
		logz.Fatal("PORT is not set")
	}
	caCertFile := os.Getenv("CA_CERT_FILENAME")
	if caCertFile == "" {
		logz.Fatal("CA_CERT_FILENAME is not set")
	}
	serverCertFile := os.Getenv("SERVER_CERT_FILENAME")
	if serverCertFile == "" {
		logz.Fatal("SERVER_CERT_FILENAME is not set")
	}
	serverKeyFile := os.Getenv("SERVER_KEY_FILENAME")
	if serverKeyFile == "" {
		logz.Fatal("SERVER_KEY_FILENAME is not set")
	}
	logz.Info("Configuration is ready",
		zap.String("logger", "zap"),
	)

	// channels for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	shutdown := make(chan error, 2)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	//creds, _ := credentials.NewServerTLSFromFile(certFile, keyFile)
	tlsCredentials, err := loadTLSCredentials(caCertFile, serverCertFile, serverKeyFile)
	if err != nil {
		logz.Fatal("Cannot load TLS credentials: ",
			zap.Error(err),
		)
	}

	// add middleware with zap logger
	grpc_zap.ReplaceGrpcLoggerV2(logz)
	s := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(logz),
		),
	)

	go func() {
		lis, err := net.Listen("tcp", ":"+port)
		if err != nil {
			shutdown <- err
		}
		defer lis.Close()

		pb.RegisterPlacesServer(s, &server{})
		if err := s.Serve(lis); err != nil {
			shutdown <- err
		}
	}()

	// Lister for interruption and gracefully stop the server
	select {
	case x := <-interrupt:
		logz.Info("Received signal", zap.String("signal", x.String()))
	case err := <-shutdown:
		logz.Error("Received an error from server", zap.Error(err))
	}
	logz.Info("Stopping the server...")
	s.GracefulStop()
}
