package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	address := "127.0.0.1:8080"

	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandle)
	mux.HandleFunc("/api", mainApiHandle)

	fmt.Println("Start server on", address)
	err := http.ListenAndServe(address, mux)
	if err != nil {
		panic(err)
	}
}

func mainHandle(response http.ResponseWriter, request *http.Request) {

	if request.Method == http.MethodGet {

		body := "Hello. I am Server.\r\n"
		body += fmt.Sprintf("Method: %s\r\n", request.Method)
		body += "Headers:\r\n"
		for key, header := range request.Header {
			body += fmt.Sprintf("- %s: %v\r\n", key, header)
		}
		body += "Query parameters:\r\n"
		for key, value := range request.URL.Query() {
			body += fmt.Sprintf("- %s: %v\r\n", key, value)
		}
		response.Write([]byte(body))
		return
	}

	if request.Method == http.MethodPost {

		body := "Hello. I am Server.\r\n"
		body += fmt.Sprintf("Method: %s\r\n", request.Method)
		body += "Headers:\r\n"
		for key, header := range request.Header {
			body += fmt.Sprintf("- %s: %v\r\n", key, header)
		}
		body += "Query parameters:\r\n"
		for key, value := range request.URL.Query() {
			body += fmt.Sprintf("- %s: %v\r\n", key, value)
		}
		body += "Form parameters:\r\n"
		if err := request.ParseForm(); err != nil {
			response.Write([]byte(err.Error()))
			return
		}
		for key, value := range request.Form {
			body += fmt.Sprintf("- %s: %v\r\n", key, value)
		}
		body += "Body content:\r\n"
		content, err := io.ReadAll(request.Body)
		if err != nil {
			response.Write([]byte(err.Error()))
			return
		}
		body += string(content)

		response.Write([]byte(body))
		return
	}

	http.Error(response, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
}

func mainApiHandle(response http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodGet {
		http.Error(response, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	data := []byte("Hello. I am Server API.")
	response.Write(data)
}
