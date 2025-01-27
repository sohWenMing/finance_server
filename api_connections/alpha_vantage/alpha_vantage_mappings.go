package alphavantage

import (
	"time"

	"github.com/sohWenMing/finance_server/conversions"
)

type FinancialReport struct {
	Symbol        string `json:"symbol"`
	AnnualReports []struct {
		FiscalDateEnding                       string `json:"fiscalDateEnding"`
		ReportedCurrency                       string `json:"reportedCurrency"`
		TotalAssets                            string `json:"totalAssets"`
		TotalCurrentAssets                     string `json:"totalCurrentAssets"`
		CashAndCashEquivalentsAtCarryingValue  string `json:"cashAndCashEquivalentsAtCarryingValue"`
		CashAndShortTermInvestments            string `json:"cashAndShortTermInvestments"`
		Inventory                              string `json:"inventory"`
		CurrentNetReceivables                  string `json:"currentNetReceivables"`
		TotalNonCurrentAssets                  string `json:"totalNonCurrentAssets"`
		PropertyPlantEquipment                 string `json:"propertyPlantEquipment"`
		AccumulatedDepreciationAmortizationPPE string `json:"accumulatedDepreciationAmortizationPPE"`
		IntangibleAssets                       string `json:"intangibleAssets"`
		IntangibleAssetsExcludingGoodwill      string `json:"intangibleAssetsExcludingGoodwill"`
		Goodwill                               string `json:"goodwill"`
		Investments                            string `json:"investments"`
		LongTermInvestments                    string `json:"longTermInvestments"`
		ShortTermInvestments                   string `json:"shortTermInvestments"`
		OtherCurrentAssets                     string `json:"otherCurrentAssets"`
		OtherNonCurrentAssets                  string `json:"otherNonCurrentAssets"`
		TotalLiabilities                       string `json:"totalLiabilities"`
		TotalCurrentLiabilities                string `json:"totalCurrentLiabilities"`
		CurrentAccountsPayable                 string `json:"currentAccountsPayable"`
		DeferredRevenue                        string `json:"deferredRevenue"`
		CurrentDebt                            string `json:"currentDebt"`
		ShortTermDebt                          string `json:"shortTermDebt"`
		TotalNonCurrentLiabilities             string `json:"totalNonCurrentLiabilities"`
		CapitalLeaseObligations                string `json:"capitalLeaseObligations"`
		LongTermDebt                           string `json:"longTermDebt"`
		CurrentLongTermDebt                    string `json:"currentLongTermDebt"`
		LongTermDebtNoncurrent                 string `json:"longTermDebtNoncurrent"`
		ShortLongTermDebtTotal                 string `json:"shortLongTermDebtTotal"`
		OtherCurrentLiabilities                string `json:"otherCurrentLiabilities"`
		OtherNonCurrentLiabilities             string `json:"otherNonCurrentLiabilities"`
		TotalShareholderEquity                 string `json:"totalShareholderEquity"`
		TreasuryStock                          string `json:"treasuryStock"`
		RetainedEarnings                       string `json:"retainedEarnings"`
		CommonStock                            string `json:"commonStock"`
		CommonStockSharesOutstanding           string `json:"commonStockSharesOutstanding"`
	} `json:"annualReports"`
	QuarterlyReports []struct {
		FiscalDateEnding                       string `json:"fiscalDateEnding"`
		ReportedCurrency                       string `json:"reportedCurrency"`
		TotalAssets                            string `json:"totalAssets"`
		TotalCurrentAssets                     string `json:"totalCurrentAssets"`
		CashAndCashEquivalentsAtCarryingValue  string `json:"cashAndCashEquivalentsAtCarryingValue"`
		CashAndShortTermInvestments            string `json:"cashAndShortTermInvestments"`
		Inventory                              string `json:"inventory"`
		CurrentNetReceivables                  string `json:"currentNetReceivables"`
		TotalNonCurrentAssets                  string `json:"totalNonCurrentAssets"`
		PropertyPlantEquipment                 string `json:"propertyPlantEquipment"`
		AccumulatedDepreciationAmortizationPPE string `json:"accumulatedDepreciationAmortizationPPE"`
		IntangibleAssets                       string `json:"intangibleAssets"`
		IntangibleAssetsExcludingGoodwill      string `json:"intangibleAssetsExcludingGoodwill"`
		Goodwill                               string `json:"goodwill"`
		Investments                            string `json:"investments"`
		LongTermInvestments                    string `json:"longTermInvestments"`
		ShortTermInvestments                   string `json:"shortTermInvestments"`
		OtherCurrentAssets                     string `json:"otherCurrentAssets"`
		OtherNonCurrentAssets                  string `json:"otherNonCurrentAssets"`
		TotalLiabilities                       string `json:"totalLiabilities"`
		TotalCurrentLiabilities                string `json:"totalCurrentLiabilities"`
		CurrentAccountsPayable                 string `json:"currentAccountsPayable"`
		DeferredRevenue                        string `json:"deferredRevenue"`
		CurrentDebt                            string `json:"currentDebt"`
		ShortTermDebt                          string `json:"shortTermDebt"`
		TotalNonCurrentLiabilities             string `json:"totalNonCurrentLiabilities"`
		CapitalLeaseObligations                string `json:"capitalLeaseObligations"`
		LongTermDebt                           string `json:"longTermDebt"`
		CurrentLongTermDebt                    string `json:"currentLongTermDebt"`
		LongTermDebtNoncurrent                 string `json:"longTermDebtNoncurrent"`
		ShortLongTermDebtTotal                 string `json:"shortLongTermDebtTotal"`
		OtherCurrentLiabilities                string `json:"otherCurrentLiabilities"`
		OtherNonCurrentLiabilities             string `json:"otherNonCurrentLiabilities"`
		TotalShareholderEquity                 string `json:"totalShareholderEquity"`
		TreasuryStock                          string `json:"treasuryStock"`
		RetainedEarnings                       string `json:"retainedEarnings"`
		CommonStock                            string `json:"commonStock"`
		CommonStockSharesOutstanding           string `json:"commonStockSharesOutstanding"`
	} `json:"quarterlyReports"`
}

