package main

import (
	"log"
	"io/ioutil"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	jwt "github.com/dgrijalva/jwt-go"
	pb "github.com/yuichi1004/grpc-experiments/fibo"
)

const (
	address     = "localhost:50051"
)

var (
	privateKey []byte
)

func init() {
	privateKey, _ = ioutil.ReadFile("../creds/demo.rsa")
}

func GenToken() string {
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims["exp"] = time.Now().Unix() + 36000
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		panic(err)
	}
	return tokenString
}

func main() {
	cred, err := credentials.NewClientTLSFromFile("../creds/server.crt", "localhost")
	if err != nil {
		panic(err)
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(address,
		grpc.WithTransportCredentials(cred),
		grpc.WithPerRPCCredentials(JWTCreds{}),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewFibonacciClient(conn)

	// Contact the server and print out its response.
	for i := 0; i < 20; i ++ {
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
	return map[string]string {
		"authorization": "Bearer " + GenToken(),
	},nil
}

func (_ JWTCreds) RequireTransportSecurity() bool {
	return false
}

