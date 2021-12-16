package main

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"fmt"
)

type TestServiceImpl struct {
}

func (u *TestServiceImpl) Hello(_ context.Context, param string) (string, error) {
	fmt.Println("param->" + param)
	return "hello " + param, nil
}

func (u *TestServiceImpl) Reference() string {
	return "testServiceImpl"
}

func main() {
	config.SetProviderService(&TestServiceImpl{})

	rootConfig := config.NewRootConfigBuilder().
		SetApplication(config.NewApplicationConfigBuilder().
			SetName("provider").
			Build()).
		AddRegistry("zk", config.NewRegistryConfigBuilder().
			SetAddress("zookeeper://localhost:2181").
			Build()).
		SetProvider(config.NewProviderConfigBuilder().
			AddService("testServiceImpl", config.NewServiceConfigBuilder().
				SetInterface("com.example.demo.TestService").
				SetProtocolIDs("dubboKey").
				Build()).
			Build()).
		AddProtocol("dubboKey", config.NewProtocolConfigBuilder().
			SetName("dubbo").
			SetPort("20000").
			Build()).
		Build()

	if err := config.Load(config.WithRootConfig(rootConfig)); err != nil {
		panic(err)
	}

	select {}
}
