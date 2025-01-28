package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sohWenMing/finance_server/config"
	database "github.com/sohWenMing/finance_server/internal/database/connection"
	"github.com/sohWenMing/finance_server/internal/database/sqlc_generated"
)

func main() {
	runConfig := config.Config{}
	// runContext := context.Background()

	db, err := database.ConnectToDB("../.env")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	runConfig.RegisterQueries(db)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ticker := scanner.Text()
		fmt.Printf("user entered %s\n", ticker)
		if ticker == "exit" {
			os.Exit(0)
		}
		entries, err := runConfig.Queries.GetBalanceSheetRecordsByTicker(context.Background(), ticker)
		if err != nil {
			fmt.Printf("error: %v", err)
		}
		fmt.Printf("num entries for ticker %s: %d\n", ticker, len(entries))
		if len(entries) == 0 {
			fmt.Printf("no entries for ticker %s: no file generated\n", ticker)
			continue
		}
		csvErr := generateCSVFile(ticker, entries)
		if csvErr != nil {
			fmt.Println(csvErr.Error())
		}

	}
}

func generateCSVFile(ticker string, entries []sqlc_generated.BalanceSheet) (err error) {

	timestamp := time.Now().Format(time.DateOnly)
	filename := fmt.Sprintf("%s_%s.csv", timestamp, ticker)
	file, err := os.Create(fmt.Sprintf("../generated_data/%s", filename))
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	headers := []string{
		"ticker",
		"fiscal_date_ending",
		"total_assets",
		"intangible_assets",
		"total_liabilities",
		"common_stock",
		"common_stock_shares_outstanding",
	}
	data := [][]string{}
	for _, entry := range entries {
		data = append(data, mapCSVRow(entry))
	}
	writer.Write(headers)
	for _, dataRow := range data {
		writer.Write(dataRow)
	}
	return nil
}

func mapCSVRow(input sqlc_generated.BalanceSheet) []string {
	row := []string{}
	row = append(row, input.Ticker)
	row = append(row, input.FiscalDateEnding.Format(time.DateOnly))
	row = append(row, fmt.Sprintf("%d, ", input.TotalAssets))
	row = append(row, fmt.Sprintf("%d, ", input.IntangibleAssets))
	row = append(row, fmt.Sprintf("%d, ", input.TotalLiabilities))
	row = append(row, fmt.Sprintf("%d, ", input.CommonStock))
	row = append(row, fmt.Sprintf("%d, ", input.CommonStockSharesOutstanding))
	return row
}
