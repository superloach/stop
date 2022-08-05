package main

import (
	"fmt"
	"time"
)

import "github.com/superloach/stop"

func main() {
	h := stop.Go(func() int {
		i := <-stop.Go(func() int {
			fmt.Println("waiting on context")
			<-stop.Context()

			fmt.Println("context closed")
			stop.Yield[int]() <- 123
			fmt.Println("yielded")
			return 456
		})
		return i
	})

	go func() {
		fmt.Println("sleeping 3s")
		time.Sleep(time.Second * 3)
		fmt.Println("stopping context")
		stop.Stop(h)

		fmt.Println("context stopped")
	}()

	fmt.Println("waiting on handle")
	fmt.Println(<-h)
	fmt.Println("waiting on handle")
	fmt.Println(<-h)
}
