package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

func main() {
	var counter int32

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		atomic.AddInt32(&counter, 1)
		log.Println("Current counter", counter)
		time.Sleep(10 * time.Millisecond)
		bytes, _ := ioutil.ReadAll(request.Body)
		writer.WriteHeader(http.StatusOK)
		writer.Write(bytes)
		atomic.AddInt32(&counter, -1)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
