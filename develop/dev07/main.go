package main

import (
	"fmt"
	"reflect"
	"time"
)

func or1(channels ...<-chan interface{}) <-chan interface{} {
	branches := make([]reflect.SelectCase, 0, len(channels))
	for _, ch := range channels {
		vc := reflect.ValueOf(ch)
		branches = append(branches, reflect.SelectCase{
			Dir: reflect.SelectRecv, Chan: vc,
		})
	}

	ret := make(chan interface{})
	go func() {
		defer close(ret)
		reflect.Select(branches)
	}()

	return ret
}

func or2(channels ...<-chan interface{}) <-chan interface{} {
	ret := make(chan interface{})
	for _, ch := range channels {
		go func(ch <-chan interface{}) {
			<-ch

			select {
			case <-ret:
			default:
				close(ret)
			}
		}(ch)
	}

	return ret
}

func main() {
	// Реализация sig имеет недостаток в виде утечки горутин :(
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or1(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("1 done after %v\n", time.Since(start))

	start = time.Now()
	<-or2(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("2 done after %v\n", time.Since(start))
}
