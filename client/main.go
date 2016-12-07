package main

import (
	"log"
	"os"

	"github.com/yuichi1004/grpc-experiments/fibo"
	"github.com/yuichi1004/grpc-experiments/token"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	TokenAddress string
	FiboAddress  string
	CAPath       string
)

func init() {
	TokenAddress = os.Getenv("TOKEN_SERVER_ADDR")
	FiboAddress = os.Getenv("FIBO_SERVER_ADDR")
	CAPath = os.Getenv("CA_CERT_PATH")
}

func main() {
	ctx := context.Background()

	ca, err := credentials.NewClientTLSFromFile(CAPath, "")
	if err != nil {
		panic(err)
	}

	// Get token
	tokenConn, err := grpc.Dial(TokenAddress,
		grpc.WithTransportCredentials(ca),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer tokenConn.Close()
	tokenClient := token.NewTokenClient(tokenConn)
	resp, err := tokenClient.GetToken(ctx, &token.TokenRequest{Subject: "me", Scope: []string{"fibo"}})
	if err != nil {
		panic(err)
	}

	// Connect to fibo server
	fiboConn, err := grpc.Dial(FiboAddress,
		grpc.WithTransportCredentials(ca),
		grpc.WithPerRPCCredentials(TokenCreds{resp.Token}),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer fiboConn.Close()
	fiboClient := fibo.NewFibonacciClient(fiboConn)

	// Calc fibonacci series
	for i := 0; i < 20; i++ {
		r, err := fiboClient.GetN(ctx, &fibo.FibonacciRequest{N: int64(i)})
		if err != nil {
			log.Fatalf("could not get fibo: %v", err)
		}
		log.Printf("Result: %d, %d", i, r.Result)
	}
}

// Token credentials
type TokenCreds struct {
	Token string
}

func (cred TokenCreds) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": cred.Token,
	}, nil
}

func (_ TokenCreds) RequireTransportSecurity() bool {
	return false
}
