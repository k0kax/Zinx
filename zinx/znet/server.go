package znet

import (
	"Zinx/zinx/utils"
	"Zinx/zinx/ziface"
	"fmt"
	"net"
	_ "runtime"
)

//实现层 类模块

// IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	//服务器名称
	Name string

	//服务器绑定的ip
	IPVersion string

	//服务器监听的IP
	IP string

	//服务器监听的端口
	Port int

	//当前server的管理模块，用来绑定当前msgID和对应的业务处理Api
	MsgHandler ziface.IMsgHandle

	//当前server的链接管理器
	ConnMgr ziface.IConnManager

	//该server创建链接之后自动调hook函数-OnConnStart
	OnConnStart func(conn ziface.IConnection)
	//该server销毁链接之后自动调hook函数-OnConnStop
	OnConnStop func(conn ziface.IConnection)
}

//// 定义当前客户端链接所绑定的handle api (目前handle写死，日后用户自定义)
//func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
//	//回显业务
//
//	fmt.Println("[Conn Handle] CallbackToClient ...")
//	if _, err := conn.Write(data[:cnt]); err != nil {
//		fmt.Println("write back buf err", err)
//		return errors.New("CallBackToClient error")
//	}
//
//	return nil
//}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s,listenner at IP:%s,Port:%d is starting\n",
		utils.GlobalObject.Name,
		utils.GlobalObject.Host,
		utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version:%s,MaxConn:%d,MaxPackageSize:%d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	fmt.Printf("[Start] Server Listenner at IP:%s,Port %d,is starting\n", s.IP, s.Port)

	go func() {
		//0.开启消息队列及worker工作池
		s.MsgHandler.StartWorkerPool()

		//1.获得一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}

		//2.监听服务器地址
		listrenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
		}

		//监听成功
		fmt.Println("start Zinx server success", s.Name, "success,Listenning")
		var cid uint32
		cid = 0

		//3.阻塞的等待用户链接，处理客户端链接业务
		for {
			//如果有客户端链接过来，阻塞会返回
			conn, err := listrenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			//设置最大连接个数的判断，如果超过最大连接则关闭新的链接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//TODO 给客户端响应一个超出最大连接的错误包
				fmt.Println(" ===========>Too Many Connection MaxConn<=================", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}
			//将处理新连接的业务方法和conn进行绑定，得到我们的链接模块
			//delConn := NewConnection(s, conn, cid, s.MsgHandler)
			delConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			//启动当前的链接业务处理通道
			go delConn.Start()

			//v0.9版本之前的没有去除此回显机制

			//客户端已经建立连接，做一些业务，做一个基本的最大512字节的回显业务
			//go func() {
			//	for {
			//		buf := make([]byte, 512)
			//		cnt, err := conn.Read(buf)
			//		if err != nil {
			//			fmt.Println("recv buf err", err)
			//		}
			//
			//		fmt.Printf("recv client buf %s,cnt %d\n", buf, cnt)
			//		//回显功能
			//		if _, err := conn.Write(buf[:cnt]); err != nil {
			//			fmt.Println("write back buf err", err)
			//		}
			//	}
			//}()
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	//TODO 将一些服务器的资源、状态或者一些已经开辟的链接信息 进行停止或者回收
	//清空链接，回收资源
	fmt.Println("[STOP]Zinx server name:", s.Name)
	s.ConnMgr.ClearConn()
}

// 运行服务器
func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//TODO 做一些启动服务器外的业务

	//阻塞状态
	select {}
}

// 添加路由:给当前的服务注册一个路由方法，供客户端的链接处理使用
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Succ")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

/*
初始化Server的模块
*/

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}

	return s
}

// 注册OnConnStart钩子方法
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 注册OnConnStop钩子方法
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用OnConnStart钩子方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("------->Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

// 调用OnConnStop钩子方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("------->Call OnConnStop()...")
		s.OnConnStop(conn)
	}
}
