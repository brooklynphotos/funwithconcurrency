package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// some context to let the servers know if we are canceling
	ctx, cancel := context.WithCancel(context.Background())

	go getHellowWorldServer(ctx)
	helloNameSvr := getHelloNameServer()
	echoSvr := getEchoServer()

	go helloNameSvr.ListenAndServe()
	go echoSvr.ListenAndServe()
	fmt.Println("all servers are started")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	cancel()
}

func getHellowWorldServer(ctx context.Context) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`Hello, world!`))
	})

	server := &http.Server{Addr: ":7000", Handler: mux}

	// shutdown using a context
	go func() {
		<-ctx.Done() // this is called when the corresponding cancel function is called
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutCtx); err != nil {
			fmt.Printf("error shutting down hello world: %s\n", err)
		}
	}()
	fmt.Printf("Hello world is starting")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("error starting hello world server: %s\n", err)
	}
	fmt.Println("Hello server closing")
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
