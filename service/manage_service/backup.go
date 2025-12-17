package manage_service

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/macar-x/cashlenx-server/mapper/cash_flow_mapper"
	"github.com/macar-x/cashlenx-server/mapper/category_mapper"
)

// BackupData represents the structure of backup data
type BackupData struct {
	Version    string                   `json:"version"`
	Timestamp  string                   `json:"timestamp"`
	CashFlows  []map[string]interface{} `json:"cash_flows"`
	Categories []map[string]interface{} `json:"categories"`
}

// EntityStats represents statistics for a single entity type (cash_flows or categories)
type EntityStats struct {
	Success    int      `json:"success"`
	Failed     int      `json:"failed"`
	FailedList []string `json:"failed_list,omitempty"`
}

// OperationStats represents comprehensive statistics for an operation (backup/restore/truncate)
type OperationStats struct {
	CashFlows  EntityStats `json:"cash_flows"`
	Categories EntityStats `json:"categories"`
}

// CreateBackup creates a backup of all database data
func CreateBackup(filePath string) (OperationStats, error) {
	stats := OperationStats{
		CashFlows:  EntityStats{Success: 0, Failed: 0, FailedList: []string{}},
		Categories: EntityStats{Success: 0, Failed: 0, FailedList: []string{}},
	}

	if filePath == "" {
		return stats, errors.New("file path cannot be empty")
	}

	// Get all categories (no pagination limit - get everything)
	categories := category_mapper.INSTANCE.GetAllCategories(0, 0)
	stats.Categories.Success = len(categories)

	// Convert categories to map format for JSON serialization
	categoryMaps := make([]map[string]interface{}, len(categories))
	for i, cat := range categories {
		categoryMaps[i] = map[string]interface{}{
			"Id":         cat.Id.Hex(),
			"Name":       cat.Name,
			"ParentId":   cat.ParentId.Hex(),
			"Remark":     cat.Remark,
			"CreateTime": cat.CreateTime,
			"ModifyTime": cat.ModifyTime,
		}
	}

	// Get all cash flows (no pagination limit - get everything)
	cashFlows := cash_flow_mapper.INSTANCE.GetAllCashFlows(0, 0)
	stats.CashFlows.Success = len(cashFlows)

	// Convert cash flows to map format for JSON serialization
	cashFlowMaps := make([]map[string]interface{}, len(cashFlows))
	for i, cf := range cashFlows {
		cashFlowMaps[i] = map[string]interface{}{
			"Id":          cf.Id.Hex(),
			"CategoryId":  cf.CategoryId.Hex(),
			"BelongsDate": cf.BelongsDate,
			"FlowType":    cf.FlowType,
			"Amount":      cf.Amount,
			"Description": cf.Description,
			"Remark":      cf.Remark,
			"CreateTime":  cf.CreateTime,
			"ModifyTime":  cf.ModifyTime,
		}
	}

	// Create backup structure
	backup := BackupData{
		Version:    "1.0.0",
		Timestamp:  time.Now().Format(time.RFC3339),
		CashFlows:  cashFlowMaps,
		Categories: categoryMaps,
	}

	// Write to file
	file, err := os.Create(filePath)
	if err != nil {
		return stats, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(backup); err != nil {
		return stats, err
	}

	return stats, nil
}
