package client

import (
	"context"
	pb "github.com/adarocket/proto/proto-gen/cardano"
	"google.golang.org/grpc"
	"time"
)

// CardanoClient - is a client to call laptop service RPCs of cardano
type CardanoClient struct {
	service pb.CardanoClient
}

// NewCardanoClient - returns a new cardano client
func NewCardanoClient(cc *grpc.ClientConn) *CardanoClient {
	service := pb.NewCardanoClient(cc)
	return &CardanoClient{service}
}

// GetStatistic - get statistics of cardano nodes
func (informClient *CardanoClient) GetStatistic(nodeUUID string) (resp *pb.SaveStatisticRequest, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.GetStatisticRequest{
		Uuid: nodeUUID,
	}

	return informClient.service.GetStatistic(ctx, req)
}
