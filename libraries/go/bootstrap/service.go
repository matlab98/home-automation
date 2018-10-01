package bootstrap

import (
	"fmt"
	"github.com/go-redis/redis"
	"home-automation/libraries/go/client"
	"home-automation/libraries/go/config"
	"log"
	"os"
)

type service struct{
	ControllerName string
	Config *config.Config
	Redis *redis.Client
}

func NewService(controllerName string) (*service, error) {
	service := &service{
		ControllerName: controllerName,
	}

	// Create API Client
	apiGateway := os.Getenv("API_GATEWAY")
	if apiGateway == "" {
		return nil, fmt.Errorf("API_GATEWAY env var not set")
	}
	apiClient, err := client.New(apiGateway)
	if err != nil {
		return nil, err
	}

	// Load config
	var configRsp map[string]interface{}
	_, err = apiClient.Get(fmt.Sprintf("service.config/read/%s", controllerName), &configRsp)
	if err != nil {
		return nil, err
	}
	service.Config = &config.Config{
		Map: configRsp,
	}

	// Connect to Redis
	if service.Config.Has("redis.host") {
		host := service.Config.Get("redis.host").String()
		port := service.Config.Get("redis.port").Int()
		addr := fmt.Sprintf("%s:%d", host, port)
		log.Printf("Connecting to Redis at address %s\n", addr)
		service.Redis = redis.NewClient(&redis.Options{
			Addr: addr,
			Password: "",
			DB: 0,
		})
	}

	return service, nil
}
