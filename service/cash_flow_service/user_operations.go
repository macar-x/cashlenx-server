package cash_flow_service

import (
	"errors"
	"strings"
	"time"

	"github.com/macar-x/cashlenx-server/mapper/cash_flow_mapper"
	"github.com/macar-x/cashlenx-server/mapper/category_mapper"
	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/util"
	"github.com/macar-x/cashlenx-server/validation"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User-specific operations for data isolation

// QueryByIdForUser retrieves a cash flow by ID, ensuring it belongs to the user
func QueryByIdForUser(plainId string, userId string) (model.CashFlowEntity, error) {
	// Validate ID
	if err := validation.ValidateID(plainId); err != nil {
		return model.CashFlowEntity{}, err
	}

	// Validate and convert userId
	userObjectId := util.Convert2ObjectId(userId)
	if userObjectId == primitive.NilObjectID {
		return model.CashFlowEntity{}, errors.New("invalid user ID")
	}

	cashFlowEntity := cash_flow_mapper.INSTANCE.GetCashFlowByObjectIdAndUser(plainId, userObjectId)
	if cashFlowEntity.IsEmpty() {
		return model.CashFlowEntity{}, errors.New("cash_flow not found or access denied")
	}
	return cashFlowEntity, nil
}

// QueryByDateForUser retrieves cash flows for a specific date for the user
func QueryByDateForUser(belongsDate string, userId string) ([]model.CashFlowEntity, error) {
	// Validate date
	if err := validation.ValidateDate(belongsDate); err != nil {
		return []model.CashFlowEntity{}, err
	}

	// Validate and convert userId
	userObjectId := util.Convert2ObjectId(userId)
	if userObjectId == primitive.NilObjectID {
		return []model.CashFlowEntity{}, errors.New("invalid user ID")
	}

	// Parse the date string
	parsedDate, err := util.ParseDate(belongsDate)
	if err != nil {
		return []model.CashFlowEntity{}, errors.New("belongs_date error, try format like 19700101, 1970-01-01, or 1970/01/01")
	}

	// Use UTC time for consistent querying
	startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 999999999, time.UTC)

	matchedCashFlowList := cash_flow_mapper.INSTANCE.GetCashFlowsByDateRangeAndUser(startOfDay, endOfDay, userObjectId)
	return matchedCashFlowList, nil
}

// QueryByDateRangeForUser retrieves cash flows within a date range for the user
func QueryByDateRangeForUser(fromDateStr, toDateStr string, userId string) ([]model.CashFlowEntity, error) {
	// Validate and convert userId
	userObjectId := util.Convert2ObjectId(userId)
	if userObjectId == primitive.NilObjectID {
		return []model.CashFlowEntity{}, errors.New("invalid user ID")
	}

	// Parse date strings
	fromDate, err := util.ParseDate(fromDateStr)
	if err != nil {
		return []model.CashFlowEntity{}, errors.New("from_date error, try format like 19700101, 1970-01-01, or 1970/01/01")
	}

	toDate, err := util.ParseDate(toDateStr)
	if err != nil {
		return []model.CashFlowEntity{}, errors.New("to_date error, try format like 19700101, 1970-01-01, or 1970/01/01")
	}

	// Use UTC time
	startDate := time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, time.UTC)
	endDate := time.Date(toDate.Year(), toDate.Month(), toDate.Day(), 23, 59, 59, 999999999, time.UTC)

	matchedCashFlowList := cash_flow_mapper.INSTANCE.GetCashFlowsByDateRangeAndUser(startDate, endDate, userObjectId)
	return matchedCashFlowList, nil
}

