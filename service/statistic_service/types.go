package statistic_service

// Summary represents a financial summary for a period
type Summary struct {
	Period             string             `json:"period"`
	PeriodType         string             `json:"period_type"`
	Income             float64            `json:"income"`
	Expense            float64            `json:"expense"`
	Balance            float64            `json:"balance"`
	TransactionCount   int                `json:"transaction_count"`
	IncomeCount        int                `json:"income_count"`
	ExpenseCount       int                `json:"expense_count"`
	AverageTransaction float64            `json:"average_transaction"`
	Categories         map[string]float64 `json:"categories"`
}

// CategoryBreakdownItem represents a single category in the breakdown
type CategoryBreakdownItem struct {
	Category   string  `json:"category"`
	Amount     float64 `json:"amount"`
	Percentage float64 `json:"percentage"`
	Count      int     `json:"count"`
}

// Breakdown represents category breakdown analysis
type Breakdown struct {
	Period            string                  `json:"period"`
	TotalExpense      float64                 `json:"total_expense"`
	TotalIncome       float64                 `json:"total_income"`
	ExpenseCategories []CategoryBreakdownItem `json:"expense_categories"`
	IncomeCategories  []CategoryBreakdownItem `json:"income_categories"`
}

// TrendDataPoint represents a single data point in the trend
type TrendDataPoint struct {
	Date    string  `json:"date"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
}

// TrendAnalysis represents the trend analysis results
type TrendAnalysis struct {
	IncomeTrend           string  `json:"income_trend"`
	ExpenseTrend          string  `json:"expense_trend"`
	AverageMonthlyExpense float64 `json:"average_monthly_expense"`
}

// Trends represents spending trends over time
type Trends struct {
	Period     string           `json:"period"`
	PeriodType string           `json:"period_type"`
	DataPoints []TrendDataPoint `json:"data_points"`
	Trends     TrendAnalysis    `json:"trends"`
}

// TopExpense represents a single top expense
type TopExpense struct {
	ID          string  `json:"id"`
	Date        string  `json:"date"`
	Category    string  `json:"category"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Percentage  float64 `json:"percentage"`
}

// TopExpenses represents the top N expenses
type TopExpenses struct {
	Period       string       `json:"period"`
	Limit        int          `json:"limit"`
	TotalExpense float64      `json:"total_expense"`
	Expenses     []TopExpense `json:"expenses"`
}
