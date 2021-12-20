package blockchain

import (
	"database/sql"
	"github.com/adarocket/alerter/internal/msgsender"
	"google.golang.org/grpc"
)

type NodesBlockChain interface {
	ConnectServices(clientConn *grpc.ClientConn, db *sql.DB)
	CreateInfoStatMsg() (map[msgsender.KeyMsg]msgsender.BodyMsg, error)
}
