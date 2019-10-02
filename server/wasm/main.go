package main

import (
	"flag"
	"github.com/DharmaOfCode/gorp/server/wasm/api"
	"github.com/go-redis/redis"
	"strconv"
	"sync"
	"syscall/js"
	//"time"
)

var (
	listen = flag.String("listen", ":8080", "listen address")
	dir    = flag.String("dir", ".", "directory to serve")
)

func add(this js.Value, i []js.Value) interface{} {
	println("running")
	value1 := js.Global().Get("document").Call("getElementById", i[0].String()).Get("value").String()
	value2 := js.Global().Get("document").Call("getElementById", i[1].String()).Get("value").String()

	int1, _ := strconv.Atoi(value1)
	int2, _ := strconv.Atoi(value2)

	js.Global().Get("document").Call("getElementById", i[2].String()).Set("value", int1+int2)
	println("woot")
	return ""
}

func registerCallbacks() {
	js.Global().Set("add", js.FuncOf(add))
}

func LoggerWorker(messages *[]string, wg *sync.WaitGroup){
	defer wg.Done()

	for m, _ := range *messages{
		println(&m)
	}

}

func Tester(s string){
	println(s)
}


func main() {
	c := make(chan struct{}, 0)

	println("WASM Go Inialized")

	helpers.Follow("")

	var wg sync.WaitGroup
	wg.Add(1)

	// register functions
	registerCallbacks()
	<-c
}


