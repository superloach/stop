package stop

import (
	"reflect"
	"runtime"
)

var _ = func() struct{} {
	println("warning: this program is using the library github.com/superloach/stop.")
	println("stop is not stable, and is only intended for educational use.")
	return struct{}{}
}

type ctxKey struct {
	Entry uintptr
	Handle uintptr
	Data *ctxData
}

type ctxMap map[ctxKey]*ctxData

var ctxs = ctxMap{}

type ctxData struct {
	Entry uintptr
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

		ctx, ok := ctxs[ctxKey{Entry: entry}]
		if !ok {
			ctx = &ctxData{
				Context: make(chan struct{}),
				Refs:    0,
				Handle:  reflect.ValueOf(handle),
			}
		}

			keys := []ctxKey{{
				Entry: entry,
			}, {
				Handle: ctx.Handle.Pointer(),
			}, {
				Data: ctx,
			}}

			for _, k := range keys {
				ctxs[k] = ctx
			}

		ctx.Refs++

		val := fn()
		handle <- val

		ctx.Refs--

		if ctx.Refs == 0 {
			for _, k := range keys {
				delete(ctxs, k)
			}
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

	ctx, ok := ctxs[ctxKey{Entry: entry}]
	if !ok {
		return findEntry(depth + 2)
	}

	return entry, ctx
}

func Stop[T any](h <-chan T) {
	ctx, ok := ctxs[ctxKey{
		Handle: reflect.ValueOf(h).Pointer(),
	}]
	if ok {
		ctx.Context <- struct{}{}
	}
}
