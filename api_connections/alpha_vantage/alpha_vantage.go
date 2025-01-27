package alphavantage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func GetBalanceSheetInformation(ticker string, client http.Client, ctx context.Context, apiKey string) (relevantMetrics []QuarterlyMetric, err error) {
	requestPath := fmt.Sprintf("https://www.alphavantage.co/query?function=BALANCE_SHEET&symbol=%s&apikey=%s", ticker, apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestPath, nil)
	if err != nil {
		return relevantMetrics, err
	}
	res, err := client.Do(req)
	if err != nil {
		return relevantMetrics, err
	}
	if res.StatusCode != 200 {
		return relevantMetrics, errors.New("status did not return with 200")
	}

	var financialReportJSON FinancialReport
	decoder := json.NewDecoder(res.Body)
	jsonDecodeErr := decoder.Decode(&financialReportJSON)
	if jsonDecodeErr != nil {
		return relevantMetrics, jsonDecodeErr
	}
	if len(financialReportJSON.QuarterlyReports) == 0 {
		return relevantMetrics, fmt.Errorf("no quarterly reports for %s returned", ticker)
		/*
			 			this will most likely be the return error in the event that there was an error.
						Issue with dealing with status is that even in the event of rate limiting or
						no information, 200 status will be returned but no relevant infomration will be returned
		*/
	}
	quarterlyReports := financialReportJSON.QuarterlyReports
	for _, quarterlyReport := range quarterlyReports {
		quarterlyMetric, err := MapQtrReportToRelevantMetrics(ticker, quarterlyReport)
		if err != nil {
			return []QuarterlyMetric{}, err
		}
		relevantMetrics = append(relevantMetrics, quarterlyMetric)
	}
	return relevantMetrics, nil

}
