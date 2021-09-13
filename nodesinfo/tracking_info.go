package nodesinfo

import (
	"github.com/adarocket/alerter/cache"
	"github.com/adarocket/alerter/client"
	"github.com/adarocket/alerter/inform"
	pb "github.com/adarocket/alerter/proto"
	"google.golang.org/grpc"
	"log"
	"time"
)

var serverURL = "165.22.92.139:5300"
var informClient *client.ControllerClient
var authClient *client.AuthClient
var cardanoClient *client.CardanoClient
var chiaClient *client.ChiaClient
var notifyClient *client.NotifierClient

const timeout = 15

func StartTracking() {
	notifyClient, err := client.NewNotifierClient()
	if err != nil {
		log.Println(err)
		return
	}

	for {
		err = auth()
		if err == nil {
			break
		}

		msg := pb.SendNotifier{
			TypeMessage: "controller down",
			Value:       err.Error(),
			Frequency:   "max",
		}
		if err := notifyClient.SendMessage(&msg); err != nil {
			log.Println(err)
		}

		time.Sleep(time.Second * 5)
	}

	cacheInstance := cache.GetCacheInstance()
	for _ = range time.Tick(time.Second * timeout) {
		nodes, err := GetNodes()
		if err != nil {
			log.Println(err)
			msg := pb.SendNotifier{
				TypeMessage: "cant get nodes",
				Value:       err.Error(),
				Frequency:   "max",
			}

			if err := notifyClient.SendMessage(&msg); err != nil {
				log.Println(err)
			}
			continue
		}

		var msges []*pb.SendNotifier
		for _, SendNotifier := range nodes {
			msges, err = inform.CheckFieldsOfNode(SendNotifier)
			if err != nil {
				log.Println(err)
			}

			err = notifyClient.SendMessages(msges)
			if err != nil {
				log.Println(err)
			}
		}

		cacheInstance.AddNewInform(nodes)
	}

	return
}

func auth() error {
	clientConn, err := grpc.Dial(serverURL, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}
	authClient = client.NewAuthClient(clientConn)

	token, err := authClient.Login("admin1", "secret")
	if err != nil {
		log.Println(err.Error())
		return err
	}

	setupInterceptorAndClient(token)

	return nil
}

func setupInterceptorAndClient(accessToken string) {
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
	chiaClient = client.NewChiaClient(clientConn)
}

func authMethods() map[string]bool {
	//const informerServicePath = "/Common.Controller/"

	return map[string]bool{
		"/cardano.Cardano/" + "GetStatistic":  true, //cardano.Cardano
		"/Common.Controller/" + "GetNodeList": true, //Common.Controller
	}
}

func GetNodes() (map[string]interface{}, error) {
	resp, err := informClient.GetNodeList()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	cardanoNodes := make(map[string]interface{}, 10)

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

			cardanoNodes[resp.NodeAuthData.Uuid] = resp
		}
	}

	return cardanoNodes, nil
}
