package cash_flow_service

import (
	"errors"
	"time"

	"github.com/macar-x/cashlenx-server/mapper/cash_flow_mapper"
	"github.com/macar-x/cashlenx-server/mapper/category_mapper"
	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/util"
	"github.com/macar-x/cashlenx-server/validation"
	"github.com/shopspring/decimal"
)

// SaveIncome creates a new income cash flow record
// Note: Could be merged with SaveOutcome into a single SaveCashFlow(flowType, ...) function
func SaveIncome(belongsDate, categoryName string, amount float64, description string) (model.CashFlowEntity, error) {
	// Validate inputs
	if err := validation.ValidateCategoryName(categoryName); err != nil {
		return model.CashFlowEntity{}, err
	}

	if err := validation.ValidateAmount(amount); err != nil {
		return model.CashFlowEntity{}, err
	}

	if belongsDate != "" {
		if err := validation.ValidateDate(belongsDate); err != nil {
			return model.CashFlowEntity{}, err
		}
	}

	if err := validation.ValidateDescription(description); err != nil {
		return model.CashFlowEntity{}, err
	}

	// Round to 2 decimal places
	amount, _ = decimal.NewFromFloat(amount).Round(2).Float64()

	// Required parameter: category
	categoryEntity := category_mapper.INSTANCE.GetCategoryByName(categoryName)
	if categoryEntity.IsEmpty() {
		return model.CashFlowEntity{}, errors.New("category does not exist")
	}

	// Optional parameter: date (default to today)
	var date time.Time
	if belongsDate != "" {
		// Parse the provided date using our multi-format parser
		parsedDate, err := util.ParseDate(belongsDate)
		if err != nil {
			return model.CashFlowEntity{}, errors.New("belongs_date error, try format like 19700101, 1970-01-01, or 1970/01/01")
		}
		// Use UTC time for consistent storage
		date = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC)
	} else {
		// Use today's date in UTC
		today := time.Now().UTC()
		date = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	}

	newCashFlowId := cash_flow_mapper.INSTANCE.InsertCashFlowByEntity(model.CashFlowEntity{
		CategoryId:  categoryEntity.Id,
		BelongsDate: date,
		FlowType:    model.FlowTypeIncome,
		Amount:      amount,
		Description: description,
	})
	if newCashFlowId == "" {
		return model.CashFlowEntity{}, errors.New("cash_flow create failed")
	}

	newCashFlow := cash_flow_mapper.INSTANCE.GetCashFlowByObjectId(newCashFlowId)
	return newCashFlow, nil
}

func IsIncomeRequiredFiledSatisfied(categoryName string, amount float64) bool {
	if categoryName == "" {
		return false
	}
	if amount == 0 {
		return false
	}

	return true
}
