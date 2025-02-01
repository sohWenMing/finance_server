package integrationtesting

import (
	"testing"
	"time"

	"github.com/google/uuid"
	alphavantage "github.com/sohWenMing/finance_server/api_connections/alpha_vantage"
	"github.com/sohWenMing/finance_server/internal/database/sqlc_generated"
	testhelpers "github.com/sohWenMing/finance_server/test_helpers"
)

func TestAlphaVantageDBIntegrationDemo(t *testing.T) {
	relevantMetrics, err := alphavantage.GetBalanceSheetInformation("IBM", client, testContext, "demo")
	numRecords := 0
	testhelpers.AssertNoError(t, err)
	for _, metric := range relevantMetrics {
		param := mapRelevantMetricParam(metric)
		_, err := TestConfig.Queries.CreateBalanceSheetRecord(testContext, param)
		testhelpers.AssertNoError(t, err)
		if err == nil {
			numRecords += 1
		}

	}
	retrievedRecords, err := TestConfig.Queries.GetBalanceSheetRecordsByTicker(testContext, "IBM")
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertIntVals(t, len(retrievedRecords), numRecords)
	deleteErr := TestConfig.Queries.DeleteBalanceSheetRecordsByTicker(testContext, "IBM")
	testhelpers.AssertNoError(t, deleteErr)
	existingRecords, err := TestConfig.Queries.GetBalanceSheetRecordsByTicker(testContext, "IBM")
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertIntVals(t, len(existingRecords), 0)
}

// func TestAlphaVantageDBIntegrationActual(t *testing.T) {
// 	tickers := []string{
// 		// "C",
// 		// "BAC",
// 		// "JPM",
// 	}
// 	//set the tickers that the api will call

// 	for _, ticker := range tickers {
// 		fmt.Printf("retrieving information for: %s\n", ticker)
// 		numRecords := 0
// 		relevantMetrics, err := alphavantage.GetBalanceSheetInformation(ticker, client, testContext, apiKey)
// 		testhelpers.AssertNoError(t, err)
// 		if err != nil {
// 			continue
// 		}
// 		//make the api call to get the information for the ticker and expect no error

// 		for _, metric := range relevantMetrics {
// 			param := mapRelevantMetricParam(metric)
// 			_, err := testConfig.Queries.CreateBalanceSheetRecord(testContext, param)
// 			//for each record, map and create the Param for the database and attempt to write to DB

// 			testhelpers.AssertNoError(t, err)
// 			if err != nil {
// 				continue
// 			}
// 			numRecords += 1
// 		}
// 		retrievedRecords, err := testConfig.Queries.GetBalanceSheetRecordsByTicker(testContext, ticker)
// 		testhelpers.AssertNoError(t, err)
// 		testhelpers.AssertIntVals(t, len(retrievedRecords), numRecords)
// 		time.Sleep(10 * time.Second)
// 	}
//}

func mapRelevantMetricParam(metric alphavantage.QuarterlyMetric) (param sqlc_generated.CreateBalanceSheetRecordParams) {
	mappedParam := sqlc_generated.CreateBalanceSheetRecordParams{
		ID:                           uuid.New(),
		Ticker:                       metric.Ticker,
		FiscalDateEnding:             metric.FiscalDateEnding,
		TotalAssets:                  metric.TotalAssets,
		IntangibleAssets:             metric.IntangibleAssets,
		TotalLiabilities:             metric.TotalLiabilities,
		CommonStock:                  metric.CommonStock,
		CommonStockSharesOutstanding: metric.CommonStockSharesOutstanding,
		CreatedOn:                    time.Now(),
		UpdatedOn:                    time.Now(),
	}
	return mappedParam

}
