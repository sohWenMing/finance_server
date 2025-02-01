package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sohWenMing/finance_server/config"
	database "github.com/sohWenMing/finance_server/internal/database/connection"
	programio "github.com/sohWenMing/finance_server/program_io"
	"github.com/sohWenMing/finance_server/server"
)

const k = ".env"

func main() {

	db, err := database.ConnectToDB(k)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB Connection started")
	defer db.Close()
	// init the server, send up channels

	config := config.InitConfig()
	config.SetJWTValidDuration(10 * time.Minute)
	//make jwt tokens in main function only valid for 10 minutes
	config.RegisterQueries(db)
	registerJWTSecretErr := config.RegisterJwtSecret(k)
	//init the config, register all the queries and connect to the database
	if registerJWTSecretErr != nil {
		log.Fatal(registerJWTSecretErr)
	}

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
