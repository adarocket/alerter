package client

import (
	"context"
	"time"

	pb "github.com/adarocket/proto/proto-gen/auth"

	"google.golang.org/grpc"
)

// AuthClient is a client to call authentication RPC
type AuthClient struct {
	service pb.AuthServiceClient
}

// NewAuthClient returns a new auth client
func NewAuthClient(cc *grpc.ClientConn) *AuthClient {
	service := pb.NewAuthServiceClient(cc)
	return &AuthClient{service}
}

// Login login user and returns the access token
func (client *AuthClient) Login(username string, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.LoginRequest{
		Username: username,
		Password: password,
	}

	res, err := client.service.Login(ctx, req)
	if err != nil {
		return "", err
	}

	return res.GetAccessToken(), nil
}
