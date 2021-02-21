package znet

import (
	"fmt"
	"zinx/utils"
	"zinx/ziface"
)

//消息处理模块的实现
type MsgHandle struct {
	//存放每个MsgID对应的处理方法
	Apis map[uint32]ziface.IRoute
	//负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作worker池的worker数量
	WorkerPoolSize uint32
}

//初始化
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRoute),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID=", request.GetMsgID(), " is NOT FOUND! Need Register!")
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRoute) {
	//判断id知否注册
	if _, ok := mh.Apis[msgID]; ok {
		fmt.Println("repeat api msgID=", msgID)
		panic("repeat api msgID")
	}
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID =", msgID, " succ!")
}

//启动一个Worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ { //创建若干Goroutine堵塞等待request
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go mh.StartWorker(i, mh.TaskQueue[i])
	}
}

//启动一个Worker工作流程
func (mh *MsgHandle) StartWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("WorkID=", workerID, " is started...")
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//将消息交个TaskQueue
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//消息平均分配给不同的worker
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize //取与
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(), " request MsgID=", request.GetMsgID(), " to WorkerID=", workerID)
	mh.TaskQueue[workerID] <- request
}
