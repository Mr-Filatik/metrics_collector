package main

import (
	"fmt"
	"net/http"
)

func main() {
	address := "127.0.0.1:8080"
	var handler Handler

	fmt.Println("Start server on", address)
	err := http.ListenAndServe(address, handler)
	if err != nil {
		panic(err)
	}
}

type Handler struct{}

func (h Handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	data := []byte("Hello. I am Server.")
	response.Write(data)
}
