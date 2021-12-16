package noderepo

import (
	"fmt"
	"github.com/adarocket/alerter/internal/client"
	"github.com/adarocket/alerter/internal/database"
	"github.com/adarocket/alerter/internal/structs"
	"google.golang.org/grpc"
)

type NodeRepository struct {
	notifyClient *client.NotifierClient
	clientConn   *grpc.ClientConn
	nodes        []structs.NodesBlockChain
}

func InitNodeRepository(notifyClient *client.NotifierClient, nodesBlockChains []structs.NodesBlockChain,
	clientConn *grpc.ClientConn, db database.ModelAlertNode) NodeRepository {

	for _, chain := range nodesBlockChains {
		chain.Init(clientConn, db)
	}

	return NodeRepository{
		notifyClient: notifyClient,
		clientConn:   clientConn,
		nodes:        nodesBlockChains,
	}
}

func (r *NodeRepository) ProcessStatistic() {
	for _, node := range r.nodes {
		statMsg, _ := node.CreateInfoStatMsg()
		fmt.Println(statMsg)
	}
}
