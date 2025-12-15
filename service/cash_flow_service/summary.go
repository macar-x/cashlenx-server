package cash_flow_service

import (
	"errors"
	"strings"
	"time"

	"github.com/macar-x/cashlenx-server/mapper/cash_flow_mapper"
	"github.com/macar-x/cashlenx-server/mapper/category_mapper"
	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/util"
)

// Summary represents financial summary data
type Summary struct {
	TotalIncome       float64
	TotalExpense      float64
	Balance           float64
	TransactionCount  int
	CategoryBreakdown map[string]float64
}

// GetSummary returns financial summary for a given period
func GetSummary(period, date string) (*Summary, error) {
	validPeriods := map[string]bool{
		"daily":   true,
		"monthly": true,
		"yearly":  true,
	}

	if !validPeriods[period] {
		return nil, errors.New("invalid period: must be daily, monthly, or yearly")
	}

	var fromDate, toDate time.Time
	var err error

	// Parse date based on period
	switch period {
	case "daily":
		// Date format: YYYY-MM-DD
		fromDate = util.FormatDateFromStringWithoutDash(date)
		if fromDate.IsZero() {
			return nil, errors.New("invalid date format for daily, use YYYY-MM-DD")
		}
		toDate = fromDate
	case "monthly":
		// Date format: YYYY-MM
		parts := strings.Split(date, "-")
		if len(parts) != 2 {
			return nil, errors.New("invalid date format for monthly, use YYYY-MM")
		}
		fromDate, err = time.Parse("2006-01", date)
		if err != nil {
			return nil, errors.New("invalid date format for monthly, use YYYY-MM")
		}
		toDate = fromDate.AddDate(0, 1, -1) // Last day of month
	case "yearly":
		// Date format: YYYY
		fromDate, err = time.Parse("2006", date)
		if err != nil {
			return nil, errors.New("invalid date format for yearly, use YYYY")
		}
		toDate = fromDate.AddDate(1, 0, -1) // Last day of year
	}

	// Query transactions for period
	summary := &Summary{
		CategoryBreakdown: make(map[string]float64),
	}

	currentDate := fromDate
	for !currentDate.After(toDate) {
		dayResults := cash_flow_mapper.INSTANCE.GetCashFlowsByBelongsDate(currentDate)

		for _, cashFlow := range dayResults {
			summary.TransactionCount++

			if cashFlow.FlowType == model.FlowTypeIncome {
				summary.TotalIncome += cashFlow.Amount
			} else {
				summary.TotalExpense += cashFlow.Amount
			}

			// Get category name for breakdown
			category := category_mapper.INSTANCE.GetCategoryByObjectId(cashFlow.CategoryId.Hex())
			if !category.IsEmpty() {
				summary.CategoryBreakdown[category.Name] += cashFlow.Amount
			}
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	summary.Balance = summary.TotalIncome - summary.TotalExpense

	return summary, nil
}
