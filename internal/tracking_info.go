package internal

import (
	"database/sql"
	"github.com/adarocket/alerter/internal/blockchainrepo"
	"github.com/adarocket/alerter/internal/blockchainrepo/blockchain"
	"github.com/adarocket/alerter/internal/client"
	"github.com/adarocket/alerter/internal/config"
	"google.golang.org/grpc"
	"log"
	"time"
)

var authClient *client.AuthClient

func StartTracking(timeoutCheck int, notifyAddr string, db *sql.DB) {
	conn, _ := auth()

	notifyClient, err := client.NewNotifierClient(notifyAddr)
	if err != nil {
		log.Println(err)
		return
	}

	blockchains := []blockchain.NodesBlockChain{
		&blockchain.Cardano{Blockchain: "cardano"},
	}

	nodeRep := blockchainrepo.InitNodeRepository(notifyClient, blockchains, conn, db)
	for {
		nodeRep.ProcessStatistic()
		time.Sleep(time.Second * time.Duration(timeoutCheck))
	}
}

func auth() (*grpc.ClientConn, error) {
	loadConfig, err := config.LoadConfig() // load config загружвется несколько раз
	if err != nil {
		return nil, err
	}

	clientConn, err := grpc.Dial(loadConfig.ControllerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}
	authClient = client.NewAuthClient(clientConn)

	token, err := authClient.Login(loadConfig.AuthClientLogin, loadConfig.AuthClientPassword)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	conn := setupInterceptorAndClient(token, loadConfig.ControllerAddr)

	return conn, nil
}

func setupInterceptorAndClient(accessToken, serverURL string) *grpc.ClientConn {
	transportOption := grpc.WithInsecure()

	interceptor, err := client.NewAuthInterceptor(authMethods(), accessToken)
	if err != nil {
		log.Fatal("cannot create auth interceptor: ", err)
	}

	clientConn, err := grpc.Dial(serverURL, transportOption, grpc.WithUnaryInterceptor(interceptor.Unary()))
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	return clientConn
}

func authMethods() map[string]bool {
	return map[string]bool{
		"/cardano.Cardano/" + "GetStatistic":  true, //cardano.Cardano
		"/Common.Controller/" + "GetNodeList": true, //Common.Controller
	}
}
