package blockchainrepo

import (
	"database/sql"
	"github.com/adarocket/alerter/internal/blockchainrepo/blockchain"
	"github.com/adarocket/alerter/internal/client"
	"github.com/adarocket/alerter/internal/msgsender"
	"google.golang.org/grpc"
)

type NodeRepository struct {
	notificator msgsender.MsgSender
	clientConn  *grpc.ClientConn
	blockChains []blockchain.NodesBlockChain
}

func InitNodeRepository(notifyClient *client.NotifierClient, nodesBlockChains []blockchain.NodesBlockChain,
	clientConn *grpc.ClientConn, db *sql.DB) NodeRepository {

	for _, chain := range nodesBlockChains {
		chain.Init(clientConn, db)
	}

	return NodeRepository{
		notificator: msgsender.CreateMsgSender(notifyClient),
		clientConn:  clientConn,
		blockChains: nodesBlockChains,
	}
}

func (r *NodeRepository) ProcessStatistic() {
	statsMsges := make(map[msgsender.KeyMsg]msgsender.BodyMsg)
	for _, node := range r.blockChains {
		statMsg, _ := node.CreateInfoStatMsg()
		for key, bodyMsg := range statMsg {
			statsMsges[key] = bodyMsg
		}
	}

	r.notificator.AddNotifiersToStack(statsMsges) // TODO remove map struct
}
