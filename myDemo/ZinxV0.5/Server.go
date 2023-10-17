package main

import (
	"Zinx/zinx/ziface"
	"Zinx/zinx/znet"
	"fmt"
)

/*
  基于zinx框架开发的服务端程序
*/

//打包压缩7z a zinx-v0.4.7z Zinx

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
	fmt.Println("Call Router Handle...")
	//先读取客户端的数据，然后ping...ping...ping
	fmt.Println("recv from client :msgId=", request.GetMsgID(), ",data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

//// 在处理conn业务的钩子方法Hook
//func (this *PingRouetr) PostHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PreHandle...")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping\n"))
//	if err != nil {
//		fmt.Println("call back After ping err")
//	}
//}

func main() {
	//1.创建一个server句柄，只用zinx的api

	s := znet.NewServer("[zinx V0.1]")

	//2.给当前的zinx框架添加一个自定义的router
	s.AddRouter(&PingRouetr{})

	//3.启动server
	s.Serve()
}
