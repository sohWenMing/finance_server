package server

import (
	"log"
	"net/http"
	"os"
)

func InitServer() (returnedServer *http.Server, returnedPort string) {
	port := loadPort()
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", pingHandler)
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("error during server operation: %v\n", err)
		}
	}()
	return server, port
}

func loadPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":8080"
	}
	return port
}
