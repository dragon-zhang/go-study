package main

import (
	"context"
	dubbo "dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"github.com/dragon-zhang/go-study/dubbo-proto/api"
	triple "github.com/dubbogo/triple/pkg/common/constant"
)

type DemoServiceImpl struct {
	api.UnimplementedDemoServiceServer
}

func (u *DemoServiceImpl) SayHello(context context.Context, request *api.HelloRequest) (*api.HelloReply, error) {
	return &api.HelloReply{Message: "Hello " + request.GetName() + ", response from provider"}, nil
}

func (u *DemoServiceImpl) Reference() string {
	return "demoServiceImpl"
}

func main() {
	config.SetProviderService(&DemoServiceImpl{})

	rootConfig := config.NewRootConfigBuilder().
		SetApplication(config.NewApplicationConfigBuilder().
			SetName("provider").
			Build()).
		AddRegistry("zk", config.NewRegistryConfigBuilder().
			SetAddress("zookeeper://localhost:2181").
			Build()).
		SetProvider(config.NewProviderConfigBuilder().
			AddService("demoServiceImpl", config.NewServiceConfigBuilder().
				SetInterface("org.apache.dubbo.demo.DemoService").
				SetProtocolIDs(triple.TRIPLE).
				SetSerialization(dubbo.ProtobufSerialization).
				Build()).
			Build()).
		AddProtocol(triple.TRIPLE, config.NewProtocolConfigBuilder().
			SetName(triple.TRIPLE).
			SetPort("20000").
			Build()).
		Build()

	if err := config.Load(config.WithRootConfig(rootConfig)); err != nil {
		panic(err)
	}

	select {}
}
