package blockchainrepo

import (
	"database/sql"
	"github.com/adarocket/alerter/internal/blockchainrepo/blockchain"
	"github.com/adarocket/alerter/internal/client"
	"github.com/adarocket/alerter/internal/config"
	"github.com/adarocket/alerter/internal/msgsender"
	"google.golang.org/grpc"
	"log"
)

type NodeRepository struct {
	notificator msgsender.MsgSender
	clientConn  *grpc.ClientConn
	blockChains []blockchain.NodesBlockChain
	db          *sql.DB
	connectFunc func(conf config.Config) (*grpc.ClientConn, error)
}

func InitNodeRepository(notifyClient *client.NotifierClient,
	nodesBlockChains []blockchain.NodesBlockChain,
	connectFunc func(conf config.Config) (*grpc.ClientConn, error),
	conf config.Config,
	db *sql.DB) NodeRepository {
	clientConn, _ := connectFunc(conf)

	return NodeRepository{
		notificator: msgsender.CreateMsgSender(notifyClient),
		clientConn:  clientConn,
		blockChains: nodesBlockChains,
		db:          db,
		connectFunc: connectFunc,
	}
}

func (r *NodeRepository) ConnectNodeRepositoryServices(conf config.Config) {
	conn, _ := r.connectFunc(conf)

	for _, chain := range r.blockChains {
		chain.ConnectServices(conn, r.db)
	}
}

func (r *NodeRepository) ProcessStatistic() {
	statsMsges := make(map[msgsender.KeyMsg]msgsender.BodyMsg)
	for _, node := range r.blockChains {
		statMsg, err := node.CreateInfoStatMsg()
		if err != nil {
			log.Println(err)
			continue
		}

		for key, bodyMsg := range statMsg {
			statsMsges[key] = bodyMsg
		}
	}

	r.notificator.AddNotifiersToStack(statsMsges)
}
