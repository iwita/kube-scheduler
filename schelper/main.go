package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	//if we crash the go code, we get the file name and row number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//Connect to mongoDB
	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = client.Connect(context.TODO())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//Use a Collection

	// collection = client.Database("mydb").Collection("machines")

	// fmt.Println("Machine Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}

	grpcServer := grpc.NewServer(opts...)
	s := socket.Server{}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

	socket.RegisterSocketServiceServer(grpcServer)
}
