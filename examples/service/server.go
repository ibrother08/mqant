package main

import (
	"github.com/liangdas/mqant"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/examples/proto/examples/greeter"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	basemodule "github.com/liangdas/mqant/module/base"
	"github.com/liangdas/mqant/registry"
	"github.com/liangdas/mqant/registry/consul"
	"github.com/nats-io/nats.go"
)

type Greeter struct {
}

func (g *Greeter) Hello(in *greeter.Request) (out *greeter.Response, err error) {
	log.LogBeego().Info("xxxxxx %s", "zouwo ")
	out = &greeter.Response{
		Msg: "success",
	}
	return
}
func (g *Greeter) Stream(in *greeter.Request) (out *greeter.Response, err error) {
	return
}
func main() {
	// 服务实例
	app := mqant.CreateApp(
		module.Debug(true),
	)
	// 配置加载
	consulURL := "127.0.0.1:8500"
	// connect to consul
	rs := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{consulURL}
	})

	// 读取一些基础配置
	natsURL := "nats://127.0.0.1:4222"
	natsUser := ""
	natsPassword := ""
	nc, err := nats.Connect(natsURL, nats.ErrorHandler(nil), nats.UserInfo(natsUser, natsPassword), nats.MaxReconnects(10000))
	if err != nil {
		//panic("no connect nats")
		//return
	}
	_ = app.OnConfigurationLoaded(func(app module.App) {
		_ = app.UpdateOptions(
			module.Nats(nc),     //指定nats rpc
			module.Registry(rs), //指定服务发现
		)
	})
	s := &Server{}
	_ = app.Run(s)
}

type Server struct {
	basemodule.BaseModule
	version string
	// 模块名字
	Name string
}

// GetApp module.App
func (m *Server) GetApp() module.App {
	return m.App
}

// OnInit() 初始化配置
func (s *Server) OnInit(app module.App, settings *conf.ModuleSettings) {
	s.BaseModule.OnInit(s, app, settings)
	srv := &Greeter{}
	greeter.RegisterGreeterTcpHandler(&s.BaseModule, srv)
}

// Run() 运行服务
func (s *Server) Run(closeSig chan bool) {
	//创建MongoDB连接实例
}

// 销毁服务
func (s *Server) OnDestroy() {
	//一定别忘了继承
	s.BaseModule.OnDestroy()
	s.GetServer().OnDestroy()
}

// Version() 获取当前服务的代码版本
func (s *Server) Version() string {
	//可以在监控时了解代码版本
	return s.version
}

func (s *Server) GetType() string {
	return "greeter"
}
