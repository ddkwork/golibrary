package events

type (
	//Interface interface {
	//	SetEvent(event any)   //接受任意类型的消息事件
	//	SetEventsCap(cap int) //设置消息事件数量
	//	Events() <-chan any   //信道先进先出队列存储事件信号
	//	HandleEvent()         //处理消息事件
	//}
	Object struct{ event chan any }
)

// func NewObject() *Object               { return &Object{event: make(chan any, 1000)} }
// func New() Interface                   { return NewObject() }
func (o *Object) SetEventsCap(cap int) { o.event = make(chan any, cap) }
func (o *Object) SetEvent(event any)   { o.event <- event }
func (o *Object) Events() <-chan any   { return o.event }
func (o *Object) HandleEvent()         { panic("implement me") }
