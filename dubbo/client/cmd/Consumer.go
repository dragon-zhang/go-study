package main

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"fmt"
)

type TestServiceImpl struct {
	Hello func(ctx context.Context, str string) (string, error)
}

var testServiceImpl = new(TestServiceImpl)

func main() {
	config.SetConsumerService(testServiceImpl)

	rootConfig := config.NewRootConfigBuilder().
		SetApplication(config.NewApplicationConfigBuilder().
			SetName("consumer").
			Build()).
		AddRegistry("zk", config.NewRegistryConfigBuilder().
			SetAddress("zookeeper://localhost:2181").
			Build()).
		SetConsumer(config.NewConsumerConfigBuilder().
			AddReference("TestServiceImpl", config.NewReferenceConfigBuilder().
				SetInterface("com.example.demo.TestService").
				SetProtocol("dubbo").
				SetRegistryIDs("zk").
				Build()).
			SetRequestTimeout("3s").
			Build()).
		Build()

	if err := config.Load(config.WithRootConfig(rootConfig)); err != nil {
		panic(err)
	}

	reply, err := testServiceImpl.Hello(context.Background(), "test")
	if err != nil {
		panic(err)
	}
	fmt.Printf("client response result: %v\n", reply)
}
