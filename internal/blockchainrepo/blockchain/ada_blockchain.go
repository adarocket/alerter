package blockchain

import (
	"database/sql"
	"github.com/adarocket/alerter/internal/client"
	"github.com/adarocket/alerter/internal/database/model"
	"github.com/adarocket/alerter/internal/msgsender"
	"github.com/adarocket/alerter/internal/nodesinfo"
	"github.com/adarocket/proto/proto-gen/cardano"
	"google.golang.org/grpc"
	"log"
)

type Cardano struct {
	Blockchain      string
	db              model.ModelAlertNode
	cardanoClient   *client.CardanoClient
	informClient    *client.ControllerClient
	oldCardanoNodes map[string]*cardano.SaveStatisticRequest
}

func (c *Cardano) ConnectServices(clientConn *grpc.ClientConn, dbConn *sql.DB) {
	c.informClient = client.NewControllerClient(clientConn)
	c.cardanoClient = client.NewCardanoClient(clientConn)
	c.db = model.NewAlertNodeInstance(dbConn)

	if len(c.oldCardanoNodes) == 0 {
		c.oldCardanoNodes = make(map[string]*cardano.SaveStatisticRequest)
	}
}

func (c *Cardano) CreateInfoStatMsg() (map[msgsender.KeyMsg]msgsender.BodyMsg, error) {
	resp, err := c.informClient.GetNodeList()
	if err != nil {
		log.Println(err)
		return map[msgsender.KeyMsg]msgsender.BodyMsg{}, err
	}

	messages := map[msgsender.KeyMsg]msgsender.BodyMsg{}

	for _, node := range resp.NodeAuthData {
		if node.Blockchain == c.Blockchain || node.Blockchain == "" { // temp " " == cardano
			alerts, err := c.db.GetAlertsByNodeUuid(node.Uuid)
			if err != nil {
				log.Println(err)
				return messages, err
			} else if len(alerts) == 0 {
				//continue
			}

			cardanoNode, err := c.cardanoClient.GetStatistic(node.Uuid)
			if err != nil {
				log.Println(err)
				continue
			}
			msg, _ := nodesinfo.CheckFieldsOfNode(cardanoNode, c.oldCardanoNodes[node.Uuid], alerts)
			for keyMsg, bodyMsg := range msg {
				messages[keyMsg] = bodyMsg
			}

			c.oldCardanoNodes[node.Uuid] = cardanoNode
		}
	}

	return messages, nil
}
