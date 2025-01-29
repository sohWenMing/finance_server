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
)

var TestConfig = config.Config{}
var client = http.Client{
	Timeout: 30 * time.Second,
}
var testContext = context.Background()
var apiKey string

func TestMain(m *testing.M) {
	envvars.LoadEnv("../.env")
	apiKey = os.Getenv("ALPHA_VANTAGE_API_KEY")
	fmt.Printf("apiKey: %s\n", apiKey)

	db, err := database.ConnectToDB("../.env")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	TestConfig.RegisterQueries(db)
	code := m.Run()
	os.Exit(code)

}
