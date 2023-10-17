package main

import (
	"Zinx/zinx/znet"
	"fmt"
	"io"
	"net"
	"time"
)

/*模拟客户端*/
func main() {
	//1.直接远程连接服务器
	fmt.Println("Client start...")
	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err,exit")
		return
	}

	for {
		//发送封包的message消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("测试信息")))
		if err != nil {
			fmt.Println("Pack error:", err)
			return
		}

		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("Write error:", err)
			return
		}

		//1.从服务器接受数据 message id=0 ping...ping...ping

		//1.先读取流中的head,得到Id dataLen

		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error:", err)
			break
		}
		//将二进制的head拆包到msg结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error", err)
		}

		if msgHead.GetMsgLen() > 0 {
			//msg里有数据
			//2.再根据dataLen进行第二次读取，将data读取出来
			msg := msgHead.(*znet.Message) //类型转换
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error:", err)
				return
			}

			fmt.Println("------->Recv Server Msg :ID=", msg.Id, ",Len=", msg.GetMsgLen(), ",data=", string(msg.GetData()))
		}

		//cpu阻塞，防止无限循环
		time.Sleep(1 * time.Second)

	}

}
