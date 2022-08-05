package stop

import (
	"reflect"
	"runtime"
	"unsafe"
)

var ctxs = trimap[uintptr, unsafe.Pointer, *ctxData]{}

type ctxData struct {
	Context chan struct{}
	Refs    uint64
}

var never = &ctxData{
	Context: make(chan struct{}),
	Refs:    0,
}

func GoNothing(fn func()) chan struct{} {
	return Go(func() struct{} {
		fn()
		return struct{}{}
	})
}

func Go[T any](fn func() T) chan T {
	handle := make(chan T)

	go func() {
		// get caller of fn, if Go was called from fn
		pc, _, _, ok := runtime.Caller(4)
		if !ok {
			// get this func
			pc, _, _, ok = runtime.Caller(0)
			if !ok {
				panic("caller not ok")
			}
		}

		entry := runtime.FuncForPC(pc).Entry()

		_, ctx, ok := ctxs.GetT(entry)
		if !ok {
			ctx = &ctxData{
				Context: make(chan struct{}),
				Refs:    0,
			}

			trimap[uintptr, unsafe.Pointer, *ctxData](ctxs).Set(
				entry,
				reflect.ValueOf(handle).UnsafePointer(),
				ctx,
			)
		}

		ctx.Refs++

		handle <- fn()

		ctx.Refs--
		if ctx.Refs == 0 {
			ctxs.DelV(ctx)
		}
	}()

	return handle
}

func Context() <-chan struct{} {
	return cancel(3)
}

func cancel(depth int) <-chan struct{} {
	// get caller of fn (from Go)
	pc, _, _, ok := runtime.Caller(depth)
	if !ok {
		println("caller not ok")
		return nil
	}

	entry := runtime.FuncForPC(pc).Entry()

	_, ctx, ok := ctxs.GetT(entry)
	if !ok {
		next := cancel(depth + 2)
		if next != nil {
			return next
		}

		println("never")
		return never.Context
	}

	return ctx.Context
}

func Stop[T any](h chan T) {
	id := reflect.ValueOf(h).UnsafePointer()

	_, ctx, ok := ctxs.GetU(id)
	if !ok {
		panic("ctx not ok")
	}

	close(ctx.Context)
	ctxs.DelU(id)
}
