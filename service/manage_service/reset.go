package manage_service

import (
	"github.com/macar-x/cashlenx-server/mapper/cash_flow_mapper"
	"github.com/macar-x/cashlenx-server/mapper/category_mapper"
)

// ResetDatabase clears all data from the database
func ResetDatabase() error {
	// This is a dangerous operation - truncate all data
	// First truncate cash flows (dependent data)
	if err := cash_flow_mapper.INSTANCE.TruncateCashFlows(); err != nil {
		return err
	}

	// Then truncate categories (parent data)
	if err := category_mapper.INSTANCE.TruncateCategories(); err != nil {
		return err
	}

	return nil
}

// TruncateDatabase is an alias for ResetDatabase - clears all data from the database
func TruncateDatabase() error {
	return ResetDatabase()
}
