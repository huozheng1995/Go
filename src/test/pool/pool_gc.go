package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var pool = sync.Pool{
	New: func() interface{} {
		return "123"
	},
}

func main() {
	t := pool.Get().(string)
	fmt.Println(t)

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
