package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/yuichi1004/grpc-experiments/token"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	certPath := os.Getenv("TOKEN_CERT_PATH")
	keyPath := os.Getenv("TOKEN_KEY_PATH")
	ServerCert, err = tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		panic(err)
	}

	ServerPort = os.Getenv("TOKEN_SERVER_PORT")
}

type server struct{}

// Token
func (s *server) GetToken(ctx context.Context, in *token.TokenRequest) (*token.TokenReply, error) {
	claims := Claims{
		Subject: in.Subject,
		Scopes:  in.Scope,
	}
	tokenPayload, err := json.Marshal(&claims)
	if err != nil {
		return nil, fmt.Errorf("failed to issue access token")
	}
	return &token.TokenReply{Token: string(tokenPayload)}, nil
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
	token.RegisterTokenServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
