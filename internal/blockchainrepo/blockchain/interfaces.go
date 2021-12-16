package blockchain

import (
	"github.com/adarocket/alerter/internal/database/db"
	"github.com/adarocket/alerter/internal/msgsender"
	"google.golang.org/grpc"
)

type NodesBlockChain interface {
	Init(clientConn *grpc.ClientConn, db db.ModelAlertNode)
	CreateInfoStatMsg() (map[msgsender.KeyMsg]msgsender.BodyMsg, error)
}
