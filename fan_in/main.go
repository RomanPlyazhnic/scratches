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

func fanOut(ctx context.Context, in chan int, numChan int) []chan int {
	out := make([]chan int, numChan)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := range numChan {
			out[i] = pipeline(ctx, in)
		}
	}()
	wg.Wait()

	return out
}

func pipeline(ctx context.Context, in chan int) (out chan int) {
	out = make(chan int)

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}

				out <- v
			}
		}
	}()

	return out
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	in := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			in <- i
		}
		close(in)
	}()

	chs := fanOut(ctx, in, 3)
	out := fanIn(ctx, chs)

	for i := range out {
		fmt.Println(i)
	}
}
