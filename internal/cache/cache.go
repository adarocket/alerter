package cache

import (
	"errors"
	"github.com/adarocket/proto/proto-gen/cardano"
	"github.com/adarocket/proto/proto-gen/chia"
	"log"
)

var (
	instance Cache
)

type Cache struct {
	CardanoNodes map[string]*cardano.SaveStatisticRequest
	ChiaNodes    map[string]*chia.SaveStatisticRequest
}

func init() {
	instance = Cache{}
	instance.CardanoNodes = make(map[string]*cardano.SaveStatisticRequest)
}

func GetCacheInstance() Cache {
	return instance
}

func (c *Cache) AddCardanoData(data interface{}) error {
	if newMapCardano, true := data.(map[string]*cardano.SaveStatisticRequest); true {
		for key, node := range newMapCardano {
			c.CardanoNodes[key] = node
		}
		return nil
	}

	if newCardanoNode, true := data.(*cardano.SaveStatisticRequest); true {
		uuid := newCardanoNode.NodeAuthData.Uuid
		c.CardanoNodes[uuid] = newCardanoNode
		return nil
	}

	return errors.New("invalid type, only mapNodes and node supported")
}

func (c *Cache) GetCardanoNodes(uuids []string) map[string]*cardano.SaveStatisticRequest {
	result := make(map[string]*cardano.SaveStatisticRequest)
	for _, uuid := range uuids {
		if node, isExist := c.CardanoNodes[uuid]; isExist {
			result[uuid] = node
		} else {
			log.Println("node uuid: ", uuid, "does not exist")
		}
	}

	return result
}
