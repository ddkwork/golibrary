package events

type (
	Object struct{ event chan any }
)

func (o *Object) SetEventsCap(cap int) { o.event = make(chan any, cap) }
func (o *Object) SetEvent(event any)   { o.event <- event }
func (o *Object) Events() <-chan any   { return o.event }
func (o *Object) HandleEvent()         { panic("implement me") }
