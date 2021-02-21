package ziface

//IRequest接口：把客户端请求的链接信息，和请求的数据包装到一个Request中
type IRequest interface {
	GetConnection() IConnection
	GetData() []byte
	GetMsgID() uint32
}
