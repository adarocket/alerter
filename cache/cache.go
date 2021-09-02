package cache

import (
	"github.com/adarocket/proto/proto-gen/cardano"
	"github.com/adarocket/proto/proto-gen/chia"
	"sync"
)

var (
	once     sync.Once
	instance Cache
)

func GetCacheInstance() Cache {
	once.Do(func() {
		instance = Cache{}
	})

	return instance
}

type Cache struct {
	CardanoNodes map[string]*cardano.SaveStatisticRequest
	ChiaNodes    map[string]*chia.SaveStatisticRequest
}

func (c *Cache) AddNewInform(newMap interface{}) {
	if newMapCardano, isTr := newMap.(map[string]*cardano.SaveStatisticRequest); isTr {
		for _, node := range newMapCardano {
			c.CardanoNodes[node.NodeAuthData.Uuid] = node
		}
	}

	if newMapChia, isTr := newMap.(map[string]*chia.SaveStatisticRequest); isTr {
		for _, node := range newMapChia {
			c.ChiaNodes[node.NodeAuthData.Uuid] = node
		}
	}
	// add error if type not exist
}

func (c *Cache) GetMapByType(newNodeMap interface{}) (oldNodeMap interface{}) {
	if _, isTr := newNodeMap.(map[string]*cardano.SaveStatisticRequest); isTr {
		return c.CardanoNodes
	}

	if _, isTr := newNodeMap.(map[string]*chia.SaveStatisticRequest); isTr {
		return c.ChiaNodes
	}

	return nil
}

func (c *Cache) GetOldNodeByType(newNode interface{}) (oldNode interface{}) {
	if newNodeT, isTr := newNode.(*cardano.SaveStatisticRequest); isTr {
		return c.CardanoNodes[newNodeT.NodeAuthData.Uuid]
	}

	if newNodeT, isTr := newNode.(*chia.SaveStatisticRequest); isTr {
		return c.ChiaNodes[newNodeT.NodeAuthData.Uuid]
	}

	return nil
}
