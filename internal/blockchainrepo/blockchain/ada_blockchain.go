package blockchain

import (
	"database/sql"
	"github.com/adarocket/alerter/internal/cache"
	"github.com/adarocket/alerter/internal/client"
	"github.com/adarocket/alerter/internal/database/model"
	"github.com/adarocket/alerter/internal/msgsender"
	"github.com/adarocket/alerter/internal/nodesinfo"
	"google.golang.org/grpc"
	"log"
)

type Cardano struct {
	Blockchain    string
	db            model.ModelAlertNode
	cardanoClient *client.CardanoClient
	informClient  *client.ControllerClient
}

func (c *Cardano) Init(clientConn *grpc.ClientConn, dbConn *sql.DB) {
	c.informClient = client.NewControllerClient(clientConn) // common controller вынести выше
	c.cardanoClient = client.NewCardanoClient(clientConn)
	c.db = model.NewAlertNodeInstance(dbConn)
}

func (c *Cardano) CreateInfoStatMsg() (map[msgsender.KeyMsg]msgsender.BodyMsg, error) {
	resp, err := c.informClient.GetNodeList()
	if err != nil {
		log.Println(err)
		return map[msgsender.KeyMsg]msgsender.BodyMsg{}, nil
	}

	cacheInstance := cache.GetCacheInstance()
	messages := map[msgsender.KeyMsg]msgsender.BodyMsg{}

	for _, node := range resp.NodeAuthData {
		if node.Blockchain == c.Blockchain || node.Blockchain == "" { // temp " " == cardano
			alerts, err := c.db.GetAlertsByNodeUuid(node.Uuid)
			if err != nil {
				log.Println(err)
				continue
			} else if len(alerts) == 0 {
				continue
			}

			cardanoNode, _ := c.cardanoClient.GetStatistic(node.Uuid)
			oldCardanoNode := cacheInstance.GetCardanoNodes([]string{node.Uuid})

			msg, _ := nodesinfo.CheckFieldsOfNode(cardanoNode, oldCardanoNode[node.Uuid], alerts)
			for keyMsg, bodyMsg := range msg {
				messages[keyMsg] = bodyMsg
			}

			cacheInstance.AddCardanoData(cardanoNode)
		}
	}

	return messages, nil
}
