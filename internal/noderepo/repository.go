package noderepo

import (
	"github.com/adarocket/alerter/internal/client"
	"github.com/adarocket/alerter/internal/database"
	"github.com/adarocket/alerter/internal/nodesinfo/msgsender"
	"github.com/adarocket/alerter/internal/structs"
	"google.golang.org/grpc"
)

type NodeRepository struct {
	notificator msgsender.MsgSender
	clientConn  *grpc.ClientConn
	blockChains []structs.NodesBlockChain
}

func InitNodeRepository(notifyClient *client.NotifierClient, nodesBlockChains []structs.NodesBlockChain,
	clientConn *grpc.ClientConn, db database.ModelAlertNode) NodeRepository {

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
