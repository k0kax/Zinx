package utils

import (
	"Zinx/zinx/ziface"
	"encoding/json"
	"io/ioutil"
)

/*
定义一个储存一切Zinx框架的全局变量，供其他模块使用
一切参数由通过zinx.json由用户进行配置
*/
type GlobalObj struct {

	/*
		server
	*/
	TcpServer ziface.IServer //当前zinx全局的server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int            //当前服务器主机监听的端口
	Name      string         //当前服务器的名称

	/*
		Zinx
	*/

	Version          string //当前zinx的版本号
	MaxConn          int    //当前服务器主机允许的最大连接数
	MaxPackageSize   uint32 //当前zinx框架数据包的最大值
	WorkerPoolSize   uint32 //当前业务工作worker池的Goroutine数量
	MaxWorkerTaskLen uint32 //zinx框架允许用户开辟多少个worker(限定条件)

}

// 定义一个全局的对外对象GloableObj
var GlobalObject *GlobalObj

// 加载用户自定义的zinx.json的方法
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}

	//将json文件解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 初始化当前的GlobalObject
func init() {
	//如果配置文件没有加载，默认的值
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "V0.9",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        10,
		MaxPackageSize: 4096,

		WorkerPoolSize:   10,   //worker工作池的队列的大小个数
		MaxWorkerTaskLen: 1024, //每个worker对应的消息队列的任务的数量最大值
	}

	//尝试从conf/zinx.json加载一些用户自定义的参数
	GlobalObject.Reload()
}
