package statistic_service

import (
	"fmt"
)

// ExportForUser exports the user's cash flow data to Excel
// Only exports data belonging to the specified user
func ExportForUser(fromDate, toDate, filePath, userId string) error {
	// TODO: Implement user-specific export
	// 1. Get all categories for user
	// 2. Get all cash flows for user in date range
	// 3. Create Excel file with user's data only
	return fmt.Errorf("ExportForUser not yet implemented - coming soon")
}
