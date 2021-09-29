package config

import (
	"github.com/bykovme/goconfig"
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
}

const cConfigPath = "alerter.conf"

// var loadedConfig Config

func LoadConfig() (loadedConfig Config, err error) {
	usrHomePath, err := goconfig.GetUserHomePath()
	if err == nil {
		err = goconfig.LoadConfig(usrHomePath+cConfigPath, &loadedConfig)
		if err != nil {
			return loadedConfig, err
		}
	} else {
		return loadedConfig, err
	}
	return loadedConfig, err
}
