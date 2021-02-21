package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

type GlobalObj struct {
	//Server
	TcpServer ziface.IServer
	Host      string
	TcpPort   int
	Name      string
	//Zinx
	Version          string
	MaxConn          int    //当前服务器主机允许最大链接数
	MaxPackageSize   uint32 //当前Zinx数据包的最大值
	WorkerPoolSize   uint32 //当前业务工作Worker池的Goroutine数量
	MaxWorkerTaskLen uint32 //Zinx允许用户最多开辟多少个Worker
	MaxMsgChanLen    uint32
}

var GlobalObject *GlobalObj

//从zinx.json去加载用于自定义的参数
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//json文件解析到sturct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

//提供一个init方法，初始化当前的GlobalObject
func init() { //读取用户配置好的zinx.json 参数替换
	//默认值
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.4",
		TcpPort:          8001,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
	}
	GlobalObject.Reload()
}
