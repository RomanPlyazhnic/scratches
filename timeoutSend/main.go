package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	in := make(chan int)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer func() {
			close(in)
			wg.Done()
		}()

		for i := 0; i < 30; i++ {
			time.Sleep(100 * time.Millisecond)

			select {
			case in <- i:
				fmt.Println("sent", i)
			case <-ctx.Done():
				return
			}
		}
	}()

	read(ctx, in)

	wg.Wait()
}

func read(ctx context.Context, ch <-chan int) {
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return
			}
			fmt.Println(v)
		case <-ctx.Done():
			fmt.Println("TIMEOUT")
			return
		}
	}
}
