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

	createUserPath := fmt.Sprintf("%s/createUser", basePath)
	createUserReqBody, responseJson, shouldReturn := createAndRetrieveUser(t, createUserPath)
	if shouldReturn {
		return
	}
	userIdUUID, err := uuid.Parse(responseJson.UserId)
	defer testConfig.Queries.DeleteUserById(context.Background(), userIdUUID)
	testhelpers.AssertNoError(t, err)
	retrievedUser, err := testConfig.Queries.GetUserById(context.Background(), userIdUUID)
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertStringVals(t, retrievedUser.Email, createUserReqBody.Email)
}

func TestLoginUserHandler(t *testing.T) {
	createUserPath := fmt.Sprintf("%s/createUser", basePath)

	type testStruct struct {
		name         string
		isExpectPass bool
	}

	tests := []testStruct{
		{
			name:         "test login user should pass",
			isExpectPass: true,
		},
		{
			name:         "test login user should fail",
			isExpectPass: false,
		},
	}
	for _, test := range tests {
		runLoginTest(t, createUserPath, test.isExpectPass)
	}
}

func runLoginTest(t *testing.T, createUserPath string, isExpectPass bool) {
	createUserReqBody, responseJson, shouldReturn := createAndRetrieveUser(t, createUserPath)
	if shouldReturn {
		return
	}

	userIdUUID, err := uuid.Parse(responseJson.UserId)
	testhelpers.AssertNoError(t, err)
	if err != nil {
		return
	}
	defer testConfig.Queries.DeleteUserById(context.Background(), userIdUUID)
	loginUserPath := fmt.Sprintf("%s/loginUser", basePath)
	if !isExpectPass {
		createUserReqBody.Password = createUserReqBody.Password + "add fail"
	}
	bodyString, err := json.Marshal(createUserReqBody)
	testhelpers.AssertNoError(t, err)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, loginUserPath, bytes.NewReader(bodyString))
	testhelpers.AssertNoError(t, err)
	if err != nil {
		return
	}
	res, err := client.Do(req)
	testhelpers.AssertNoError(t, err)
	if err != nil {
		return
	}

	switch isExpectPass {
	case true:
		testhelpers.AssertIntVals(t, res.StatusCode, 200)
		var loginResponseJson usermapping.LoginResponse
		decoder := json.NewDecoder(res.Body)
		decodeErr := decoder.Decode(&loginResponseJson)
		testhelpers.AssertNoError(t, decodeErr)
		if decodeErr != nil {
			return
		}
		testhelpers.AssertStringVals(t, loginResponseJson.AccessToken, "placeholder_access_token")
		testhelpers.AssertBool(t, loginResponseJson.IsSuccess, true)
		return
	case false:
		testhelpers.AssertIntVals(t, res.StatusCode, 401)
	}
}

func createAndRetrieveUser(t *testing.T, path string) (usermapping.UserJSON, usermapping.CreatedUserResponse, bool) {
	createUserReqBody := usermapping.UserJSON{
		Email:    "wenming.soh@gmail.com",
		Password: "password",
	}

	bodyString, err := json.Marshal(createUserReqBody)
	testhelpers.AssertNoError(t, err)

	req, reqErr := http.NewRequest(http.MethodPost, path, bytes.NewReader(bodyString))
	testhelpers.AssertNoError(t, reqErr)
	if reqErr != nil {
		return usermapping.UserJSON{}, usermapping.CreatedUserResponse{}, true
	}
	res, resErr := client.Do(req)
	testhelpers.AssertNoError(t, resErr)
	if resErr != nil {
		return usermapping.UserJSON{}, usermapping.CreatedUserResponse{}, true
	}
	var responseJson usermapping.CreatedUserResponse

	decoder := json.NewDecoder(res.Body)
	jsonDecodeErr := decoder.Decode(&responseJson)

	testhelpers.AssertNoError(t, jsonDecodeErr)
	if jsonDecodeErr != nil {
		return usermapping.UserJSON{}, usermapping.CreatedUserResponse{}, true
	}
	testhelpers.AssertBool(t, responseJson.IsSuccess, true)
	return createUserReqBody, responseJson, false
}