type QuarterlyReportStruct struct {
	FiscalDateEnding                       string `json:"fiscalDateEnding"`
	ReportedCurrency                       string `json:"reportedCurrency"`
	TotalAssets                            string `json:"totalAssets"`
	TotalCurrentAssets                     string `json:"totalCurrentAssets"`
	CashAndCashEquivalentsAtCarryingValue  string `json:"cashAndCashEquivalentsAtCarryingValue"`
	CashAndShortTermInvestments            string `json:"cashAndShortTermInvestments"`
	Inventory                              string `json:"inventory"`
	CurrentNetReceivables                  string `json:"currentNetReceivables"`
	TotalNonCurrentAssets                  string `json:"totalNonCurrentAssets"`
	PropertyPlantEquipment                 string `json:"propertyPlantEquipment"`
	AccumulatedDepreciationAmortizationPPE string `json:"accumulatedDepreciationAmortizationPPE"`
	IntangibleAssets                       string `json:"intangibleAssets"`
	IntangibleAssetsExcludingGoodwill      string `json:"intangibleAssetsExcludingGoodwill"`
	Goodwill                               string `json:"goodwill"`
	Investments                            string `json:"investments"`
	LongTermInvestments                    string `json:"longTermInvestments"`
	ShortTermInvestments                   string `json:"shortTermInvestments"`
	OtherCurrentAssets                     string `json:"otherCurrentAssets"`
	OtherNonCurrentAssets                  string `json:"otherNonCurrentAssets"`
	TotalLiabilities                       string `json:"totalLiabilities"`
	TotalCurrentLiabilities                string `json:"totalCurrentLiabilities"`
	CurrentAccountsPayable                 string `json:"currentAccountsPayable"`
	DeferredRevenue                        string `json:"deferredRevenue"`
	CurrentDebt                            string `json:"currentDebt"`
	ShortTermDebt                          string `json:"shortTermDebt"`
	TotalNonCurrentLiabilities             string `json:"totalNonCurrentLiabilities"`
	CapitalLeaseObligations                string `json:"capitalLeaseObligations"`
	LongTermDebt                           string `json:"longTermDebt"`
	CurrentLongTermDebt                    string `json:"currentLongTermDebt"`
	LongTermDebtNoncurrent                 string `json:"longTermDebtNoncurrent"`
	ShortLongTermDebtTotal                 string `json:"shortLongTermDebtTotal"`
	OtherCurrentLiabilities                string `json:"otherCurrentLiabilities"`
	OtherNonCurrentLiabilities             string `json:"otherNonCurrentLiabilities"`
	TotalShareholderEquity                 string `json:"totalShareholderEquity"`
	TreasuryStock                          string `json:"treasuryStock"`
	RetainedEarnings                       string `json:"retainedEarnings"`
	CommonStock                            string `json:"commonStock"`
	CommonStockSharesOutstanding           string `json:"commonStockSharesOutstanding"`
}

