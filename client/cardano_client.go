package client

import (
	"context"
	pb "github.com/adarocket/proto/proto-gen/cardano"
	"google.golang.org/grpc"
	"time"
)

// InformClient is a client to call laptop service RPCs
type CardanoClient struct {
	service pb.CardanoClient
}

// NewInformClient -
func NewCardanoClient(cc *grpc.ClientConn) *CardanoClient {
	service := pb.NewCardanoClient(cc)
	return &CardanoClient{service}
}

// GetStatistic -
func (informClient *CardanoClient) GetStatistic(nodeUUID string) (resp *pb.SaveStatisticRequest, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.GetStatisticRequest{
		Uuid: nodeUUID,
	}

	return informClient.service.GetStatistic(ctx, req)
}
