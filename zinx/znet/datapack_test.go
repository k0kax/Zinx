package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//测试模块

//运行该测试代码需要，先将globalobj.go的初始化配置代码关闭

// 测试datapack封包，拆包的单元测试
func TestDataPack_Pack(t *testing.T) {
	/*
		模拟的服务器
	*/

	//1.创建socketTCP
	listenner, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	//创建一个go承载 负责从客户端处理业务

	go func() {
		//2.从客户端读取数据，拆包处理
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
			}

			go func(conn net.Conn) {
				//处理客户端的请求

				//--------》拆包——————————》

				//定义一个拆包的对象dp
				dp := NewDataPack()
				for {
					//1.第一次从conn读，把head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head err", err)
						break
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err", err)
						return
					}

					if msgHead.GetMsgLen() > 0 {
						//Msg有数据，需要二次读取

						//2.第二次从conn读，根据head的datalen再读取data内容
						msg := msgHead.(*Message)                //类型断言
						msg.Data = make([]byte, msg.GetMsgLen()) //给切片分配空间

						//根据dataLen的长度再次从io流读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err", err)
							return
						}

						//完整的消息读取完毕
						fmt.Println("------->Recv MsgId:", msg.Id, ",dataLen:", msg.DataLen, ",data:", string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	/*
		模拟客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err", err)
	}
	//创建一个封包对象

	dp := NewDataPack()

	//模拟粘包过程，封装两个包一同发送
	//封装第一个message包
	msg1 := &Message{
		Data:    []byte{'Z', 'i', 'n', 'x'},
		Id:      1,
		DataLen: 4,
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error:", err)
		return
	}
	//封装第二个message包
	msg2 := &Message{
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
		Id:      2,
		DataLen: 5,
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 error:", err)
		return
	}
	//将两个包黏在一起
	sendData1 = append(sendData1, sendData2...)
	//一次性发送给服务端
	conn.Write(sendData1)

	//客户端阻塞
	select {}
}
