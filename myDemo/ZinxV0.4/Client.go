package main

import (
	"fmt"
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
		//2.连接调用write写数据
		_, err = conn.Write([]byte("Hello Zinx V0.1...."))
		if err != nil {
			fmt.Println("write conn err", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error:", err)
			return
		}

		fmt.Printf("server call back:%s,cnt = %d\n", buf, cnt)

		//cpu阻塞，防止无限循环
		time.Sleep(1 * time.Second)

	}

}
