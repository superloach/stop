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
		fmt.Println("sleeping 3s")
		time.Sleep(time.Second * 3)

		fmt.Println("stopping handle")
		stop.Stop(handle)
	}()

	for val := range handle {
		fmt.Println("value from handle:", val)
	}
}
