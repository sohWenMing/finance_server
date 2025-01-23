package main

import (
	"fmt"
	"net/http"
	"os"

	envvars "github.com/sohWenMing/finance_server/env_vars"
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

	envVarErr := envvars.LoadEnv(".env")
	if envVarErr != nil {
		fmt.Printf("error: %v\n", envVarErr)
	}
	DbString := os.Getenv("DB_STRING")
	fmt.Printf("DbString: %s\n", DbString)

	fmt.Printf("server started: listening on port %d\n", port)
	go programio.InitStdoutExit(doneChan)
	<-exitChan
	fmt.Print("program successfully exited")
}
