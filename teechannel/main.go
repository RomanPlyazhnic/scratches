package main

import (
	"fmt"
	"sync"
)

func main() {
	in := make(chan int)

	go func() {
		defer close(in)
		for i := range 10 {
			in <- i
		}
	}()

	out1, out2 := tee(in)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()

		for v := range out1 {
			fmt.Println("chan1", v)
		}
	}()
	go func() {
		defer wg.Done()

		for v := range out2 {
			fmt.Println("chan2", v)
		}
	}()

	wg.Wait()
}

func tee[T any](in chan T) (out1, out2 chan T) {
	out1 = make(chan T)
	out2 = make(chan T)

	go func() {
		defer close(out1)
		defer close(out2)

		for v := range in {
			var out1, out2 = out1, out2
			for range 2 {
				select {
				case out1 <- v:
					out1 = nil
				case out2 <- v:
					out2 = nil
				}
			}
		}
	}()

	return out1, out2
}
