package manage_service

import (
	"time"

	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/service/cash_flow_service"
	"github.com/macar-x/cashlenx-server/service/category_service"
	"github.com/macar-x/cashlenx-server/util"
)

// InitializeDemoData initializes the database with demo categories and transactions
func InitializeDemoData() error {
	// Create default categories
	categories := []string{
		"Salary",
		"Freelance",
		"Food & Dining",
		"Transportation",
		"Shopping",
		"Entertainment",
		"Healthcare",
		"Utilities",
	}

	for _, catName := range categories {
		// Check if category already exists
		_, err := category_service.CreateService("", catName)
		if err != nil {
			util.Logger.Warnw("category creation skipped", "category", catName, "error", err)
		}
	}

	// Create sample transactions for the past week
	today := time.Now()

	// Sample income
	_, _ = cash_flow_service.SaveIncome(
		today.AddDate(0, 0, -7).Format(model.DateFormatYYYYMMDD),
		"Salary",
		5000.00,
		"Monthly salary",
	)

	// Sample expenses
	expenses := []struct {
		daysAgo     int
		category    string
		amount      float64
		description string
	}{
		{1, "Food & Dining", 45.50, "Lunch"},
		{1, "Transportation", 20.00, "Bus fare"},
		{2, "Food & Dining", 32.00, "Groceries"},
		{3, "Entertainment", 50.00, "Movie tickets"},
		{4, "Shopping", 120.00, "Clothes"},
		{5, "Healthcare", 80.00, "Pharmacy"},
		{6, "Utilities", 150.00, "Electricity bill"},
	}

	for _, exp := range expenses {
		date := today.AddDate(0, 0, -exp.daysAgo).Format(model.DateFormatYYYYMMDD)
		_, _ = cash_flow_service.SaveOutcome(date, exp.category, exp.amount, exp.description)
	}

	return nil
}
