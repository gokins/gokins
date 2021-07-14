package engine

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	tm := time.NewTimer(time.Second * 5)
	select {
	case <-tm.C:
		fmt.Println(<-tm.C)
	default:
		fmt.Println("default")
	}
	fmt.Println(tm)
}

func TestSync(t *testing.T) {
	lk := sync.RWMutex{}
	go func() {
		for {
			lk.Lock()
			fmt.Println(1)
			lk.Unlock()
		}
	}()
	go func() {
		for {
			lk.Lock()
			fmt.Println(2)
			lk.Unlock()
		}
	}()
}
