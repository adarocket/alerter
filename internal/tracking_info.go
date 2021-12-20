package internal

import (
	"database/sql"
	"github.com/adarocket/alerter/internal/auth"
	"github.com/adarocket/alerter/internal/blockchainrepo"
	"github.com/adarocket/alerter/internal/blockchainrepo/blockchain"
	"github.com/adarocket/alerter/internal/client"
	"github.com/adarocket/alerter/internal/config"
	"log"
	"time"
)

func StartTracking(loadConfig config.Config, db *sql.DB) {
	notifyClient, err := client.NewNotifierClient(loadConfig.NotifierAddr)
	if err != nil {
		log.Println(err)
		return
	}

	blockchains := []blockchain.NodesBlockChain{
		&blockchain.Cardano{Blockchain: "cardano"},
	}

	nodeRep := blockchainrepo.InitNodeRepository(notifyClient, blockchains, auth.Auth, loadConfig, db)
	for {
		nodeRep.ConnectNodeRepositoryServices(loadConfig)
		nodeRep.ProcessStatistic()
		time.Sleep(time.Second * time.Duration(loadConfig.TimeoutCheck))
	}
}
