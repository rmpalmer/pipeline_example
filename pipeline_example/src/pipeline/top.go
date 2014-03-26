package main 

import (
	"fmt"
	"sync"
)

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _,n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(in <- chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n*n
		}
		close(out)
	}()
	return out
}

func merge(done <-chan struct {}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	
	// the output channel 
	out := make(chan int)
	
	// the function for merging; will one for each input channel
	output := func(c <-chan int) {
	
		defer wg.Done()
		// read all input until done?
		for n := range c {
			select {
				case out <- n:
				case <-done:
					return
			}
		}
	
		// ..then signal to wait group
		wg.Done()
	}
	
	// wait to get a done message from each of the input channels
	wg.Add(len(cs))
	
	// launch a consumer goroutine for each input channel
	for _,c := range cs {
		go output(c)
	}
	
	// final goroutine that waits for completions
	go func() {
		wg.Wait()
		close(out)
	} ()
	
	// give the output channel back to the caller.
	return out
}

func main() {
	fmt.Printf("Begin\n")
	
	in := gen(12, 3)
	c1 := sq(in)
	c2 := sq(in)
	
	// a channel to tell the producers to stop
	done := make (chan struct{})
	defer close(done)
	
	// instead of consuming from all channels, consume only one
	out := merge(done, c1, c2)
	fmt.Println(<-out)
		
	fmt.Printf("End\n")
}

