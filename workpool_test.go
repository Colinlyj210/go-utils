package go_utils

import (
	"fmt"
	"testing"
	"time"
)

func TestNewWorkPool(t *testing.T) {
	pool := NewWorkPool(10, 1000)
	for i := 0; i < 10000; i++ {
		pool.Execute(func(w *WorkPool, args ...interface{}) {
			fmt.Println("i=", args)
		}, i)
	}

	time.Sleep(time.Second * 5)

}