// QueryAllForUser queries all cash flows for a user with optional filtering and pagination
func QueryAllForUser(
	userId string,
	cashType string,
	categoryId string,
	description string,
	exactDescription string,
	fromDateStr string,
	toDateStr string,
	limit int,
	offset int,
) ([]*model.CashFlowEntity, int64, error) {
	// Validate and convert userId
	userObjectId := util.Convert2ObjectId(userId)
	if userObjectId == primitive.NilObjectID {
		return nil, 0, errors.New("invalid user ID")
	}

	// Get total count for this user
	totalCount := cash_flow_mapper.INSTANCE.CountAllCashFlowsByUser(userObjectId)

	// Get paginated results for this user
	cashFlows := cash_flow_mapper.INSTANCE.GetAllCashFlowsByUser(userObjectId, limit, offset)

	// Parse date filters if provided
	var fromDate, toDate time.Time
	var err error

	if fromDateStr != "" {
		fromDate, err = util.ParseDate(fromDateStr)
		if err != nil {
			return nil, 0, err
		}
	}

	if toDateStr != "" {
		toDate, err = util.ParseDate(toDateStr)
		if err != nil {
			return nil, 0, err
		}
	}

	// Apply filters
	var filteredResults []*model.CashFlowEntity
	for i := range cashFlows {
		entity := cashFlows[i]
		match := true

		// Filter by cash type
		if cashType != "" && entity.FlowType != cashType {
			match = false
		}

		// Filter by category ID
		if categoryId != "" && entity.CategoryId.Hex() != categoryId {
			match = false
		}

		// Filter by exact description
		if exactDescription != "" && entity.Description != exactDescription {
			match = false
		}

		// Filter by fuzzy description
		if description != "" && exactDescription == "" {
			if !strings.Contains(entity.Description, description) {
				match = false
			}
		}

		// Filter by date range
		if !fromDate.IsZero() && entity.BelongsDate.Before(fromDate) {
			match = false
		}
		if !toDate.IsZero() && entity.BelongsDate.After(toDate) {
			match = false
		}

		if match {
			filteredResults = append(filteredResults, &entity)
		}
	}

	return filteredResults, totalCount, nil
}

// DeleteByIdForUser deletes a cash flow by ID, ensuring it belongs to the user
func DeleteByIdForUser(plainId string, userId string) (model.CashFlowEntity, error) {
	// Validate ID
	if err := validation.ValidateID(plainId); err != nil {
		return model.CashFlowEntity{}, err
	}

	// Validate and convert userId
	userObjectId := util.Convert2ObjectId(userId)
	if userObjectId == primitive.NilObjectID {
		return model.CashFlowEntity{}, errors.New("invalid user ID")
	}

	// Check if it exists and belongs to user
	existCashFlowEntity := cash_flow_mapper.INSTANCE.GetCashFlowByObjectIdAndUser(plainId, userObjectId)
	if existCashFlowEntity.IsEmpty() {
		return model.CashFlowEntity{}, errors.New("cash_flow not found or access denied")
	}

	// Delete it
	deletedEntity := cash_flow_mapper.INSTANCE.DeleteCashFlowByObjectIdAndUser(plainId, userObjectId)
	if deletedEntity.IsEmpty() {
		return model.CashFlowEntity{}, errors.New("cash_flow delete failed")
	}
	return deletedEntity, nil
}

// DeleteByDateForUser deletes cash flows for a specific date for the user
func DeleteByDateForUser(belongsDate string, userId string) ([]model.CashFlowEntity, error) {
	// Validate date
	if err := validation.ValidateDate(belongsDate); err != nil {
		return []model.CashFlowEntity{}, err
	}

	// Validate and convert userId
	userObjectId := util.Convert2ObjectId(userId)
	if userObjectId == primitive.NilObjectID {
		return []model.CashFlowEntity{}, errors.New("invalid user ID")
	}

	// Parse date
	deleteDate := util.FormatDateFromStringWithoutDash(belongsDate)
	if deleteDate.IsZero() {
		return []model.CashFlowEntity{}, errors.New("belongs_date error, try format like 19700101")
	}

	cashFlowList := cash_flow_mapper.INSTANCE.DeleteCashFlowsByBelongsDateAndUser(deleteDate, userObjectId)
	return cashFlowList, nil
}

