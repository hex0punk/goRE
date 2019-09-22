package main

import (
	"flag"
	"github.com/go-redis/redis"
	"strconv"
	"sync"
	"syscall/js"
	"time"
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

func RedisTest(){
	pubsub := rdb.Subscribe("mychannel1")

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive()
	if err != nil {
		panic(err)
	}

	// Go channel which receives messages.
	ch := pubsub.Channel()

	// Publish a message.
	err = rdb.Publish("mychannel1", "hello").Err()
	if err != nil {
		panic(err)
	}

	time.AfterFunc(time.Second, func() {
		// When pubsub is closed channel is closed too.
		_ = pubsub.Close()
	})

	// Consume messages.
	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
	}
}

func main() {
	c := make(chan struct{}, 0)

	println("WASM Go Inialized")

	var wg sync.WaitGroup
	wg.Add(1)

	// register functions
	registerCallbacks()
	<-c
}


