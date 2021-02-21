package ziface

//消息管理抽象层
type IMsgHandle interface {
	DoMsgHandler(request IRequest)
	AddRouter(msgID uint32, router IRoute)
	StartWorkerPool() //启动worker工作池
	SendMsgToTaskQueue(request IRequest)
}
