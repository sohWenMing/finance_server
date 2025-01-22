package server

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	testhelpers "github.com/sohWenMing/finance_server/test_helpers"
)

var (
	portChan     chan int
	doneChan     chan struct{}
	client       *http.Client
	receivedPort int
)

func TestMain(m *testing.M) {

	portChan = make(chan int)
	doneChan = make(chan struct{})
	client = &http.Client{
		Timeout: 5 * time.Second,
	}

	go func(portChan chan int, doneChan chan struct{}) {
		InitServer(true, portChan, doneChan)
	}(portChan, doneChan)
	//Init on server has to be done on separate goroutine, so as to not block

	receivedPort = <-portChan
	//block execution until receivedPort is received from the portChan
	code := m.Run()
	doneChan <- struct{}{}
	os.Exit(code)
}

func TestPing(t *testing.T) {

	path := fmt.Sprintf("http://localhost:%d/ping", receivedPort)

	req, err := http.NewRequest(http.MethodGet, path, nil)
	testhelpers.AssertNoError(t, err)

	res, err := client.Do(req)
	testhelpers.AssertNoError(t, err)

	testhelpers.AssertIntVals(t, res.StatusCode, 200)
	doneChan <- struct{}{}

}
