package server

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	testhelpers "github.com/sohWenMing/finance_server/test_helpers"
)

func TestPing(t *testing.T) {
	portChan := make(chan int)
	doneChan := make(chan struct{})
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	go func(portChan chan int, doneChan chan struct{}) {
		InitServer(true, portChan, doneChan)
	}(portChan, doneChan)
	//Init on server has to be done on separate goroutine, so as to not block

	receivedPort := <-portChan
	path := fmt.Sprintf("http://localhost:%d/ping", receivedPort)

	req, err := http.NewRequest(http.MethodGet, path, nil)
	testhelpers.AssertNoError(t, err)

	res, err := client.Do(req)
	testhelpers.AssertNoError(t, err)

	testhelpers.AssertIntVals(t, res.StatusCode, 200)
	doneChan <- struct{}{}

}
