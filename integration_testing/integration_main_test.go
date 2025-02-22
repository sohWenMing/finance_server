package integrationtesting

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/sohWenMing/finance_server/config"
	envvars "github.com/sohWenMing/finance_server/env_vars"
	database "github.com/sohWenMing/finance_server/internal/database/connection"
	"github.com/sohWenMing/finance_server/server"
)

var (
	portChan             chan int
	doneChan             chan struct{}
	exitChan             chan struct{}
	receivedPort         int
	basePath             string
	createUserPath       string
	loginUserPath        string
	testJWtAccessPath    string
	testRefreshTokenPath string
)

var TestConfig = config.InitConfig()

//don't set the duration for JWT validity here, set it within the tests so that can check positive and negative cases

var client = http.Client{
	Timeout: 30 * time.Second,
}
var testContext = context.Background()
var apiKey string

func TestMain(m *testing.M) {
	envvars.LoadEnv("../.env")
	apiKey = os.Getenv("ALPHA_VANTAGE_API_KEY")

	portChan = make(chan int)
	doneChan = make(chan struct{})
	exitChan = make(chan struct{})

	db, err := database.ConnectToDB("../.env")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	TestConfig.RegisterQueries(db)

	TestConfig.RegisterJwtSecret("./.env")

	go func(portChan chan int, doneChan chan struct{}) {
		server.InitServer(true, portChan, doneChan, exitChan, http.Dir(".."), TestConfig)
	}(portChan, doneChan)
	//Init on server has to be done on separate goroutine, so as to not block

	receivedPort = <-portChan
	basePath = fmt.Sprintf("http://localhost:%d", receivedPort)
	createUserPath = fmt.Sprintf("%s/createUser", basePath)
	loginUserPath = fmt.Sprintf("%s/loginUser", basePath)
	testJWtAccessPath = fmt.Sprintf("%s/testJWTAccess", basePath)
	testRefreshTokenPath = fmt.Sprintf("%s/refreshToken", basePath)
	code := m.Run()
	os.Exit(code)

}
