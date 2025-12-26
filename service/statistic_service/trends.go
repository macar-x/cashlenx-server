package statistic_service

import (
	"fmt"
)

// GetTrendsForUser gets spending trends for the specified period
// Only includes transactions belonging to the specified user
func GetTrendsForUser(period, date, userId string) (*Trends, error) {
	// TODO: Implement user-specific trends
	// 1. Parse period and date
	// 2. Get all cash flows for user in period
	// 3. Group by month/week
	// 4. Calculate trends (increasing, decreasing, stable)
	// 5. Calculate averages
	return nil, fmt.Errorf("GetTrendsForUser not yet implemented - coming soon")
}
