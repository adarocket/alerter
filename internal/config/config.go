package config

import (
	"github.com/bykovme/goconfig"
	"log"
	"sync"
)

// Config - structure of config file
type Config struct {
	NotifierAddr       string `json:"notifier_addr"`
	SqlLitePathDB      string `json:"sql_lite_path_db"`
	ControllerAddr     string `json:"controller_addr"`
	WebServerAddr      string `json:"web_server_addr"`
	SecretKey          string `json:"secret_key"`
	AuthClientLogin    string `json:"auth_client_login"`
	AuthClientPassword string `json:"auth_client_password"`
	TimeoutCheck       int    `json:"timeout_check"`
}

const cConfigPath = "alerter.conf"

var once sync.Once

var config Config

func LoadConfig() (loadedConfig Config, err error) {
	once.Do(func() {
		//usrHomePath, err := goconfig.GetUserHomePath()
		if err == nil {
			err = goconfig.LoadConfig(cConfigPath, &loadedConfig)
			if err != nil {
				log.Fatal(err)
			}
			config = loadedConfig
		} else {
			log.Fatal(err)
		}
	})

	return config, err
}
