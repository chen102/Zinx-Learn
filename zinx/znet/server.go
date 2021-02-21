package znet

import (
	//	"errors"
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

//IServer的接口实现
type Server struct {
	Name      string
	IPVersion string //服务器绑定ip版本
	IP        string
	Port      int

	//当前Server添加一个router，server注册的链接对应的处理业务
	//Router ziface.IRoute
	//当前Server的消息管理模块
	MsgHandler ziface.IMsgHandle
	//连接管理器
	ConnMgr ziface.IConnManager
	//Server创建链接后自动调用Hook函数
	OnConnStart func(conn ziface.IConnection)
	//Server销毁链接之前自动调用Hook函数
	OnConnStop func(conn ziface.IConnection)
}

//定义当前客户端链接的所绑定handle api
//属于业务代码--前期测试用的
//func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
////回显业务
//fmt.Println("[Conn Handle] CallBackToClient ...")
//if _, err := conn.Write(data[:cnt]); err != nil {
//fmt.Println("write back buf err:", err)
//return errors.New("CallBack error")
//}
//return nil
//}
func (s Server) Start() {
	fmt.Printf("[Zinx] Server Name:%s,listenner at IP : %s Port:%d is starting...\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	go func() {

		//开启消息队列，及Worker工作池
		s.MsgHandler.StartWorkerPool()
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addt error:", err)
			return
		}
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, "err:", err)
			return
		}
		fmt.Println("start Zinx server succ", s.Name, "succ,Listenning")
		var cid uint32
		cid = 0 //connID
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err:", err)
				continue
			}
			//设置最大连接个数的判断
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				fmt.Println("MaxConn is:", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}
			//将处理新连接的业务方法和conn进行绑定得到我们的链接模块
			dealConn := NewConnection(&s, conn, cid, s.MsgHandler)
			cid++

			go dealConn.Start()
		}

	}()
}
func (s Server) Stop() {
	fmt.Println("[STOP]Zinx server name", s.Name)
	s.ConnMgr.ClearConn()
}
func (s Server) Run() {
	s.Start()
	//阻塞状态
	select {}
}
func (s *Server) AddRouter(msgID uint32, router ziface.IRoute) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router succ!!")
}

//初始化Server
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		//Router:    nil,
		MsgHandler: NewMsgHandle(), //初始化消息模块
		ConnMgr:    NewConnManager(),
	}
	return s
}
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc

}
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("Call OnConnStart()...")
		s.OnConnStop(conn)

	}
}
