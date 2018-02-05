package tachikoma

type Emitter interface {
	Start(c chan interface{})
}

type EmitterFnProto func(c chan interface{})

type EmitterFn struct {
	fn EmitterFnProto
}

func (e EmitterFn) Start(c chan interface{}) {
	e.fn(c)
}

func NewEmitterFn(fn EmitterFnProto) EmitterFn {
	return EmitterFn{fn: fn}
}
