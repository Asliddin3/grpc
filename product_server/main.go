package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "gitlab.com/go/grpc/genproto/product"
	db "gitlab.com/go/grpc/postgres"
)

type server struct {
	pb.ProductServiceServer
}

func (*server) GetProducts(ctx context.Context, req *pb.Empty) (*pb.ListProductResponse, error) {
	product, err := db.GetProducts()
	if err != nil {
		fmt.Println("error while selecting from database", err)
		return &pb.ListProductResponse{}, err
	}
	var productsResp pb.ListProductResponse
	for _, prod := range product {
		productsResp.Products = append(productsResp.Products, &pb.Product{
			Id:       prod.ID,
			Name:     prod.Name,
			Category: prod.Category,
			Type:     prod.Type,
		})
	}
	fmt.Println(&productsResp)
	return &productsResp, nil
}
func (*server) UpdateProduct(ctx context.Context, req *pb.ProductReq) (*pb.Product, error) {
	product, err := db.UpdateProduct(req)
	if err != nil {
		return &pb.Product{}, err
	}
	return &pb.Product{
		Id:       product.ID,
		Name:     product.Name,
		Category: product.Category,
		Type:     product.Type,
	}, nil
}

func (*server) GetProductInfo(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	index := int64(req.Id)
	product, err := db.GetProductInfo(int64(index))
	if err != nil {
		return &pb.Product{}, err
	}
	productresp := pb.Product{
		Id:       int64(product.ID),
		Name:     string(product.Name),
		Category: string(product.Category),
		Type:     string(product.Type),
	}

	return &productresp, nil
}

func (*server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	product, err := db.CreateProduct(&db.Product{
		Name:       req.Name,
		Categoryid: int(req.Categoryid),
		Typeid:     int(req.Typeid),
	})

	if err != nil {
		fmt.Println(err)
		return &pb.Product{}, err
	}
	fmt.Println(product)
	strores, err := db.CreateStores(product.ID, req.Stores)
	fmt.Println(req.Categoryid)
	fmt.Println(strores)
	if err != nil {
		return &pb.Product{}, err
	}
	var storesResp []*pb.Store
	for _, store := range strores {
		storeResp := pb.Store{
			Id:   store.Id,
			Name: store.Name,
		}
		for _, addres := range store.Addresses {
			adr := pb.Address{
				Id:       int64(addres.Id),
				District: addres.District,
				Street:   addres.Street,
			}
			storeResp.Addresses = append(storeResp.Addresses, &adr)
		}
		storesResp = append(storesResp, &storeResp)
	}
	productResp := pb.Product{
		Id:     product.ID,
		Name:   product.Name,
		Stores: storesResp,
	}
	fmt.Println(productResp)
	return &productResp, nil
}

func (*server) DelelteProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ListProductResponse, error) {
	id := req.Id
	err := db.DelelteProduct(id)
	if err != nil {
		return &pb.ListProductResponse{}, err
	}
	products, err := db.GetProducts()
	if err != nil {
		fmt.Println("error while geting products", err)
		return &pb.ListProductResponse{}, err
	}
	var productsResp pb.ListProductResponse
	for _, prod := range products {
		productsResp.Products = append(productsResp.Products, &pb.Product{
			Id:       prod.ID,
			Name:     prod.Name,
			Category: prod.Category,
			Type:     prod.Type,
		})
	}
	return &productsResp, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterProductServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
