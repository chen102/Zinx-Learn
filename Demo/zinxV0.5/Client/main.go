package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
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
		dp := znet.NewDataPack()
		biniartMsg, _ := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx client Test Message")))
		_, err := conn.Write(biniartMsg)
		if err != nil {
			fmt.Println("clent write err:", err)
			return
		}
		binaryHead := make([]byte, dp.GetHandLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head err: ", err)
		}
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error: ", err)
			break
		}
		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read Messge data error: ", err)
			}
			fmt.Println("Recv Server Msg: ID=", msg.Id, " len=", msg.DataLen, " data=", string(msg.Data))
		}
		//cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
