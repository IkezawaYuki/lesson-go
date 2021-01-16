package main

import (
	"fmt"
	"time"
)

const (
	CHAN_CAPACITY  = 0
	MAX_GOROUTINES = 60000
	MEASURE_COUNT  = 100
)

const durationInSeconds = 100

func main() {
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	end := make(chan struct{})
	result := make(chan int)

	go pingPong(ch1, ch2, end, result)
	go pingPong(ch2, ch1, end, result)

	start := time.Now()
	ch1 <- struct{}{}
	<-time.Tick(time.Second * durationInSeconds)
	close(end)
	r1 := <-result
	r2 := <-result
	elapsed := time.Now().Sub(start)

	fmt.Printf("Elapsed Time = %v\n", elapsed)
	fmt.Printf("%d per second\n", r1/durationInSeconds)
	fmt.Printf("%d per second\n", r2/durationInSeconds)
}

func pingPong(in <-chan struct{}, out chan<- struct{},
	end <-chan struct{}, result chan<- int) {
	for i := 0; ; i++ {
		select {
		case v := <-in:
			select {
			case out <- v:
			case <-end:
				result <- i
				return
			}
		case <-end:
			result <- i
			return
		}
	}
}

func mai2() {
	next := make(chan int)
	final := make(chan int)
	go pipe(next, 0, final)
	next <- 0
	i := <-final
	if i != 0 {
		panic("i != 0")
	}
	fmt.Printf("\n%d goroutine are created\n", MAX_GOROUTINES)
	oneByOneSending(next, final)
	continuousSending(next, final)
}

func oneByOneSending(next chan<- int, final <-chan int) {
	var total int64

	for v := 1; v <= MEASURE_COUNT; v++ {
		start := time.Now()
		next <- v
		<-final
		end := time.Now()

		diff := end.Sub(start)
		fmt.Printf("%3d: %v\n", v, diff)
		total += diff.Nanoseconds()
	}
	fmt.Printf("average round trip time = %d nano seconds\n", total/MEASURE_COUNT)
	fmt.Printf("average switching trip time = %d nano seconds\n", total/MEASURE_COUNT*MAX_GOROUTINES)
}

func continuousSending(next chan<- int, final <-chan int) {
	start := time.Now()
	go func() {
		for i := 0; i < MEASURE_COUNT; i++ {
			next <- i
		}
	}()

	for i := 0; i < MEASURE_COUNT; i++ {
		<-final
	}
	end := time.Now()
	diff := end.Sub(start)
	fmt.Printf("elapse time for sending %d values ... %d nano seconds\n",
		MEASURE_COUNT, diff.Nanoseconds())
}

func pipe(prev <-chan int, stages int, final chan<- int) {
	next := make(chan int, CHAN_CAPACITY)
	stages++
	if stages%10000 == 0 {
		time.Sleep(time.Second)
		fmt.Printf("%d\n", stages)
	}

	if stages == MAX_GOROUTINES {
		for v := range prev {
			final <- v
		}
	} else {
		go pipe(next, stages, final)
		for v := range prev {
			next <- v
		}
	}
}
