package channel

import (
	"fmt"
	"strconv"
	"testing"
)

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			fmt.Printf("send x: %d \n", x)
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}
func Test_2(t *testing.T) {
	c := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println("receive x: " + strconv.Itoa(<-c))
		}
		quit <- 0
	}()
	fibonacci(c, quit)
}
