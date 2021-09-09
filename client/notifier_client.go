package client

import (
	"context"
	pb "github.com/adarocket/alerter/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
	"time"
)

var target = "127.0.0.1:5300"

type NotifierClient struct {
	client pb.ReverseClient
}

func NewNotifierClient() (*NotifierClient, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		grpclog.Errorf("fail to dial: %v", err)
		log.Println(err)
		return nil, err
	}

	client := pb.NewReverseClient(conn)

	return &NotifierClient{client}, nil
}

func (c *NotifierClient) SendMessage(msg string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.client.Do(ctx, &pb.Request{Message: msg})
	if err != nil {
		grpclog.Errorf("fail to dial: %v", err)
		log.Println(err)
		return err
	}

	return nil
}
