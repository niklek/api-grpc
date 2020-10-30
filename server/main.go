package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"context"
	"strconv"
	pb "api-grpc/places"
)

const port = ":50051"

type server struct {
	pb.UnimplementedPlacesServer
}

func (s *server) GetById(ctx context.Context, in *pb.PlaceIdRequest) (*pb.PlaceResponse, error) {
	log.Printf("Received: get place by id: %v", in.GetId())
	return &pb.PlaceResponse{Message: "Place id: " + strconv.Itoa(int(in.GetId()))}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPlacesServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

