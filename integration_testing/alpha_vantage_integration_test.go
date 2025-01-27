package integrationtesting

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	alphavantage "github.com/sohWenMing/finance_server/api_connections/alpha_vantage"
	"github.com/sohWenMing/finance_server/config"
	envvars "github.com/sohWenMing/finance_server/env_vars"
	database "github.com/sohWenMing/finance_server/internal/database/connection"
	testhelpers "github.com/sohWenMing/finance_server/test_helpers"
)

var testConfig = config.Config{}
var client = http.Client{
	Timeout: 30 * time.Second,
}
var testContext = context.Background()

func TestMain(m *testing.M) {
	envvars.LoadEnv("../.env")
	apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	fmt.Printf("apiKey: %s\n", apiKey)

	db, err := database.ConnectToDB("../.env")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	testConfig.RegisterQueries(db)
	code := m.Run()
	os.Exit(code)

}

func TestAlphaVantageDBIntegration(t *testing.T) {
	_, err := alphavantage.GetBalanceSheetInformation("IBM", client, testContext, "demo")
	numRecords := 0
	testhelpers.AssertNoError(t, err)
	// for _, metric := range relevantMetrics {
	// 	param := sqlc_generated.CreateBalanceSheetRecordParams{
	// 		ID:                           uuid.New(),
	// 		Ticker:                       "IBM",
	// 		FiscalDateEnding:             metric.FiscalDateEnding,
	// 		TotalAssets:                  metric.TotalAssets,
	// 		IntangibleAssets:             metric.IntangibleAssets,
	// 		TotalLiabilities:             metric.TotalLiabilities,
	// 		CommonStock:                  metric.CommonStock,
	// 		CommonStockSharesOutstanding: metric.CommonStockSharesOutstanding,
	// 		CreatedOn:                    time.Now(),
	// 		UpdatedOn:                    time.Now(),
	// 	}
	// 	record, err := testConfig.Queries.CreateBalanceSheetRecord(testContext, param)
	// 	fmt.Printf("created record: %v\n", record)
	// 	testhelpers.AssertNoError(t, err)
	// 	if err == nil {
	// 		numRecords += 1
	// 	}
	// }
	fmt.Printf("Number of records inserted into database: %d", numRecords)

}
