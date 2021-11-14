package interpreter

type LazyObject struct {
	val Value
}

func (*LazyObject) Type() String {
	return "internal+lazy"
}

var lazyNativeMap = NativeMap{}

func (*LazyObject) Natives() NativeMap {
	return lazyNativeMap
}

func (L *LazyObject) Resolve() Value {
	return L.val
}

func (L *LazyObject) WakeUp(val Value) {
	L.val = val
}

func Lazy() *LazyObject {
	return &LazyObject{}
}