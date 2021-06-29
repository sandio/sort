package main

import (
	"fmt"
	"log"
	"net"

	"github.com/sandio/sort/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const serverPort = "localhost:10001"
const clientPort = "localhost:10000"

func main() {
	client, conn := newSortingRobotClient()
	defer conn.Close()
	grpcServer, lis := newFulfillmentServer(client)

	fmt.Printf("gRPC server started. Listening on %s\n", serverPort)
	grpcServer.Serve(lis)
}

func newSortingRobotClient() (gen.SortingRobotClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(clientPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	return gen.NewSortingRobotClient(conn), conn
}

func newFulfillmentServer(client gen.SortingRobotClient) (*grpc.Server, net.Listener) {
	lis, err := net.Listen("tcp", serverPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	gen.RegisterFulfillmentServer(grpcServer, newFulfillmentService(client))
	reflection.Register(grpcServer)

	return grpcServer, lis
}
