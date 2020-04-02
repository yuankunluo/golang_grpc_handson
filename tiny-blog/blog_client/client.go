package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/credentials"

	pb "tiny-blog/blogpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a Tine-Blog client")

	tls := true
	opts := grpc.WithInsecure()

	if tls {
		certFile := "./ssl/ca.crt" // Certificate Authority Trust Certificate
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		opts = grpc.WithTransportCredentials(creds)
		if sslErr != nil {
			log.Fatalf("Failed to load CA cert: %v\n", sslErr)
			return
		}
	}

	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Count not connect: %v\n", err)
	}
	defer conn.Close()

	c := pb.NewBlogServiceClient(conn)

	blog := &pb.Blog{
		AuthorId: "yuankunluo",
		Title:    "My First Blog",
		Content:  "Content of the first blog...",
	}
	createRes, createErr := c.CreateBlog(
		context.Background(),
		&pb.CreateBlogRequest{
			Blog: blog,
		})
	if createErr != nil {
		log.Fatalf("Error while requesting create blog: %v\n", createErr)
	}
	fmt.Printf("Blog has been created: %v\n", createRes.GetBlog())
}
