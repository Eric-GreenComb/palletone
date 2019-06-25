package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/Eric-GreenComb/palletone/bean"
)

// Server Server Config
var Server bean.ServerConfig

const cmdRoot = "core"

func init() {
	loadConfig()
}

func loadConfig() {
	viper.SetEnvPrefix(cmdRoot)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigName(cmdRoot)
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(fmt.Errorf("Fatal error when reading %s config file:%s", cmdRoot, err))
		os.Exit(1)
	}

	Server.Port = viper.GetString("server.port")
	Server.Mode = viper.GetString("server.mode")
}
