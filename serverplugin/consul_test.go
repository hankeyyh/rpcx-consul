package serverplugin

import (
	"context"
	"fmt"
	metrics "github.com/rcrowley/go-metrics"
	rpcx_consul_client "github.com/rpcxio/rpcx-consul/client"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
	"testing"
	"time"
)

func TestConsulRegistry(t *testing.T) {
	s := server.NewServer()

	r := NewConsulRegisterPlugin(
		WithConsulServiceAddress("tcp@127.0.0.1:8972"),
		WithConsulServers([]string{"127.0.0.1:8500"}),
		WithConsulBasePath("/rpcx_test"),
		WithConsulMetrics(metrics.NewRegistry()),
		WithConsulUpdateInterval(time.Minute),
	)
	err := r.Start()
	if err != nil {
		return
	}
	defer func() {
		s.Close()
		r.Stop()
	}()
	s.Plugins.Add(r)

	s.RegisterName("Arith", new(Arith), "")
	s.Serve("tcp", "127.0.0.1:8972")

	//if err := r.Stop(); err != nil {
	//	t.Fatal(err)
	//}
}

func TestConsulDiscovery(t *testing.T) {
	d, _ := rpcx_consul_client.NewConsulDiscovery("/rpcx_test", "Arith", []string{"127.0.0.1:8500"}, nil)
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	args := Args{
		A: 1,
		B: 2,
	}
	reply := new(Reply)
	for {
		err := xclient.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(reply.C)
		time.Sleep(time.Second)
	}
}
