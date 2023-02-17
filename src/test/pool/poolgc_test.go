package pool

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

var pool = sync.Pool{
	New: func() interface{} {
		return "123"
	},
}

func Test_2(t *testing.T) {
	val := pool.Get().(string)
	fmt.Println(val)

	pool.Put("234")
	pool.Put("345")
	pool.Put("456")
	pool.Put("567")

	runtime.GC()
	time.Sleep(1 * time.Second)

	t2 := pool.Get().(string)
	fmt.Println(t2)

	runtime.GC()
	time.Sleep(1 * time.Second)

	t2 = pool.Get().(string)
	fmt.Println(t2)
}
