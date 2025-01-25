package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sohWenMing/finance_server/config"
	database "github.com/sohWenMing/finance_server/internal/database/connection"
	usermapping "github.com/sohWenMing/finance_server/mapping/user_mapping"
	testhelpers "github.com/sohWenMing/finance_server/test_helpers"
)

var (
	portChan     chan int
	doneChan     chan struct{}
	exitChan     chan struct{}
	client       *http.Client
	receivedPort int
	basePath     string
)

var testConfig = config.Config{}

func TestMain(m *testing.M) {

	db, err := database.ConnectToDB("../.env")
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()
	portChan = make(chan int)
	doneChan = make(chan struct{})
	exitChan = make(chan struct{})
	client = &http.Client{
		Timeout: 5 * time.Second,
	}

	testConfig.RegisterQueries(db)

	go func(portChan chan int, doneChan chan struct{}) {
		InitServer(true, portChan, doneChan, exitChan, http.Dir(".."), testConfig)
	}(portChan, doneChan)
	//Init on server has to be done on separate goroutine, so as to not block

	receivedPort = <-portChan
	basePath = fmt.Sprintf("http://localhost:%d", receivedPort)
	//block execution until receivedPort is received from the portChan
	code := m.Run()
	doneChan <- struct{}{}
	//send done signal to server to close server, when all tests are done
	os.Exit(code)
}

func TestPing(t *testing.T) {

	path := fmt.Sprintf("%s/ping", basePath)

	req, err := http.NewRequest(http.MethodGet, path, nil)
	testhelpers.AssertNoError(t, err)

	res, err := client.Do(req)
	testhelpers.AssertNoError(t, err)

	testhelpers.AssertIntVals(t, res.StatusCode, 200)

}

func TestFileServer(t *testing.T) {

	path := fmt.Sprintf("%s/app/", basePath)

	req, err := http.NewRequest(http.MethodGet, path, nil)
	testhelpers.AssertNoError(t, err)

	res, err := client.Do(req)
	testhelpers.AssertNoError(t, err)

	resBody, err := io.ReadAll(res.Body)
	testhelpers.AssertNoError(t, err)
	if !strings.Contains(string(resBody), "let's get this") {
		t.Errorf("%s", string(resBody))
	}
}

func TestCreateUserHandler(t *testing.T) {

	path := fmt.Sprintf("%s/createUser", basePath)
	createUserReqBody := usermapping.CreateUserJSON{
		Email:    "wenming.soh@gmail.com",
		Password: "password",
	}

	bodyString, err := json.Marshal(createUserReqBody)
	testhelpers.AssertNoError(t, err)

	req, reqErr := http.NewRequest(http.MethodPost, path, bytes.NewReader(bodyString))
	testhelpers.AssertNoError(t, reqErr)
	if reqErr != nil {
		return
	}
	res, resErr := client.Do(req)
	testhelpers.AssertNoError(t, resErr)
	if resErr != nil {
		return
	}
	var responseJson usermapping.CreatedUserResponse

	decoder := json.NewDecoder(res.Body)
	jsonDecodeErr := decoder.Decode(&responseJson)

	testhelpers.AssertNoError(t, jsonDecodeErr)
	if jsonDecodeErr != nil {
		return
	}
	testhelpers.AssertBool(t, responseJson.IsSuccess, true)
	userIdUUID, err := uuid.Parse(responseJson.UserId)
	testhelpers.AssertNoError(t, err)
	retrievedUser, err := testConfig.Queries.GetUserById(context.Background(), userIdUUID)
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertStringVals(t, retrievedUser.Email, createUserReqBody.Email)
	testConfig.Queries.DeleteUserById(context.Background(), userIdUUID)
}
