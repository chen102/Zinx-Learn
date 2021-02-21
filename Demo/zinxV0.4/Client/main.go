package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("client start...")
	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "192.168.110.221:8001")
	if err != nil {
		fmt.Println("client start err")
		return
	}
	for {
		_, err := conn.Write([]byte("Hello Zinx V0.2..."))
		if err != nil {
			fmt.Println("write conn err:", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf err:", err)
			return
		}
		fmt.Printf("server call back:%s,cnt=%d\n", buf, cnt)

		//cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
