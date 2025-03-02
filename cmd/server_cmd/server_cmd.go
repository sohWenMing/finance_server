package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sohWenMing/finance_project/internal/server"
)

func main() {

	quitChan := make(chan (os.Signal), 1)
	signal.Notify(quitChan, syscall.SIGTERM, syscall.SIGINT)
	/*
		signal.Notify is useful to log when user exits through terminal, using something like control C to exit the operation
	*/

	server, port := server.InitServer()
	//server.InitServer() has a go routine within it that starts the server to make it listen in the background

	fmt.Printf("server is listening on: %s\n", port)

	<-quitChan
	//make the function wait for a quit signal to be sent to carry on the shutdown of the server in the main function

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	// cancelFunc, though not explicitly used in operation, should be called at end of function to avoid context leak
	err := server.Shutdown(ctx)
	if err != nil {
		fmt.Printf("error occured during shutdown: %v\n", err)
	}
	fmt.Println("server shutdown successfully")
}
