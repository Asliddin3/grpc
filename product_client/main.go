package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "gitlab.com/go/grpc/genproto/product"
)

func main() {
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.CreateProduct(ctx, &pb.CreateProductRequest{
		Name:       "new name",
		Categoryid: 2,
		Typeid:     1,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Product name: %s", r.Name)
}
