package client

import (
	"context"
	"time"

	pbCommon "github.com/adarocket/proto/proto-gen/common"

	"google.golang.org/grpc"
)

// ControllerClient is a client to call laptop service RPCs
type ControllerClient struct {
	service pbCommon.ControllerClient
}

// NewControllerClient -
func NewControllerClient(cc *grpc.ClientConn) *ControllerClient {
	service := pbCommon.NewControllerClient(cc)
	return &ControllerClient{service}
}

// GetNodeList - returns list of all nodes
func (informClient *ControllerClient) GetNodeList() (resp *pbCommon.GetNodeListResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pbCommon.GetNodeListRequest{}

	return informClient.service.GetNodeList(ctx, req)
}
