package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var configInfo *ConfigInfo

type ConfigInfo struct {
	ApiServiceConfig   *ApiServiceConfig
}

func NewConfig() *ConfigInfo {
	return &ConfigInfo{
		ApiServiceConfig: &ApiServiceConfig{
			GrpcHost: viper.GetString("grpc.api_service.host"),
		},
	}
}

func GetConfig() *ConfigInfo {
	if configInfo == nil {
		fmt.Println("config is null")
		configInfo = NewConfig()
	}
	return configInfo
}
