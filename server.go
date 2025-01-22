package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/sohWenMing/finance_server/handlers"
)

func InitServer(isTest bool, portChan chan<- int, doneChan <-chan struct{}) {
	port := ":8080"
	if isTest {
		port = ":0"
	}
	//for testing purposes, usually will listen at port :8080, but if testing sets port to 0 so that random available port will be used

	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", handlers.PingHandler)
	mux.Handle("GET /app/", http.StripPrefix("/app", handlers.MiddleWareGenerator(http.FileServer(http.Dir(".")))))
	server := &http.Server{
		Handler: mux,
	}
	listener, err := net.Listen("tcp", port)
	//becauase net.Listen is already given an input param of "tcp", the return of listener.Addr() will will return a type of *net.TCPAddr
	if err != nil {
		log.Fatal(err)
	}
	returnedPort := listener.Addr().(*net.TCPAddr).Port
	// cast the address of the listener to type *net.TCPAddr, which has a Port field

	portChan <- returnedPort

	go func() {
		serveErr := server.Serve(listener)
		if serveErr != nil {
			if errors.Is(serveErr, http.ErrServerClosed) {
				fmt.Println("expected server close")
			}
			fmt.Print("serveErr was executed!\n")
			log.Fatal(serveErr)
		}
	}()
	// serving of the server has to be done it a separate go routine, because it is a blocking action
	<-doneChan
	//block, until server receives signal to close
	server.Close()
}
