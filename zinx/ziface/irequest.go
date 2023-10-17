package ziface

/*
	Irequest接口：
	实际上是把客户端请求和请求的数据包装到了一个request中
*/

type IRequest interface {
	//得到当前的链接
	GetConnection() IConnection

	//得到请求的消息数据\
	GetData() []byte

	//得到请求消息包的id
	GetMsgID() uint32
}
