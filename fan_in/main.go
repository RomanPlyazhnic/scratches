package main

import (
	"context"
	"fmt"
	"sync"
)

func fanIn(ctx context.Context, chs []chan int) chan int {
	out := make(chan int)

	wg := &sync.WaitGroup{}

	go func() {
		for _, ch := range chs {
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

						select {
						case out <- i:
						case <-ctx.Done():
							return
						}
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
	var chs []chan int

	for range 3 {
		chs = append(chs, make(chan int))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for i := range 3 {
			go func() {
				for j := 0; j < 5; j++ {
					chs[i] <- j
				}
				close(chs[i])
			}()
		}
	}()

	out := fanIn(ctx, chs)

	for i := range out {
		fmt.Println(i)
	}
}
