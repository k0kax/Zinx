package znet

import (
	"Zinx/zinx/utils"
	"Zinx/zinx/ziface"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

/*
	7z a zinx-v0.1.7z Zinx
	链接模块
*/

type Connection struct {
	//当前conn属于那个server
	TcpServer ziface.IServer

	//当前链接的socket TCP套接字
	Conn *net.TCPConn

	//当前链接的ID
	connID uint32

	//当前的链接状态
	isClosed bool

	//当前链接所绑定的处理业务的方法API
	//handleAPI ziface.HandleFunc

	//告知当前链接已经退出的/停止 channel(由Reader告知Writer退出)
	ExitChan chan bool

	//无缓冲的管道用于读写Goroutine之间的消息通信
	msgChan chan []byte

	//消息的管理MsgID 和对应的业务处理API关系
	MsgHandler ziface.IMsgHandle

	//链接属性集合
	property map[string]interface{}

	//保护链接属性的锁
	propertyLock sync.RWMutex
}

// 初始化链接的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		connID:     connID,
		MsgHandler: msgHandle,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		property:   make(map[string]interface{}),
	}

	//将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

// 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running...]")
	defer fmt.Println("connID=", c.connID, "[Reader is exit],remote add is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buf中，//最大512字节
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err", err)
		//	continue
		//}

		//创建拆包解包的对象
		dp := NewDataPack()

		//读取客户端的Msg head 二进制流 8字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read message head err:", err)
			break
		}
		//拆包，得到msgid msglen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack err:", err)
			break
		}

		//根据datalen 再次读取data 放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error:", err)
			}
		}
		msg.SetData(data)

		////调当前链接绑定的HandleApi
		//if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
		//	fmt.Println("ConnID", c.ConnID, "handle is error", err)
		//	break
		//}

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启工作池机制，将消息发送给Worker工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//从路由找到注册绑定的conn对应的router调用
			//根据啊绑定好的MsgID找到对应的处理api业务 执行
			go c.MsgHandler.DoMsgHandle(&req)
		}
	}
}

/*
写消息的Goroutine,专门发送客户端消息的模块
*/
func (c *Connection) StartWriter() {
	fmt.Println("[Write Goroutine is running...]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit. remote add is]", c.RemoteAddr().String())

	//不断阻塞等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error,", err)
				return
			}
		case <-c.ExitChan:
			//代表Reader已经退出了，此时Writer也要退出
			return
		}
	}
}

// 启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()...ConnID=", c.connID)
	//启动从当前来连接的读业务
	go c.StartReader()
	//TODO 启动从当前链接写数据的业务
	go c.StartWriter()

	//按照开发者传递进来的，创建链接之后需要调用的处理业务，执行对应的HOOK函数
	c.TcpServer.CallOnConnStart(c)

}

// 停止链接 结束当前的链接工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()=", c.connID)

	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//调用开发者注册的 在销毁链接前 需要执行的业务HOOK函数
	c.TcpServer.CallOnConnStop(c)

	//关闭socket链接
	c.Conn.Close()

	//告知Writer关闭
	c.ExitChan <- true

	//将当前连接从ConnMsg中删除掉
	c.TcpServer.GetConnMgr().Remove(c)

	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

// 获取当前链接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn

}

// 获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.connID

}

// 获取远程客户端的 TCP状态 Ip 端口
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 提供一个SendMsg方法 将我们发个客户端的数据先封包，后发送
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	//将data进行封包 MsgDataLen MsgID/Data

	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgID)
		return errors.New("Pack error msg")
	}

	//将数据发送给客户端
	c.msgChan <- binaryMsg

	return nil
}

// 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	//添加一个链接的属性
	c.property[key] = value
}

// 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

// 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	//删除属性
	delete(c.property, key)
}
