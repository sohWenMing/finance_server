package alphavantage

import (
	"context"
	"net/http"
	"testing"
	"time"

	testhelpers "github.com/sohWenMing/finance_server/test_helpers"
)

var client = http.Client{
	Timeout: 30 * time.Second,
}

func TestGetBalanceSheetInformation(t *testing.T) {

	type relevantMetricsToErr struct {
		ticker          string
		relevantMetrics []QuarterlyMetric
		err             error
	}

	results := []relevantMetricsToErr{}
	tickers := []string{
		"IBM",
	}
	for _, ticker := range tickers {
		result, err := GetBalanceSheetInformation(ticker, client, context.Background(), "demo")
		resultToErr := relevantMetricsToErr{
			ticker,
			result,
			err,
		}
		results = append(results, resultToErr)
		time.Sleep(1 * time.Second)
	}

	for _, result := range results {
		testhelpers.AssertStringVals(t, result.relevantMetrics[0].Ticker, "IBM")
		testhelpers.AssertNoError(t, result.err)
	}
}
