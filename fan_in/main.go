package main

import (
	"context"
	"fmt"
	"sync"
)

func fanIn(ctx context.Context, chans []chan int) chan int {
	out := make(chan int)

	wg := &sync.WaitGroup{}

	go func() {
		for _, ch := range chans {
			wg.Add(1)

			go func() {
				defer wg.Done()

				for {
					select {
					case <-ctx.Done():
						return
					case i, ok := <-ch:
						if !ok {
							return
						}
						out <- i
					}
				}
			}()
		}

		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	chans := make([]chan int, 3)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for _, ch := range chans {
			ch = make(chan int)

			go func() {
				for i := 0; i < 5; i++ {
					ch <- i
				}
				close(ch)
			}()
		}
	}()

	out := fanIn(ctx, chans)

	for i := range out {
		fmt.Println(i)
	}
}
