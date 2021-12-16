package nodesinfo

import (
	"github.com/adarocket/alerter/internal/client"
	"github.com/adarocket/alerter/internal/config"
	"github.com/adarocket/alerter/internal/database"
	"github.com/adarocket/alerter/internal/noderepo"
	"github.com/adarocket/alerter/internal/structs"
	"google.golang.org/grpc"
	"log"
	"time"
)

var authClient *client.AuthClient

func StartTracking(timeoutCheck int, notifyAddr string, db database.ModelAlertNode) {
	conn, _ := auth()

	notifyClient, err := client.NewNotifierClient(notifyAddr)
	if err != nil {
		log.Println(err)
		return
	}

	var blockchains []structs.NodesBlockChain
	cardanoStruct := structs.Cardano{Blockchain: "cardano"}
	blockchains = append(blockchains, &cardanoStruct)

	nodeRep := noderepo.InitNodeRepository(notifyClient, blockchains, conn, db)

	for {
		nodeRep.ProcessStatistic()
		time.Sleep(time.Hour * time.Duration(timeoutCheck))
	}
}

func auth() (*grpc.ClientConn, error) {
	loadConfig, err := config.LoadConfig()
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
