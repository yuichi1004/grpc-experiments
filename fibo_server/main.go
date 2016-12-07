package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	fibo "github.com/yuichi1004/grpc-experiments/fibo"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var (
	ServerCert tls.Certificate
	ServerPort string
)

type (
	Claims struct {
		Subject string
		Scopes  []string
	}
)

func init() {
	var err error
	certPath := os.Getenv("FIBO_CERT_PATH")
	keyPath := os.Getenv("FIBO_KEY_PATH")
	ServerCert, err = tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		panic(err)
	}

	ServerPort = os.Getenv("FIBO_SERVER_PORT")
}

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

type server struct{}

func (s *server) GetN(ctx context.Context, in *fibo.FibonacciRequest) (*fibo.FibonacciReply, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	token, ok := md["authorization"]
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	tokenPayload := []byte(token[0])
	var claims Claims
	if err := json.Unmarshal(tokenPayload, &claims); err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	for _, scope := range claims.Scopes {
		if scope == "fibo" {
			return &fibo.FibonacciReply{Result: Fibo(in.N)}, nil
		}
	}
	return nil, fmt.Errorf("unauthorized (require fibo scope)")
}

func main() {
	log.Printf("Listening on port %s", ServerPort)
	lis, err := net.Listen("tcp", ServerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	transportSecurity := credentials.NewServerTLSFromCert(&ServerCert)
	s := grpc.NewServer(
		grpc.Creds(transportSecurity),
	)
	fibo.RegisterFibonacciServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