// UpdateByIdForUser updates a cash flow record by ID, ensuring it belongs to the user
func UpdateByIdForUser(plainId, belongsDate, categoryName string, amount float64, description string, userId string) (model.CashFlowEntity, error) {
	// Validate ID
	if err := validation.ValidateID(plainId); err != nil {
		return model.CashFlowEntity{}, err
	}

	// Validate and convert userId
	userObjectId := util.Convert2ObjectId(userId)
	if userObjectId == primitive.NilObjectID {
		return model.CashFlowEntity{}, errors.New("invalid user ID")
	}

	// Validate optional fields if provided
	if belongsDate != "" {
		if err := validation.ValidateDate(belongsDate); err != nil {
			return model.CashFlowEntity{}, err
		}
	}

	if categoryName != "" {
		if err := validation.ValidateCategoryName(categoryName); err != nil {
			return model.CashFlowEntity{}, err
		}
	}

	if amount != 0 {
		if err := validation.ValidateAmount(amount); err != nil {
			return model.CashFlowEntity{}, err
		}
	}

	if description != "" {
		if err := validation.ValidateDescription(description); err != nil {
			return model.CashFlowEntity{}, err
		}
	}

	// Query existing record - ensure it belongs to the user
	existingEntity := cash_flow_mapper.INSTANCE.GetCashFlowByObjectIdAndUser(plainId, userObjectId)
	if existingEntity.IsEmpty() {
		return model.CashFlowEntity{}, errors.New("cash_flow not found or access denied")
	}

	// Update fields that are provided
	if belongsDate != "" {
		date := util.FormatDateFromStringWithoutDash(belongsDate)
		if date.IsZero() {
			return model.CashFlowEntity{}, errors.New("invalid date format")
		}
		existingEntity.BelongsDate = date
	}

	if categoryName != "" {
		// Note: Category lookup should also be user-specific once categories have user isolation
		categoryEntity := category_mapper.INSTANCE.GetCategoryByName(categoryName)
		if categoryEntity.IsEmpty() {
			return model.CashFlowEntity{}, errors.New("category does not exist")
		}
		existingEntity.CategoryId = categoryEntity.Id
	}

	if amount != 0 {
		// Round to 2 decimal places
		roundedAmount, _ := decimal.NewFromFloat(amount).Round(2).Float64()
		existingEntity.Amount = roundedAmount
	}

	if description != "" {
		existingEntity.Description = description
	}

	// Update modify time
	existingEntity.ModifyTime = time.Now().UTC()

	// Call mapper to update the record
	// Note: Using the regular UpdateCashFlowByEntity because we already verified ownership
	updatedEntity := cash_flow_mapper.INSTANCE.UpdateCashFlowByEntity(plainId, existingEntity)
	if updatedEntity.IsEmpty() {
		return model.CashFlowEntity{}, errors.New("failed to update cash_flow")
	}

	return updatedEntity, nil
}

// Summary represents financial summary data
type Summary struct {
	TotalIncome       float64            `json:"total_income"`
	TotalExpense      float64            `json:"total_expense"`
	Balance           float64            `json:"balance"`
	TransactionCount  int                `json:"transaction_count"`
	CategoryBreakdown map[string]float64 `json:"category_breakdown"`
}

// GetSummaryForUser returns financial summary for a given period for a specific user
func GetSummaryForUser(period, date string, userId string) (*Summary, error) {
	validPeriods := map[string]bool{
		"daily":   true,
		"monthly": true,
		"yearly":  true,
	}

	if !validPeriods[period] {
		return nil, errors.New("invalid period: must be daily, monthly, or yearly")
	}

	// Validate and convert userId
	userObjectId := util.Convert2ObjectId(userId)
	if userObjectId == primitive.NilObjectID {
		return nil, errors.New("invalid user ID")
	}

	var fromDate, toDate time.Time
	var err error

	// Parse date based on period
	switch period {
	case "daily":
		// Date format: YYYY-MM-DD or YYYYMMDD
		parsedDate, err := util.ParseDate(date)
		if err != nil {
			return nil, errors.New("invalid date format for daily")
		}
		fromDate = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC)
		toDate = fromDate
	case "monthly":
		// Date format: YYYY-MM or YYYYMM
		var parsedDate time.Time
		if strings.Contains(date, "-") {
			parsedDate, err = time.Parse("2006-01", date)
		} else if len(date) == 6 {
			parsedDate, err = time.Parse("200601", date)
		} else {
			return nil, errors.New("invalid date format for monthly, use YYYY-MM or YYYYMM")
		}
		if err != nil {
			return nil, errors.New("invalid date format for monthly")
		}
		fromDate = time.Date(parsedDate.Year(), parsedDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		toDate = fromDate.AddDate(0, 1, -1) // Last day of month
	case "yearly":
		// Date format: YYYY
		parsedDate, err := time.Parse("2006", date)
		if err != nil {
			return nil, errors.New("invalid date format for yearly, use YYYY")
		}
		fromDate = time.Date(parsedDate.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		toDate = fromDate.AddDate(1, 0, -1) // Last day of year
	}

	// Query transactions for period using user-specific methods
	summary := &Summary{
		CategoryBreakdown: make(map[string]float64),
	}

	// Use date range query for efficiency instead of iterating day by day
	cashFlows := cash_flow_mapper.INSTANCE.GetCashFlowsByDateRangeAndUser(fromDate, toDate, userObjectId)

	for _, cashFlow := range cashFlows {
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

	summary.Balance = summary.TotalIncome - summary.TotalExpense

	return summary, nil
}
