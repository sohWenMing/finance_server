package main

import (
	"fmt"
	"net/http"

	programio "github.com/sohWenMing/finance_server/program_io"
	"github.com/sohWenMing/finance_server/server"
)

func main() {
	// init the server, send up channels

	portChan := make(chan int)
	doneChan := make(chan struct{})
	exitChan := make(chan struct{})

	go server.InitServer(false, portChan, doneChan, exitChan, http.Dir("."))
	port := <-portChan

	fmt.Printf("server started: listening on port %d\n", port)
	go programio.InitStdoutExit(doneChan)
	<-exitChan
	fmt.Print("program successfully exited")
}
