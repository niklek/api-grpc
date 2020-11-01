package main

import (
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
	"context"
	"strconv"
	pb "api-grpc/places"
)

const (
	address   = "localhost:50051"
	defaultId = int64(42)
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPlacesClient(conn)

	placeId := defaultId
	if len(os.Args) > 1 {
		placeId, _ = strconv.ParseInt(os.Args[1], 10, 64)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	p, err := c.GetById(ctx, &pb.PlaceIdRequest{Id: placeId})
	if err != nil {
		log.Fatalf("could not get place by id: %d error: %v", placeId, err)
	}
	log.Printf("Place id: %#v", p.GetPlace())
}

