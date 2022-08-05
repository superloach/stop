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
	Entry  uintptr
	Handle uintptr
	Data   *ctxData
}

type ctxMap map[ctxKey]*ctxData

var ctxs = ctxMap{}

type ctxData struct {
	Entry   uintptr
	Context chan struct{}
	Refs    uint64
	Handle  reflect.Value
}

var never = make(chan struct{})

func Go[T any](fn func() T) <-chan T {
	handle := make(chan T)

	go func() {
		// get this func
		pc, _, _, _ := runtime.Caller(0)

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

		handle <- fn()
		close(handle)

		ctx.Refs--

		if ctx.Refs == 0 {
			for _, k := range keys {
				delete(ctxs, k)
			}
		}
	}()

	return handle
}

// Context fetches the context channel for the current goroutine. This channel closes if
func Context() <-chan struct{} {
	_, ctx := findEntry()
	if ctx == nil {
		return never
	}

	return ctx.Context
}

// Yield sends a value onto the return handle for the current goroutine.
func Yield[T any](val T) {
	_, ctx := findEntry()

	if ctx != nil {
		ctx.Refs++
		ctx.Handle.Interface().(chan T) <- val
		ctx.Refs--
	}
}

// Pass works like Yield, except it consumes an entire channel rather than one value.
func Pass[T any](c <-chan T) {
	_, ctx := findEntry()

	if ctx != nil {
		for val := range c {
			ctx.Refs++
			ctx.Handle.Interface().(chan T) <- val
			ctx.Refs--
		}
	}
}

func findEntry() (uintptr, *ctxData) {
	// get caller of fn (from Go)
	entry := uintptr(0)
	ctx := (*ctxData)(nil)

	depth := 0

	for {
		pc, _, _, ok := runtime.Caller(depth)
		if !ok {
			return 0, nil
		}

		entry := runtime.FuncForPC(pc).Entry()

		ctx, ok = ctxs[ctxKey{Entry: entry}]
		if ok {
			break
		}

		depth++
	}

	return entry, ctx
}

func Stop[T any](h <-chan T) {
	ctx, ok := ctxs[ctxKey{
		Handle: reflect.ValueOf(h).Pointer(),
	}]
	if ok {
		close(ctx.Context)
	}
}
