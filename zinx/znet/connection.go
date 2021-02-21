package znet

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

type Connetcion struct {
	//当前Conn属于那个Server
	TcpServer ziface.IServer
	Conn      *net.TCPConn
	ConnID    uint32
	isClosed  bool //Reader告知Writer退出
	//handleAPI ziface.HandleFunc
	//ExitChan chan bool  //使用context统一管理
	ctx    context.Context
	cancel context.CancelFunc //退出函数
	//用于读写Goroutine的通信
	msgChan     chan []byte
	msgBuffChan chan []byte
	sync.RWMutex
	//Router   ziface.IRoute //该链接处理的方法Router
	//消息管理msgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle

	//链接属性集合
	property map[string]interface{}
	//保护链接属性的锁
	propertyLock sync.RWMutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connetcion {
	c := &Connetcion{
		TcpServer: server,
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		//handleAPI: callback_api,
		MsgHandler:  msgHandler,
		msgChan:     make(chan []byte),
		msgBuffChan: make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		//ExitChan:    make(chan bool, 1),
		property: make(map[string]interface{}),
	}
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

//将conn加入ConnManager中
//链接的读业务方法
func (c *Connetcion) StartReader() {
	fmt.Println("[Reader Goroutine is runing...]")
	defer fmt.Println("connID=", c.ConnID, "Reader is exit,remote addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {

		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//fmt.Println("recv buf err:", err)
		//continue
		//}
		//if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
		//fmt.Println("ConnID ", c.ConnID, "handle is error ", err)
		//break
		//}
		select {
		case <-c.ctx.Done():
			return
		default:
			//拆包过程
			dp := NewDataPack()
			headData := make([]byte, dp.GetHandLen())
			_, err := io.ReadFull(c.GetTCPConnection(), headData)
			if err != nil {
				fmt.Println("read msg heada error:", err)
				break
			}
			msg, err := dp.Unpack(headData) //返回ziface.IMessage
			if err != nil {
				fmt.Println("unpack error", err)
				break
			}
			var data []byte
			if msg.GetMsgLen() > 0 {
				data = make([]byte, msg.GetMsgLen())
				if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
					fmt.Println("read msg data error:", err)
					break
				}
			}
			msg.SetMsgData(data)

			//得到当前conn数据的Request请求数据
			req := Request{
				conn: c,
				msg:  msg,
			}
			//判断是否已经开启工作池
			if utils.GlobalObject.WorkerPoolSize > 0 {
				c.MsgHandler.SendMsgToTaskQueue(&req) //将消息发送给Worker工作池处理
			} else {
				//根据绑定好的MsgID找到对应处理api业务 执行
				go c.MsgHandler.DoMsgHandler(&req)
			}

		}

		//从路由中，找到注册绑定的Conn对应的router调用
		//go func(request ziface.IRequest) {
		//c.Router.PreHandle(request)
		//c.Router.Handle(request)
		//c.Router.PostHandle(request)

		//}(&req)
	}
}

//写消息Goroutine
func (c *Connetcion) StartWrite() {
	fmt.Println("[Writer Gortine is running...]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")
	//阻塞等待channel消息
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error: ", err)
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				_, err := c.Conn.Write(data)
				if err != nil {
					fmt.Println("Send buff data error: ", err)
					return
				}

			} else {
				fmt.Println("msgBuffChan is Closed")
				break
			}
		case <-c.ctx.Done(): //代表Reader已经退出，此时Writer也要退出
			return
		}
	}

}

//提供一个sendMSG方法(封包发送)
func (c Connetcion) SendMsg(msgId uint32, data []byte) error {
	c.RLock()
	if c.isClosed == true {
		c.RLocker()
		return errors.New("Conn closed when send msg ")
	}
	c.RLocker()
	//封包过程
	dp := NewDataPack()
	binartMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id:", msgId)
		return errors.New("Pack error msg")
	}

	//if _, err := c.Conn.Write(binartMsg); err != nil {
	//fmt.Println("Write msg id", msgId, " error:", err)
	//return errors.New("conn Write error")
	//}
	c.msgChan <- binartMsg //传给通信chnnel
	return nil
}
func (c Connetcion) SendBuffMsg(msgId uint32, data []byte) error {
	c.RLock()
	if c.isClosed == true {
		c.RLocker()
		return errors.New("Conn closed when send msg ")
	}
	c.RLocker()
	//封包过程
	dp := NewDataPack()
	binartMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id:", msgId)
		return errors.New("Pack error msg")
	}

	//if _, err := c.Conn.Write(binartMsg); err != nil {
	//fmt.Println("Write msg id", msgId, " error:", err)
	//return errors.New("conn Write error")
	//}
	c.msgBuffChan <- binartMsg //传给通信chnnel
	return nil
}

//启动链接
func (c *Connetcion) Start() {
	fmt.Println("Conn Start()... ConnID=", c.ConnID)
	c.ctx, c.cancel = context.WithCancel(context.Background())
	//启动从当前链接的读数据的业务
	go c.StartReader()
	//启动从当前链接的写数据的业务
	go c.StartWrite()
	//按照开发者传递进来的 创建链接之后需要调用的处理业务，执行对应的Hook函数
	c.TcpServer.CallOnConnStart(c)
}

//停止链接
func (c *Connetcion) Stop() {
	fmt.Println("Conn Stop.. ConnID=", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true
	//调用开发者注册的销毁链接之前需要执行的业务
	c.TcpServer.CallOnConnStop(c)
	c.Conn.Close()
	//c.ExitChan <- true //告知Write退出
	c.cancel() //统一关闭读写
	//将当前从connMgr中删除
	c.TcpServer.GetConnMgr().Remove(c)
	//close(c.ExitChan)

	close(c.msgChan)
	close(c.msgBuffChan)
}

//获取当前链接的绑定socket conn
func (c *Connetcion) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前链接模块的链接ID
func (c *Connetcion) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的TCP状态 IP port
func (c *Connetcion) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
func (c *Connetcion) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}
func (c *Connetcion) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}

}
func (c *Connetcion) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)

}
