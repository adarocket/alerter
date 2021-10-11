package nodesinfo

import (
	"github.com/adarocket/alerter/internal/cache"
	"github.com/adarocket/alerter/internal/client"
	"github.com/adarocket/alerter/internal/config"
	"github.com/adarocket/alerter/internal/nodesinfo/msgsender"
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

	cacheInstance := cache.GetCacheInstance()
	for _ = range time.Tick(time.Second * time.Duration(timeoutCheck)) {
		nodes, err := GetNodes()
		if err != nil {
			log.Println(err)
			errSend := notifyClient.SendMessage(&notifier.SendNotifier{
				TypeMessage: "cant get nodes", Value: err.Error()})
			if errSend != nil {
				log.Println(errSend)
			}
			continue
		}

		var messages map[msgsender.KeyMsgSender]*notifier.SendNotifier
		for key, node := range nodes {
			messages, err = CheckFieldsOfNode(node, key)
			if err != nil {
				log.Println(err)
			}
		}

		msgSender.AddNotifiersToStack(messages)
		cacheInstance.AddNewInform(nodes)
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

func GetNodes() (map[cache.KeyCache]interface{}, error) {
	resp, err := informClient.GetNodeList()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	cardanoNodes := make(map[cache.KeyCache]interface{}, 10)

	for _, node := range resp.NodeAuthData {
		switch node.Blockchain {
		case "cardano":
			resp, err := cardanoClient.GetStatistic(node.Uuid)
			if err != nil {
				log.Println(err)
				continue
			}

			if resp.NodeAuthData.Uuid == "" {
				continue
			}

			key := cache.KeyCache{
				Key:      resp.NodeAuthData.Uuid,
				TypeNode: node.Blockchain,
			}
			cardanoNodes[key] = resp
		}
	}

	return cardanoNodes, nil
}
