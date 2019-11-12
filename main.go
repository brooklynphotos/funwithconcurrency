package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	helloWorldSvr := getHellowWorldServer()
	helloNameSvr := getHelloNameServer()
	echoSvr := getEchoServer()

	helloWorldSvr.ListenAndServe()
	helloNameSvr.ListenAndServe()
	echoSvr.ListenAndServe()
	fmt.Println("all servers are started")
}

func getHellowWorldServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`Hello, world!`))
	})

	return &http.Server{Addr: ":7000", Handler: mux}
}

func getHelloNameServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		name := params.Get("name")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Hello, %s!", name)))
	})

	return &http.Server{Addr: ":8000", Handler: mux}
}

func getEchoServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.Copy(w, r.Body)
	})

	return &http.Server{Addr: ":9000", Handler: mux}
}
