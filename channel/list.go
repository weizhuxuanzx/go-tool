package channel

// 桥接模式 管道队列
type Message interface {
	Reader(channel chan []byte) []byte
	Sender(channel chan []byte)
}
type MQServer struct {
	Name    string
	channel chan []byte
}

func (m MQServer) Reader(ms Message) []byte {
	return ms.Reader(m.channel)
}
func (m MQServer) Sender(ms Message) {
	ms.Sender(m.channel)
}
