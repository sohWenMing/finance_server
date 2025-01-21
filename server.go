package server

import (
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

	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", handlers.PingHandler)
	server := &http.Server{
		Handler: mux,
	}
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	returnedPort := listener.Addr().(*net.TCPAddr).Port
	// cast the address of the listener to type *net.TCPAddr, which has a Port field

	portChan <- returnedPort

	go func() {
		serveErr := server.Serve(listener)
		if serveErr != nil {
			log.Fatal(serveErr)
		}
	}()
	// serving of the server has to be done it a separate go routine, because it is a blocking action
	<-doneChan
	//block, until server receives signal to close
	server.Close()
}
