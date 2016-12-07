package main

import (
	"log"

	pb "github.com/yuichi1004/grpc-experiments/fibo"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address = "fibo.example.com:50051"
)

func init() {
}

func GenToken() string {
	return "XXXX"
}

func main() {
	ca, err := credentials.NewClientTLSFromFile("../creds/ca/cacert.pem", "fibo.example.com")
	if err != nil {
		panic(err)
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(address,
		grpc.WithTransportCredentials(ca),
		grpc.WithPerRPCCredentials(JWTCreds{}),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewFibonacciClient(conn)

	// Contact the server and print out its response.
	for i := 0; i < 20; i++ {
		r, err := c.GetN(context.Background(), &pb.FibonacciRequest{N: int64(i)})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Result: %d, %d", i, r.Result)
	}
}

type JWTCreds struct {
}

func (_ JWTCreds) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + GenToken(),
	}, nil
}

func (_ JWTCreds) RequireTransportSecurity() bool {
	return false
}
