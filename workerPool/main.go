package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const workerCount = 4

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

	fn := func(i int) int {
		time.Sleep(time.Second)
		return i
	}
	out := workerPool(ctx, in, workerCount, fn)

	for i := range out {
		fmt.Println(i)
	}
}

func workerPool(ctx context.Context, in chan int, workerCount int, fn func(int) int) (out chan int) {
	out = make(chan int)

	wg := &sync.WaitGroup{}
	go func() {
		for range workerCount {
			wg.Add(1)
			go worker(ctx, wg, in, out, fn)
		}
		wg.Wait()
		close(out)
	}()

	return out
}

func worker(ctx context.Context, wg *sync.WaitGroup, in chan int, out chan int, fn func(int) int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case v, ok := <-in:
			if !ok {
				return
			}

			select {
			case <-ctx.Done():
				return
			case out <- fn(v):
			}
		}
	}
}
