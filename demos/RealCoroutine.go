package demos

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type TestServiceImpl struct {
	Hello func(ctx context.Context, str string) (string, error)
}

var testServiceImpl = new(TestServiceImpl)

func RealCoroutine() {
	//让协程只由1个线程调度
	runtime.GOMAXPROCS(1)
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

	//获取开始时间戳
	start := time.Now().UnixMilli()
	var loop = 900
	var wg sync.WaitGroup
	wg.Add(loop)
	for i := 0; i < loop; i++ {
		i := i
		go func() {
			defer func() {
				// 必须要先声明defer，否则不能捕获到panic异常
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			defer wg.Done()
			reply, err := testServiceImpl.Hello(context.Background(), "test"+strconv.Itoa(i))
			if err != nil {
				panic(err)
			}
			fmt.Printf("client response result: %v\n", reply)
		}()
	}
	wg.Wait()
	cost := time.Now().UnixMilli() - start
	fmt.Println("go coroutine cost " + strconv.Itoa(int(cost)) + " ms")
}
