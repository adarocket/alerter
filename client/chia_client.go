package client

import (
	"context"
	pb "github.com/adarocket/proto/proto-gen/chia"
	"google.golang.org/grpc"
	"time"
)

// InformClient is a client to call laptop service RPCs
type ChiaClient struct {
	service pb.ChiaClient
}

// NewInformClient -
func NewChiaClient(cc *grpc.ClientConn) *ChiaClient {
	service := pb.NewChiaClient(cc)
	return &ChiaClient{service}
}

// GetStatistic -
func (informClient *ChiaClient) GetStatistic(nodeUUID string) (resp *pb.SaveStatisticRequest, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.GetStatisticRequest{
		Uuid: nodeUUID,
	}

	return informClient.service.GetStatistic(ctx, req)
}
