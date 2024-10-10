package main

import (
	"fmt"
	"os"
	"runtime/trace"
)

// run: go run trace.go
// run web service: go tool trace trace.out
// start chrome: cd C:/"Program Files"/Google/Chrome/Application
// .\chrome.exe  -new-window --enable-blink-features=ShadowDOMV0,CustomElementsV0,HTMLImports
func main() {

	//创建trace文件
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	//启动trace goroutine
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	//main
	fmt.Println("Hello World")
}
