package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"context"
	pb "api-grpc/places"
)

const port = ":50051"

type server struct {
	pb.UnimplementedPlacesServer
}

func (s *server) GetById(ctx context.Context, in *pb.PlaceIdRequest) (*pb.PlaceResponse, error) {

	dummyPlace := &pb.Place{
		Id: in.GetId(),
		Name: "Dummy name",
		Location: "Dummy location",
	}
	log.Printf("Received: get place by id: %v", in.GetId())
	return &pb.PlaceResponse{Place: dummyPlace}, nil
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

