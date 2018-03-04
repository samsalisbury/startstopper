package startstopper_test

import (
	"fmt"

	"github.com/samsalisbury/startstopper"
)

type SomeProcessor struct {
	Received []int
	startstopper.StartStopper
}

func (sp *SomeProcessor) Process(ch <-chan int) {
	sp.Start()
	select {
	case <-sp.Stopped():
		return
	default:
	}
	for {
		select {
		case <-sp.Stopped():
			return
		case item := <-ch:
			sp.Received = append(sp.Received, item)
		}
	}
}

func Example() {
	p := &SomeProcessor{}
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
			if i == 4 {
				p.Stop()
			}
		}
		p.Stop()
	}()
	p.Process(ch)
	fmt.Println(p.Received)
	p.Process(ch)
	fmt.Println(p.Received)
	// output:
	// [0 1 2 3 4]
	// [0 1 2 3 4 5 6 7 8 9]
}
