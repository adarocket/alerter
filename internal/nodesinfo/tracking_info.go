package nodesinfo

import (
	"github.com/adarocket/alerter/internal/cache"
	"github.com/adarocket/alerter/internal/client"
	"github.com/adarocket/alerter/internal/config"
	"github.com/adarocket/alerter/internal/controller"
	"github.com/adarocket/alerter/internal/nodesinfo/msgsender"
	"github.com/adarocket/proto/proto-gen/cardano"
	"github.com/adarocket/proto/proto-gen/notifier"
	"log"
	"time"

	"google.golang.org/grpc"
)

var informClient *client.ControllerClient
var authClient *client.AuthClient
var cardanoClient *client.CardanoClient

func StartTracking(timeoutCheck int, notifyAddr string) {
	notifyClient, err := client.NewNotifierClient(notifyAddr)
	if err != nil {
		log.Println(err)
		return
	}

	msgSender := msgsender.CreateMsgSender(notifyClient)

	for _ = range time.Tick(time.Second * time.Duration(timeoutCheck)) {
		if err := auth(); err != nil {
			if err := notifyClient.SendMessage(&notifier.SendNotifier{
				TextMessage: "controller down"}); err != nil {
				log.Println(err)
			}
			continue
		}

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

func getNodesMessages() (map[msgsender.KeyMsgSender]msgsender.ValueMsgSender, error) {
	resp, err := informClient.GetNodeList()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	cacheInstance := cache.GetCacheInstance()
	cardanoNodes := make(map[cache.KeyCache]*cardano.SaveStatisticRequest)
	nodesMessages := make(map[msgsender.KeyMsgSender]msgsender.ValueMsgSender)

	for _, node := range resp.NodeAuthData {
		switch node.Blockchain {
		case "cardano":
			key := cache.KeyCache{
				Key:      node.Uuid,
				TypeNode: node.Blockchain,
			}

			fieldNodeMessages, cardanoNode, err := getCardanoNodeMessages(key)
			if err != nil {
				log.Println(err)
				continue
			}

			for key, val := range fieldNodeMessages {
				nodesMessages[key] = val
			}

			cardanoNodes[key] = cardanoNode
		}
	}

	cacheInstance.AddNewInform(cardanoNodes)
	return nodesMessages, nil
}

func getCardanoNodeMessages(key cache.KeyCache) (
	fieldNodeMessages map[msgsender.KeyMsgSender]msgsender.ValueMsgSender,
	resp *cardano.SaveStatisticRequest, err error) {

	resp, err = cardanoClient.GetStatistic(key.Key)
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

	fieldNodeMessages, err = CheckFieldsOfNode(resp, key, alerts)
	if err != nil {
		log.Println(err)
		return
	}

	return fieldNodeMessages, resp, nil
}