type QuarterlyMetric struct {
	Ticker                       string    `json:"ticker"`
	FiscalDateEnding             time.Time `json:"fiscal_date_ending"`
	TotalAssets                  int64     `json:"total_assets"`
	IntangibleAssets             int64     `json:"intangible_assets"`
	TotalLiabilities             int64     `json:"total_liabilities"`
	CommonStock                  int64     `json:"common_stock"`
	CommonStockSharesOutstanding int64     `json:"common_stock_shares_outstanding"`
	TangibleBookValuePerShare    float32   `json:"tangible_book_value_per_share"`
}

func MapQtrReportToRelevantMetrics(ticker string, qtrReport QuarterlyReportStruct) (relevantMetrics QuarterlyMetric, err error) {

	timeStamp, err := conversions.GetDateOnlyTimeStampFromDateString(qtrReport.FiscalDateEnding)
	if err != nil {
		return relevantMetrics, err
	}

	totalAssets, err := conversions.GetUint64FromString(qtrReport.TotalAssets)
	if err != nil {
		return relevantMetrics, err
	}
	intangibleAssets, err := conversions.GetUint64FromString(qtrReport.IntangibleAssets)
	if err != nil {
		return relevantMetrics, err
	}
	totalLiabilities, err := conversions.GetUint64FromString(qtrReport.TotalLiabilities)
	if err != nil {
		return relevantMetrics, err
	}
	commonStock, err := conversions.GetUint64FromString(qtrReport.CommonStock)
	if err != nil {
		return relevantMetrics, err
	}
	commonStockSharesOutstanding, err := conversions.GetUint64FromString(qtrReport.CommonStockSharesOutstanding)
	if err != nil {
		return relevantMetrics, err
	}
	workingRelevantMetrics := QuarterlyMetric{
		Ticker:                       ticker,
		FiscalDateEnding:             timeStamp,
		TotalAssets:                  totalAssets,
		IntangibleAssets:             intangibleAssets,
		TotalLiabilities:             totalLiabilities,
		CommonStock:                  commonStock,
		CommonStockSharesOutstanding: commonStockSharesOutstanding,
	}
	return getTBVPS(workingRelevantMetrics), nil
}

func getTBVPS(input QuarterlyMetric) (output QuarterlyMetric) {
	workingRelevantMetrics := input
	workingRelevantMetrics.TangibleBookValuePerShare = (float32(input.TotalAssets) - float32(input.IntangibleAssets) - float32(input.TotalLiabilities)) / float32(input.CommonStockSharesOutstanding)
	return workingRelevantMetrics
}
