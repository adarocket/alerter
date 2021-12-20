package auth

import (
	"github.com/adarocket/alerter/internal/client"
	"github.com/adarocket/alerter/internal/config"
	"google.golang.org/grpc"
	"log"
)

func Auth(loadConfig config.Config) (*grpc.ClientConn, error) {
	clientConn, err := grpc.Dial(loadConfig.ControllerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}
	authClient := client.NewAuthClient(clientConn)

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
