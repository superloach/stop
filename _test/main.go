package main

import (
	"fmt"
	"time"
)

import "github.com/superloach/stop"

func main() {
	handle := stop.Go(func() string {
		stop.Yield("(outer) yield before inner")

		stop.Pass(stop.Go(func() string {
			stop.Yield("(inner) yield before stop")

			<-stop.Context()

			stop.Yield("(inner) yield after stop")

			return "(inner) return"
		}))

		stop.Yield("(outer) yield after inner")

		return "(outer) return"
	})

	go func() {
		fmt.Println("sleeping 1s on handle")
		time.Sleep(time.Second)

		fmt.Println("stopping handle")
		stop.Stop(handle)
	}()

	for val := range handle {
		fmt.Println("value from handle:", val)
	}

	handle2 := stop.GoNothing(func() {
		fmt.Println("waiting for stop")

		<-stop.Context()

		fmt.Println("stopped")
	})

	go func() {
		fmt.Println("sleeping 1s on handle2")
		time.Sleep(time.Second)

		fmt.Println("stopping handle2")
		stop.Stop(handle2)

		fmt.Println("stopped handle2")
	}()

	fmt.Println("waiting for handle2")
	<-handle2

	fmt.Println("handle2 closed")
}
