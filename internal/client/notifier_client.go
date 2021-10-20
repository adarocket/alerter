package client

import (
	"context"
	pb "github.com/adarocket/proto/proto-gen/notifier"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
	"time"
)

// NotifierClient is a client to send notifiers service RPCs
type NotifierClient struct {
	client pb.NotifierConnectClient
}

func NewNotifierClient(addr string) (*NotifierClient, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		grpclog.Errorf("fail to dial: %v", err)
		log.Println(err)
		return nil, err
	}

	client := pb.NewNotifierConnectClient(conn)

	return &NotifierClient{client}, nil
}

// SendMessage - send notifier to notifierServer
func (c *NotifierClient) SendMessage(msg *pb.SendNotifier) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.client.SendNotification(ctx, msg)
	if err != nil {
		grpclog.Errorf("fail to dial: %v", err)
		log.Println(err)
		return err
	}

	return nil
}
