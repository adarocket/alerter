package structs

import (
	"github.com/adarocket/alerter/internal/database"
	"github.com/adarocket/alerter/internal/nodesinfo/msgsender"
	"google.golang.org/grpc"
)

type NodesBlockChain interface {
	Init(clientConn *grpc.ClientConn, db database.ModelAlertNode)
	CreateInfoStatMsg() (map[msgsender.KeyMsg]msgsender.BodyMsg, error)
}
