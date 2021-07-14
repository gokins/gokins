package engine

import (
	"fmt"
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
