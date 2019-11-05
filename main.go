package main

import "net/http"

func main() {
	helloWorldSvr := getHellowWorldServer()

	helloWorldSvr.ListenAndServe()
}

func getHellowWorldServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`Hello, world!`))
	})

	return &http.Server{Addr: ":7000", Handler: mux}
}
