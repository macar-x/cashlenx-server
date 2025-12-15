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

	// 取小數點後兩位
	amount, _ = decimal.NewFromFloat(amount).Round(2).Float64()

	// 必填參數: 類別
	categoryEntity := category_mapper.INSTANCE.GetCategoryByName(categoryName)
	if categoryEntity.IsEmpty() {
		return model.CashFlowEntity{}, errors.New("category does not exist")
	}

	// 選填參數: 日期（默認當天）
	date := util.FormatDateFromStringWithoutDash(util.FormatDateToStringWithoutDash(time.Now()))
	if belongsDate != "" {
		date = util.FormatDateFromStringWithoutDash(belongsDate)
	}

	newCashFlowId := cash_flow_mapper.INSTANCE.InsertCashFlowByEntity(model.CashFlowEntity{
		CategoryId:  categoryEntity.Id,
		BelongsDate: date,
		FlowType:    "INCOME",
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
