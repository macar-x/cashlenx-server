package statistic_service

import (
	"fmt"
)

// GetBreakdownForUser gets category breakdown for the specified period
// Only includes transactions belonging to the specified user
func GetBreakdownForUser(period, date, userId string) (*Breakdown, error) {
	// TODO: Implement user-specific breakdown
	// 1. Parse period and date
	// 2. Get all cash flows for user in period
	// 3. Group by category and type
	// 4. Calculate percentages
	// 5. Sort by amount descending
	return nil, fmt.Errorf("GetBreakdownForUser not yet implemented - coming soon")
}
