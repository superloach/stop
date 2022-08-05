package stop

import (
	"reflect"
	"runtime"
)

var ctxs = trimap[uintptr, uintptr, *ctxData]{}

type ctxData struct {
	Context chan struct{}
	Refs    uint64
	Handle  reflect.Value
}

var never = make(chan struct{})

func Go[T any](fn func() T) <-chan T {
	handle := make(chan T)

	go func() {
		// get caller of fn, if Go was called from fn
		pc, _, _, ok := runtime.Caller(4)
		if !ok {
			// get this func
			pc, _, _, _ = runtime.Caller(0)
		}

		entry := runtime.FuncForPC(pc).Entry()

		_, ctx, ok := ctxs.GetT(entry)
		if !ok {
			ctx = &ctxData{
				Context: make(chan struct{}),
				Refs:    0,
				Handle:  reflect.ValueOf(handle),
			}

			trimap[uintptr, uintptr, *ctxData](ctxs).Set(
				entry,
				ctx.Handle.Pointer(),
				ctx,
			)
		}

		ctx.Refs++

		val := fn()
		handle <- val

		ctx.Refs--

		if ctx.Refs == 0 {
			close(ctx.Context)
			ctxs.DelV(ctx)
		}
	}()

	return handle
}

func Context() <-chan struct{} {
	_, ctx := findEntry(3)
	if ctx == nil {
		return never
	}

	return ctx.Context
}

func Yield[T any]() chan<- T {
	_, ctx := findEntry(3)

	yield := make(chan T, 1)

	if ctx != nil {
		ctx.Refs++
	}

	go func() {
		val := <-yield

		if ctx != nil {
			ctx.Handle.Interface().(chan T) <- val
			ctx.Refs--
		}

		close(yield)
	}()

	return yield
}

func findEntry(depth int) (uintptr, *ctxData) {
	// get caller of fn (from Go)
	pc, _, _, ok := runtime.Caller(depth)
	if !ok {
		return 0, nil
	}

	entry := runtime.FuncForPC(pc).Entry()

	_, ctx, ok := ctxs.GetT(entry)
	if !ok {
		return findEntry(depth + 2)
	}

	return entry, ctx
}

func Stop[T any](h <-chan T) {
	id := reflect.ValueOf(h).Pointer()

	_, ctx, ok := ctxs.GetU(id)
	if ok {
		ctx.Context <- struct{}{}
	}
}
