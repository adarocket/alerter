package cache

import (
	"sync"

	"github.com/adarocket/proto/proto-gen/cardano"
	"github.com/adarocket/proto/proto-gen/chia"
)

// FIXME: переделать кэширование
// использовать https://github.com/dgraph-io/ristretto

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

type KeyCache struct {
	Key      string
	TypeNode string
}

type Cache struct {
	CardanoNodes map[KeyCache]*cardano.SaveStatisticRequest
	ChiaNodes    map[KeyCache]*chia.SaveStatisticRequest
}

func (c *Cache) AddNewInform(newMap interface{}) {
	if newMapCardano, isTr := newMap.(map[KeyCache]*cardano.SaveStatisticRequest); isTr {
		for key, node := range newMapCardano {
			c.CardanoNodes[key] = node
		}
	}

	if newMapChia, isTr := newMap.(map[KeyCache]*chia.SaveStatisticRequest); isTr {
		for key, node := range newMapChia {
			c.ChiaNodes[key] = node
		}
	}
	// add error if type not exist
}

func (c *Cache) GetOldNodeByType(newNode interface{}, key KeyCache) (oldNode interface{}) {
	if _, isTr := newNode.(*cardano.SaveStatisticRequest); isTr {
		return c.CardanoNodes[key]
	}

	if _, isTr := newNode.(*chia.SaveStatisticRequest); isTr {
		return c.ChiaNodes[key]
	}

	return nil
}