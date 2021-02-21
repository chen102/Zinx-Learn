package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	listenner, err := net.Listen("tcp", "192.168.110.221:8003")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}
	//模拟的服务器
	go func() {
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
			}
			go func(conn net.Conn) {
				//处理客户端的请求 拆包的过程
				dp := NewDataPack()
				for {
					headData := make([]byte, dp.GetHandLen())
					_, err := io.ReadFull(conn, headData) //读头消息(8字节)
					if err != nil {
						fmt.Println("read head error:", err)
						break
					}
					msgHead, err := dp.Unpack(headData) //返回数据的长度和类型
					if err != nil {
						fmt.Println("server unpack err:", err)
					}
					//判断是否有数据可读
					if msgHead.GetMsgLen() > 0 {
						msg := msgHead.(*Message) //类型断言(msgHead是imessage)
						msg.Data = make([]byte, msg.GetMsgLen())
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err:", err)
							return
						}
						fmt.Println("Recv MsgID: ", msg.Id, " datalen: ", msg.DataLen, " data: ", string(msg.Data))

					}

				}
			}(conn)
		}

	}()
	//模拟客户端
	conn, err := net.Dial("tcp", "192.168.110.221:8003")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}
	dp := NewDataPack()
	//模拟粘包问题
	//分别封装两个包，一同发向服务端
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error", err)
		return
	}
	msg2 := &Message{
		Id:      2,
		DataLen: 5,
		Data:    []byte{'o', 'h', 'h', 'h', 'h'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg1 error", err)
		return
	}
	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)

	select {}

}
