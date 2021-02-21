package ziface

//服务器接口
type IServer interface {
	Start() //启动服务器
	Stop()  //停止服务器
	Run()   //运行服务器
	//路由功能:给当前的服务注册一个路由方法,供客户端的链接处理使用
	//AddRouter(router IRoute)
	AddRouter(msgID uint32, router IRoute)
	GetConnMgr() IConnManager //获取当前server的连接管理器
	SetOnConnStart(func(connection IConnection))
	SetOnConnStop(func(connection IConnection))
	CallOnConnStart(connection IConnection)
	CallOnConnStop(connection IConnection)
}
