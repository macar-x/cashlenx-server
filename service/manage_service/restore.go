package manage_service

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/macar-x/cashlenx-server/mapper/cash_flow_mapper"
	"github.com/macar-x/cashlenx-server/mapper/category_mapper"
	"github.com/macar-x/cashlenx-server/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RestoreBackup restores database from a backup file
func RestoreBackup(filePath string) (OperationStats, error) {
	stats := OperationStats{
		CashFlows:  EntityStats{Success: 0, Failed: 0, FailedList: []string{}},
		Categories: EntityStats{Success: 0, Failed: 0, FailedList: []string{}},
	}

	if filePath == "" {
		return stats, errors.New("file path cannot be empty")
	}

	// Read backup file
	file, err := os.Open(filePath)
	if err != nil {
		return stats, err
	}
	defer file.Close()

	// Parse JSON
	var backup BackupData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&backup); err != nil {
		return stats, err
	}

	// Update total counts for stats
	totalCategories := len(backup.Categories)
	totalCashFlows := len(backup.CashFlows)
	stats.Categories.Failed = totalCategories
	stats.CashFlows.Failed = totalCashFlows

	// Step 1: Clear existing data
	if _, err := ResetDatabase(); err != nil {
		return stats, err
	}

	// Step 2: Insert categories from backup
	for _, catMap := range backup.Categories {
		// Parse Id from backup data
		id, _ := primitive.ObjectIDFromHex(catMap["Id"].(string))
		
		// Parse ParentId from backup data
		parentId, _ := primitive.ObjectIDFromHex(catMap["ParentId"].(string))
		
		// Parse CreateTime and ModifyTime
		createTime, _ := time.Parse(time.RFC3339, catMap["CreateTime"].(string))
		modifyTime, _ := time.Parse(time.RFC3339, catMap["ModifyTime"].(string))
		
		// Create category entity from backup data, preserving all original fields
		catEntity := model.CategoryEntity{
			Id:         id,
			ParentId:   parentId,
			Name:       catMap["Name"].(string),
			Remark:     catMap["Remark"].(string),
			CreateTime: createTime,
			ModifyTime: modifyTime,
		}

		// Insert category
		if id := category_mapper.INSTANCE.InsertCategoryByEntity(catEntity); id != "" {
			stats.Categories.Success++
			stats.Categories.Failed--
		}
	}

	// Step 3: Insert cash flows from backup
	cashFlowEntities := make([]model.CashFlowEntity, totalCashFlows)
	for i, cfMap := range backup.CashFlows {
		// Parse Id from backup data
		id, _ := primitive.ObjectIDFromHex(cfMap["Id"].(string))
		
		// Parse belongs_date string to time.Time
		belongsDate, _ := time.Parse(time.RFC3339, cfMap["BelongsDate"].(string))

		// Parse CategoryId from backup data
		categoryId, _ := primitive.ObjectIDFromHex(cfMap["CategoryId"].(string))
		
		// Parse CreateTime and ModifyTime
		createTime, _ := time.Parse(time.RFC3339, cfMap["CreateTime"].(string))
		modifyTime, _ := time.Parse(time.RFC3339, cfMap["ModifyTime"].(string))
		
		// Create cash flow entity from backup data, preserving all original fields
		cfEntity := model.CashFlowEntity{
			Id:          id,
			CategoryId:  categoryId,
			BelongsDate: belongsDate,
			FlowType:    cfMap["FlowType"].(string),
			Amount:      cfMap["Amount"].(float64),
			Description: cfMap["Description"].(string),
			Remark:      cfMap["Remark"].(string),
			CreateTime:  createTime,
			ModifyTime:  modifyTime,
		}
		cashFlowEntities[i] = cfEntity
	}

	// Use bulk insert for cash flows if available
	if len(cashFlowEntities) > 0 {
		if ids, err := cash_flow_mapper.INSTANCE.BulkInsertCashFlows(cashFlowEntities); err != nil {
			// If bulk insert fails, try individual inserts
			for _, cfEntity := range cashFlowEntities {
				if id := cash_flow_mapper.INSTANCE.InsertCashFlowByEntity(cfEntity); id != "" {
					stats.CashFlows.Success++
					stats.CashFlows.Failed--
				}
			}
		} else {
			// Bulk insert succeeded
			stats.CashFlows.Success = len(ids)
			stats.CashFlows.Failed = totalCashFlows - len(ids)
		}
	}

	return stats, nil
}
