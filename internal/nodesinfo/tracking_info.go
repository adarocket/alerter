package nodesinfo

import (
	"github.com/adarocket/alerter/internal/cache"
	"github.com/adarocket/alerter/internal/client"
	"github.com/adarocket/alerter/internal/config"
	"github.com/adarocket/alerter/internal/controller"
	"github.com/adarocket/alerter/internal/nodesinfo/msgsender"
	"github.com/adarocket/proto/proto-gen/cardano"
	"github.com/adarocket/proto/proto-gen/common"
	"log"
	"time"

	"github.com/adarocket/proto/proto-gen/notifier"
	"google.golang.org/grpc"
)

var informClient *client.ControllerClient
var authClient *client.AuthClient
var cardanoClient *client.CardanoClient

func StartTracking(timeoutCheck int) {
	notifyClient, err := client.NewNotifierClient()
	if err != nil {
		log.Println(err)
		return
	}

	msgSender := msgsender.CreateMsgSender(notifyClient)

	for {
		if err := auth(); err == nil {
			break
		}
		if err := notifyClient.SendMessage(&notifier.SendNotifier{
			TypeMessage: "controller down", Value: err.Error()}); err != nil {
			log.Println(err)
		}
		time.Sleep(time.Second * time.Duration(timeoutCheck))
	}

	for _ = range time.Tick(time.Second * time.Duration(timeoutCheck)) {
		nodesMessages, err := getNodesMessages()
		if err != nil {
			log.Println(err)
			continue
		}

		msgSender.AddNotifiersToStack(nodesMessages)
	}
}

func auth() error {
	loadConfig, err := config.LoadConfig()
	if err != nil {
		return err
	}

	clientConn, err := grpc.Dial(loadConfig.ControllerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}
	authClient = client.NewAuthClient(clientConn)

	token, err := authClient.Login(loadConfig.AuthClientLogin, loadConfig.AuthClientPassword)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	setupInterceptorAndClient(token, loadConfig.ControllerAddr)

	return nil
}

func setupInterceptorAndClient(accessToken, serverURL string) {
	transportOption := grpc.WithInsecure()

	interceptor, err := client.NewAuthInterceptor(authMethods(), accessToken)
	if err != nil {
		log.Fatal("cannot create auth interceptor: ", err)
	}

	clientConn, err := grpc.Dial(serverURL, transportOption, grpc.WithUnaryInterceptor(interceptor.Unary()))
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	informClient = client.NewControllerClient(clientConn)
	cardanoClient = client.NewCardanoClient(clientConn)
}

func authMethods() map[string]bool {
	return map[string]bool{
		"/cardano.Cardano/" + "GetStatistic":  true, //cardano.Cardano
		"/Common.Controller/" + "GetNodeList": true, //Common.Controller
	}
}

func getNodesMessages() (map[msgsender.KeyMsgSender]*notifier.SendNotifier, error) {
	resp, err := informClient.GetNodeList()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	cacheInstance := cache.GetCacheInstance()
	cardanoNodes := make(map[cache.KeyCache]interface{})
	nodesMessages := make(map[msgsender.KeyMsgSender]*notifier.SendNotifier)

	for _, node := range resp.NodeAuthData {
		switch node.Blockchain {
		case "cardano":
			fieldNodeMessages, cardanoNode, err := getCardanoNodeMessages(node)
			if err != nil {
				log.Println(err)
				continue
			}

			for key, val := range fieldNodeMessages {
				nodesMessages[key] = val
			}

			key := cache.KeyCache{
				Key:      node.Uuid,
				TypeNode: node.Blockchain,
			}
			cardanoNodes[key] = cardanoNode
		}
	}

	cacheInstance.AddNewInform(cardanoNodes)
	return nodesMessages, nil
}

func getCardanoNodeMessages(node *common.NodeAuthData) (
	fieldNodeMessages map[msgsender.KeyMsgSender]*notifier.SendNotifier, resp *cardano.SaveStatisticRequest, err error) {

	resp, err = cardanoClient.GetStatistic(node.Uuid)
	if err != nil {
		log.Println(err)
		return
	}

	if resp.NodeAuthData.Uuid == "" {
		return
	}

	db := controller.GetAlertNodeControllerInstance()
	alerts, err := db.GetAlertsByNodeUuid(resp.NodeAuthData.Uuid)
	if err != nil {
		log.Println(err)
		return
	}

	if len(alerts) == 0 {
		return
	}

	key := cache.KeyCache{
		Key:      resp.NodeAuthData.Uuid,
		TypeNode: node.Blockchain,
	}
	fieldNodeMessages, err = CheckFieldsOfNode(resp, key, alerts)
	if err != nil {
		log.Println(err)
		return
	}

	return fieldNodeMessages, resp, nil
}
