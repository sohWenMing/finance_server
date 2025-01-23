package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sohWenMing/finance_server/config"
	database "github.com/sohWenMing/finance_server/internal/database/connection"
	programio "github.com/sohWenMing/finance_server/program_io"
	"github.com/sohWenMing/finance_server/server"
)

func main() {

	db, err := database.ConnectToDB(".env")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// init the server, send up channels

	config := config.Config{}
	config.RegisterQueries(db)
	//init the config, register all the queries and connect to the database

	portChan := make(chan int)
	doneChan := make(chan struct{})
	exitChan := make(chan struct{})

	go server.InitServer(false, portChan, doneChan, exitChan, http.Dir("."), config)
	port := <-portChan

	fmt.Printf("server started: listening on port %d\n", port)
	go programio.InitStdoutExit(doneChan)
	<-exitChan
	fmt.Print("program successfully exited")
}
