package nodesinfo

import (
	"github.com/adarocket/alerter/client"
	"github.com/adarocket/alerter/inform"
	"github.com/adarocket/proto/proto-gen/cardano"
	"google.golang.org/grpc"
	"log"
)

var serverURL = "165.22.92.139:5300"
var informClient *client.ControllerClient
var authClient *client.AuthClient
var cardanoClient *client.CardanoClient
var chiaClient *client.ChiaClient

const timeout = 15

func StartTracking() {
	if err := auth(); err != nil {
		return
	}

	a, _ := GetNewCardanoNodes()
	for _, request := range a {
		inform.CheckNodes(request)
	}

	/*for _ = range time.Tick(time.Second * timeout) {

	}*/

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
	const informerServicePath = "/cardano.Cardano/"

	return map[string]bool{
		informerServicePath + "GetStatistic": true,
		informerServicePath + "GetNodeList":  true,
	}
}

func GetNewCardanoNodes() (map[string]*cardano.SaveStatisticRequest, error) {
	resp, err := informClient.GetNodeList()
	if err != nil {
		log.Println(err)
		return map[string]*cardano.SaveStatisticRequest{}, err
	}
	//Statistics.
	cardanoNodes := make(map[string]*cardano.SaveStatisticRequest, 10)

	for _, node := range resp.NodeAuthData {
		switch node.Blockchain {
		default:
			resp, err := cardanoClient.GetStatistic("08e792fd-2a19-466f-9a2a-d9fd40bdf9d1")
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
