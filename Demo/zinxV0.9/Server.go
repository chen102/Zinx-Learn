package main

import "zinx/znet"
import "zinx/ziface"
import "fmt"

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}
type HelloZinxRouter struct {
	znet.BaseRouter
}

//func (this *PingRouter) PreHandle(request ziface.IRequest) {
//fmt.Println("Call Router PreHandle...")
//_, err := request.GetConnection().GetTCPConnection().Write([]byte("Before ping..."))
//if err != nil {
//fmt.Println("call back before ping error")
//}
//}
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	//_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping... ping..."))
	//if err != nil {
	//fmt.Println("call back ping ping error")
	//}
	fmt.Println("recv from client:msgID:", request.GetMsgID(), " data:", string(request.GetData()))
	err := request.GetConnection().SendMsg(100, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}

}
func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle...")
	//_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping... ping..."))
	//if err != nil {
	//fmt.Println("call back ping ping error")
	//}
	fmt.Println("recv from client:msgID:", request.GetMsgID(), " data:", string(request.GetData()))
	err := request.GetConnection().SendMsg(101, []byte("hello...hello...hello"))
	if err != nil {
		fmt.Println(err)
	}

}

//func (this *PingRouter) PostHandle(request ziface.IRequest) {
//fmt.Println("Call Router PostHandle...")
//_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping..."))
//if err != nil {
//fmt.Println("call back after ping error")
//}

//}
//创建链接之后的钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("DoConnectionBegin is Called...")
	if err := conn.SendMsg(222, []byte("DoConnection BEGIN")); err != nil {
		fmt.Println(err)
	}
}

//结束链接之前的钩子函数
func DoConnectionEnd(conn ziface.IConnection) {
	fmt.Println("DoConnectionEnd is Called...")
	fmt.Println(conn.GetConnID(), "is End ")
}
func main() {
	s := znet.NewServer("[zinx v0.9]")
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionEnd)

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	s.Run()
}
