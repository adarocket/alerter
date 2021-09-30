package client

import (
	"context"
	pb "github.com/adarocket/alerter/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
	"time"
)

var NotifyTarget = "127.0.0.1:5300"

type NotifierClient struct {
	client pb.NotifierConnectClient
}

func NewNotifierClient() (*NotifierClient, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial(NotifyTarget, opts...)
	if err != nil {
		grpclog.Errorf("fail to dial: %v", err)
		log.Println(err)
		return nil, err
	}

	client := pb.NewNotifierConnectClient(conn)

	return &NotifierClient{client}, nil
}

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

func (c *NotifierClient) SendMessages(msges []*pb.SendNotifier) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, msg := range msges {
		_, err := c.client.SendNotification(ctx, msg)
		if err != nil {
			grpclog.Errorf("fail to dial: %v", err)
			log.Println(err)
			continue
		}
	}

	return nil
}
