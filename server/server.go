package server

import "net/http"

func Serve(){
	http.Handle("/", http.FileServer(http.Dir("./server/wasm")))
	http.ListenAndServe(":1984", nil)
}
