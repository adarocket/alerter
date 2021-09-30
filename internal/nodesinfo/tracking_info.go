package nodesinfo

import (
	"github.com/adarocket/alerter/internal/cache"
	client2 "github.com/adarocket/alerter/internal/client"
	"github.com/adarocket/alerter/internal/config"
	"log"
	"time"

	pb "github.com/adarocket/alerter/proto"
	"google.golang.org/grpc"
)

var informClient *client2.ControllerClient
var authClient *client2.AuthClient
var cardanoClient *client2.CardanoClient
var chiaClient *client2.ChiaClient

const timeout = 15

func StartTracking() {
	notifyClient, err := client2.NewNotifierClient()
	if err != nil {
		log.Println(err)
		return
	}

	for {
		if err := auth(); err == nil {
			break
		}
		if err := notifyClient.SendMessage(&pb.SendNotifier{
			TypeMessage: "controller down", Value: err.Error()}); err != nil {
			log.Println(err)
		}
		time.Sleep(time.Second * 5)
	}

	cacheInstance := cache.GetCacheInstance()
	for _ = range time.Tick(time.Second * timeout) {
		nodes, err := GetNodes()
		if err != nil {
			log.Println(err)
			if err := notifyClient.SendMessage(&pb.SendNotifier{
				TypeMessage: "cant get nodes", Value: err.Error()}); err != nil {
				log.Println(err)
			}
			time.Sleep(time.Second * 5)
			continue
		}

		var messages []*pb.SendNotifier
		for key, node := range nodes {
			messages, err = CheckFieldsOfNode(node, key)
			if err != nil {
				log.Println(err)
			}
			err = notifyClient.SendMessages(messages)
			if err != nil {
				log.Println(err)
			}
		}

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
	authClient = client2.NewAuthClient(clientConn)

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

	interceptor, err := client2.NewAuthInterceptor(authMethods(), accessToken)
	if err != nil {
		log.Fatal("cannot create auth interceptor: ", err)
	}

	clientConn, err := grpc.Dial(serverURL, transportOption, grpc.WithUnaryInterceptor(interceptor.Unary()))
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	informClient = client2.NewControllerClient(clientConn)
	cardanoClient = client2.NewCardanoClient(clientConn)
	chiaClient = client2.NewChiaClient(clientConn)
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
