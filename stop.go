package stop

import "runtime"
import "reflect"
import "unsafe"

var _handles = map[unsafe.Pointer]uintptr{}
var _contexts = map[uintptr]chan struct{}{}

var _never = make(chan struct{})

func GoNothing(fn func()) chan struct{} {
	return Go(func() struct{} {
		fn()
		return struct{}{}
	})
}

func Go[T any](fn func() T) chan T {
	handle := make(chan T)

	go func() {
		pc, _, _, ok := runtime.Caller(0)
		if !ok {
			panic("caller not ok")
		}

		entry := runtime.FuncForPC(pc).Entry()
		id := reflect.ValueOf(handle).UnsafePointer()

		_handles[id] = entry
		_contexts[entry] = make(chan struct{})

		handle <- fn()
	}()

	return handle
}

func Context() <-chan struct{} {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		panic("caller not ok")
	}

	entry := runtime.FuncForPC(pc).Entry()

	ctx, ok := _contexts[entry]
	if !ok {
		println("never")
		return _never
	}

	return ctx
}

func Stop[T any](h chan T) {
	id := reflect.ValueOf(h).UnsafePointer()

	entry, ok := _handles[id]
	if !ok {
		panic("entry not ok")
	}

	ctx, ok := _contexts[entry]
	if !ok {
		panic("ctx not ok")
	}

	close(ctx)
}
