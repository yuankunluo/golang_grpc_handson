package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "tiny-blog/blogpb"
)

var collection *mongo.Collection

type server struct{}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"authod_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title`
}

func main() {
	// Set Log Level.
	// If we crush the code, we get the file name and line number.
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Hello From Blog Service Server")

	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}
	certFile := "./ssl/server.crt"
	keyFile := "./ssl/server.pem"
	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	if sslErr != nil {
		log.Fatalf("Failed loading certificates: %v\n", sslErr)
		return
	}
	opts := grpc.Creds(creds)

	// Connect to MongoDB
	mongoClient, mongoErr := mongo.NewClient(options.Client().ApplyURI("mongodb://mongoadmin:mymongopass123++@localhost:27017"))
	if mongoErr != nil {
		log.Fatalf("Failed to create MongoDB client: %v\n", mongoErr)
	}
	mongoConnError := mongoClient.Connect(context.TODO())
	if mongoConnError != nil {
		log.Fatalf("Failed to connect to MongoDB: %v\n", mongoConnError)
	}
	fmt.Println("Connect to MongoDB succeeded.")
	collection = mongoClient.Database("tiny-blog-db").Collection("blog")

	// Create Server.
	s := grpc.NewServer(opts)
	// Register Service Server.
	pb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Start serving...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v\n", err)
		}
	}()

	// Wait fro Control-C to exit.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	// Block until Control-C received.
	<-ch
	s.Stop()
	fmt.Println("Server stopped smoothly.")
	mongoClient.Disconnect(context.TODO())
	fmt.Println("Disconnect from MongoDB.")
	lis.Close()
	fmt.Println("Listener was closed smoothly.")

	fmt.Println("See you again!")
}
