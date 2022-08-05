package stop

import (
	"reflect"
	"runtime"
)

func init() {
	println("warning: this program is using the library github.com/superloach/stop, which is not stable, and is only intended for educational use.")
}

type ctxData struct {
	Context chan struct{}
	Refs    uint64

	Entry  uintptr
	Handle reflect.Value
}

type ctxMap map[uintptr]*ctxData

func (m ctxMap) Set(ctx *ctxData) {
	m[ctx.Entry] = ctx
	m[ctx.Handle.Pointer()] = ctx
}

func (m ctxMap) Del(ctx *ctxData) {
	delete(m, ctx.Entry)
	delete(m, ctx.Handle.Pointer())
}

var ctxs = ctxMap{}

func Go[T any](fn func() T) <-chan T {
	handle := make(chan T)

	go func() {
		// get this func
		pc, _, _, _ := runtime.Caller(0)

		entry := runtime.FuncForPC(pc).Entry()

		ctx, ok := ctxs[entry]
		if !ok {
			ctx = &ctxData{
				Entry:   entry,
				Context: make(chan struct{}),
				Refs:    0,
				Handle:  reflect.ValueOf(handle),
			}
		}

		ctxs.Set(ctx)

		ctx.Refs++

		handle <- fn()
		close(handle)

		ctx.Refs--

		if ctx.Refs == 0 {
			ctxs.Del(ctx)
		}
	}()

	return handle
}

func GoNothing(fn func()) <-chan struct{} {
	return Go[struct{}](func() struct{} {
		fn()
		return struct{}{}
	})
}

// Context fetches the Context channel for the current goroutine. This channel closes if Stop is called on its Handle.
func Context() <-chan struct{} {
	ctx := findEntry()
	if ctx == nil {
		panic("cannot fetch a context from a non-managed goroutine")
	}

	return ctx.Context
}

// Yield sends a value onto the Handle for the current goroutine.
func Yield[T any](val T) {
	ctx := findEntry()
	if ctx == nil {
		panic("cannot yield from a non-managed goroutine")
	}

	ctx.Refs++
	ctx.Handle.Interface().(chan T) <- val
	ctx.Refs--
}

// Pass works like Yield, except it consumes an entire channel rather than one value.
func Pass[T any](c <-chan T) {
	ctx := findEntry()
	if ctx == nil {
		panic("cannot pass from a non-managed goroutine")
	}

	ctx.Refs++
	for val := range c {
		ctx.Handle.Interface().(chan T) <- val
	}
	ctx.Refs--
}

func findEntry() *ctxData {
	// get caller of fn (from Go)
	pc, ok := uintptr(0), false

	for depth := 0; ; depth++ {
		pc, _, _, ok = runtime.Caller(depth)
		if !ok {
			// went past depth
			return nil
		}

		entry := runtime.FuncForPC(pc).Entry()

		ctx, ok := ctxs[entry]
		if ok {
			return ctx
		}
	}
}

func Stop[T any](h <-chan T) {
	ctx, ok := ctxs[reflect.ValueOf(h).Pointer()]
	if !ok {
		panic("no context found")
	}

	close(ctx.Context)
}
