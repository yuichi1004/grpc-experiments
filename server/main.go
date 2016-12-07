package main

import (
	"log"
	"net"
	"io/ioutil"
	"fmt"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	pb "github.com/yuichi1004/grpc-experiments/fibo"
)

const (
	address = "localhost:50051"
	port    = ":50051"
)

var (
	publicKey []byte
)

func init() {
	publicKey, _ = ioutil.ReadFile("../creds/demo.rsa.pub")
}

// server is used to implement helloworld.GreeterServer.
type server struct{}

func Fibo(n int64) int64 {
	a := int64(0)
	b := int64(1)
	for i := int64(0); i < n; i++ {
		tmp := b
		b = a + b
		a = tmp
	}
	return a
}

// SayHello implements helloworld.GreeterServer
func (s *server) GetN(ctx context.Context, in *pb.FibonacciRequest) (*pb.FibonacciReply, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	log.Printf("Metadata %+v", md)
	jwtToken, ok := md["authorization"]
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	creds := strings.Split(jwtToken[0], " ")
	if creds[0] != "Bearer" {
		return nil, fmt.Errorf("unauthorized")
	}

	token, err := jwt.Parse(creds[1], func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			log.Printf("Unexpected signing method: %v", t.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		return publicKey, nil
	})
	if err == nil && token.Valid {
		return &pb.FibonacciReply{Result: Fibo(in.N)}, nil
	}

	return nil, fmt.Errorf("unauthorized")
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	transportSecurity, err := credentials.NewServerTLSFromFile("../creds/ca.pem", "../creds/ca.key")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(
		grpc.Creds(transportSecurity),
	)
	pb.RegisterFibonacciServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
