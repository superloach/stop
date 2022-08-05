package main

import (
	"fmt"
	"time"
)

import "github.com/superloach/stop"

func main() {
	handle := stop.Go(func() int {
		return <-stop.Go(func() int {
			fmt.Println("waiting on context")

			<-stop.Context()
			fmt.Println("context closed")

			stop.Yield[int]() <- 123
			fmt.Println("yielded")

			return 456
		})
	})

	go func() {
		fmt.Println("sleeping 3s")
		time.Sleep(time.Second * 3)

		fmt.Println("stopping handle")
		stop.Stop(handle)
	}()

	fmt.Println("waiting for value from handle")
	fmt.Println("value from handle:", <-handle)

	fmt.Println("waiting for value from handle")
	fmt.Println("value from handle:", <-handle)
}
