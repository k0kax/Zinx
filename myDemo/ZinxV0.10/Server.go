package main

import (
	"Zinx/zinx/ziface"
	"Zinx/zinx/znet"
	"fmt"
)

/*
  基于zinx框架开发的服务端程序
*/

//打包压缩7z a zinx-v0.7.7z Zinx

// ping test 自定义路由
type PingRouetr struct {
	znet.BaseRouter
}

// Test PreHandle
//func (this *PingRouetr) PreHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PreHandle...")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("Before ping...\n"))
//	if err != nil {
//		fmt.Println("call back Before ping err")
//	}
//}

// 在处理conn业务的主方法Hook
func (this *PingRouetr) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle...")
	//先读取客户端的数据，然后ping...ping...ping
	fmt.Println("recv from client :msgId=", request.GetMsgID(), ",data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloRouetr struct {
	znet.BaseRouter
}

// 在处理conn业务的主方法Hook
func (this *HelloRouetr) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloRouter Handle...")
	//先读取客户端的数据，然后ping...ping...ping
	fmt.Println("recv from client :msgId=", request.GetMsgID(), ",data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("你好，我是丁真..."))
	if err != nil {
		fmt.Println(err)
	}

}

//创建链接之后的钩子函数

func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("=====>DoConnectionBegin is Called ....")
	if err := conn.SendMsg(202, []byte("DoConnection BEGIN")); err != nil {
		fmt.Println(err)
	}

	//给当前的链接设置属性
	fmt.Println("Set conn property......")
	conn.SetProperty("Name", "koka")
	conn.SetProperty("Github", "www.github.com")
	conn.SetProperty("Blog", "www.gitbook.com")

}

// 链接断开前需要执行的函数
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("======>DoConnectionLost is Called....")
	fmt.Println("connID = ", conn.GetConnID(), "is Lost...")

	//获取链接属性

	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name:", name)
	}

	if Github, err := conn.GetProperty("Github"); err == nil {
		fmt.Println("Github:", Github)
	}

	if Blog, err := conn.GetProperty("Blog"); err == nil {
		fmt.Println("Blog:", Blog)
	}

}

func main() {
	//1.创建一个server句柄，只用zinx的api

	s := znet.NewServer("[zinx V0.1]")

	//2.注册链接hook钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//3.给当前的zinx框架添加一个自定义的router
	s.AddRouter(0, &PingRouetr{})
	s.AddRouter(1, &HelloRouetr{})

	//4.启动server
	s.Serve()
}
