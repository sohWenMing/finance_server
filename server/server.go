package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/sohWenMing/finance_server/config"
	"github.com/sohWenMing/finance_server/handlers"
)

func InitServer(isTest bool, portChan chan<- int, doneChan <-chan struct{}, exitChan chan<- struct{}, root http.FileSystem, config config.Config) {
	port := ":8080"
	if isTest {
		port = ":0"
	}
	//for testing purposes, usually will listen at port :8080, but if testing sets port to 0 so that random available port will be used

	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", handlers.PingHandler)
	mux.Handle("GET /testJWTAccess", handlers.CheckAccessMiddleware(handlers.PingHandler, config.JwtSecret, false))

	mux.Handle("GET /app/", http.StripPrefix("/app", handlers.FileServerMiddleWare(http.FileServer(root))))
	mux.HandleFunc("POST /createUser", handlers.CreateUserHandler(config.Queries))
	mux.HandleFunc("POST /loginUser", handlers.LoginUserHandler(config.GetJWTValidDuration, config.Queries, config.JwtSecret))
	mux.HandleFunc("POST /refreshToken", handlers.CheckRefreshTokenAndGetNewJWTHandler(config.GetJWTValidDuration, config.Queries, config.JwtSecret))
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
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			log.Fatal(serveErr)
		}
	}()
	// serving of the server has to be done it a separate go routine, because it is a blocking action
	<-doneChan
	fmt.Println("signal was sent on channel, closing server")
	//block, until server receives signal to close
	server.Close()
	exitChan <- struct{}{}

}
