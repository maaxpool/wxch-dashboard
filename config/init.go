package config

import (
	"github.com/spf13/viper"
	"strings"
)

var config Config = Config{
	Debug: &DebugConfig{
		Verbose:       false,
		DisableCron:   true,
		DisableSentry: true,
		SentryDSN:     "",
	},
	DB: &DBConfig{
		Host:     "127.0.0.1",
		User:     "postgres",
		Password: "password",
		DbName:   "wxch-dashboard",
		Port:     "15432",
		SSLMode:  "disable",
		TimeZone: "UTC",
	},
	RPC: &RPCConfig{
		Listen: "127.0.0.1",
		Port:   4700,
		AdminToken: "",
	},
	EthNode: &EthNodeConfig{
		ChainId:                   3,
		InfuraHost:                "",
		WxchBridgeContractAddress: "0xbaC8fa980A71Ff221D361905999654319d46202D",
	},
}

func init() {
	instance := viper.New()

	// only for dev
	instance.AddConfigPath("/etc/wxch-dashboard")
	instance.AddConfigPath(".")

	instance.SetConfigType("yaml")
	instance.SetConfigName("config.yaml")

	instance.SetEnvPrefix("wd")
	instance.AutomaticEnv()
	instance.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := instance.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}

	err = instance.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
}

func Get() *Config {
	return &config
}
