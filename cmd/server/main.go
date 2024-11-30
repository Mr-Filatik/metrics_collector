package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	address := "127.0.0.1:8080"

	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandle)
	mux.HandleFunc("/json", mainJsonHandle)
	mux.HandleFunc("/api", mainApiHandle)
	mux.HandleFunc("/api/login", mainApiLoginHandle)
	mux.Handle("/api/account", Conveyor(http.HandlerFunc(mainApiAccountHandle), authMiddleware))

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

type MessageResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func mainJsonHandle(response http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodGet {
		http.Error(response, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	data := MessageResponse{"Test message.", "Test error."}
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}
	response.Header().Set("content-type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(jsonData)
}

func mainApiHandle(response http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodGet {
		http.Error(response, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	data := []byte("Hello. I am Server API.")
	response.Write(data)
}

const form = `<html>
    <head>
    <title></title>
    </head>
    <body>
        <form action="/api/login" method="post">
            <label>Логин <input type="text" name="login"></label>
            <label>Пароль <input type="password" name="password"></label>
            <input type="submit" value="Login">
        </form>
    </body>
</html>`

func Auth(login, password string) bool {
	return login == "user" && password == "password"
}

func mainApiLoginHandle(response http.ResponseWriter, request *http.Request) {

	if request.Method == http.MethodPost {
		login := request.FormValue("login")
		password := request.FormValue("password")
		if Auth(login, password) {
			io.WriteString(response, "Hello!")
		} else {
			http.Error(response, "Invalid login or password", http.StatusUnauthorized)
		}
	} else {
		io.WriteString(response, form)
	}
}

type Middleware func(http.Handler) http.Handler

func Conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("auth") == "user" {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Invalid auth token", http.StatusUnauthorized)
		}
	})
}

func mainApiAccountHandle(response http.ResponseWriter, request *http.Request) {

	response.Write([]byte("Account!"))
}
