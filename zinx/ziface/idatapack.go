//封包、拆包模块
//解决TCP粘包问题----针对Message进行TLV格式的封装/拆包
//先写消息的长度在写消息的ID，最后写消息的内容
//先读取固定长度的hand，消息内容的长度和消息的类型；在根据消息内容的长度，再次进行一次读写，从conn中读取消息的内容
package ziface

type IDataPack interface {
	GetHandLen() uint32
	Pack(msg IMessage) ([]byte, error)
	Unpack([]byte) (IMessage, error)
}
