package main

import "Zinx/zinx/znet"

/*
  基于zinx框架开发的服务端程序
*/

func main() {
	//创建一个server句柄，只用zinx的api

	s := znet.NewServer("[zinx V0.1]")

	s.Serve()
}
