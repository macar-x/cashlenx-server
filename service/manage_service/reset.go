package manage_service

import (
	"errors"
)

// ResetDatabase clears all data from the database
func ResetDatabase() error {
	// Note: To properly implement this, we need:
	// 1. Mapper method to delete all cash flows: DeleteAllCashFlows()
	// 2. Mapper method to delete all categories: DeleteAllCategories()

	// This is a dangerous operation and should be implemented carefully
	// with proper transaction support and confirmation

	// TODO: Add to mapper interfaces:
	//   - DeleteAllCashFlows() int64
	//   - DeleteAllCategories() int64

	return errors.New("database reset requires mapper enhancement - need DeleteAll methods")
}
