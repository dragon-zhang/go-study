package main

import (
	"context"
	dubbo "dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"fmt"
	triple "github.com/dubbogo/triple/pkg/common/constant"
	"go-projects/dubbo-proto/api"
)

var demoService = new(api.DemoServiceClientImpl)

func main() {
	config.SetConsumerService(demoService)

	rootConfig := config.NewRootConfigBuilder().
		SetApplication(config.NewApplicationConfigBuilder().
			SetName("consumer").
			Build()).
		AddRegistry("zk", config.NewRegistryConfigBuilder().
			SetAddress("zookeeper://localhost:2181").
			Build()).
		SetConsumer(config.NewConsumerConfigBuilder().
			AddReference("DemoServiceClientImpl", config.NewReferenceConfigBuilder().
				SetInterface("org.apache.dubbo.demo.DemoService").
				SetProtocol(triple.TRIPLE).
				SetRegistryIDs("zk").
				SetSerialization(dubbo.ProtobufSerialization).
				Build()).
			SetRequestTimeout("3s").
			Build()).
		Build()

	if err := config.Load(config.WithRootConfig(rootConfig)); err != nil {
		panic(err)
	}

	reply, err := demoService.SayHello(context.Background(), &api.HelloRequest{Name: "consumer"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("client response result: %v\n", reply)
}
