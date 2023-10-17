package ziface

import (
	"net"
)

// 定义连接模块的抽象层
type IConnection interface {
	//启动链接 让当前的链接准备开始工作
	Start()

	//停止链接 结束当前的链接工作
	Stop()

	//获取当前链接绑定的socket conn
	GetTCPConnection() *net.TCPConn

	//获取当前链接模块的链接ID
	GetConnID() uint32

	//获取远程客户端的 TCP状态 Ip 端口
	RemoteAddr() net.Addr

	//发送数据 将数据发送给远程的客户端
	SendMsg(msgID uint32, data []byte) error

	//设置链接属性
	SetProperty(key string, value interface{})

	//获取链接属性
	GetProperty(key string) (interface{}, error)

	//移除链接属性
	RemoveProperty(key string)
}

// 定义一个处理链接业务的方法 (3个参数 链接 发送的数据 数据长度)
type HandleFunc func(*net.TCPConn, []byte, int) error
