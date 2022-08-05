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
			return 123
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
}
