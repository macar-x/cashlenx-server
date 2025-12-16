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

// UpdateById updates a cash flow record by ID
func UpdateById(plainId, belongsDate, categoryName string, amount float64, description string) (model.CashFlowEntity, error) {
	// Validate ID
	if err := validation.ValidateID(plainId); err != nil {
		return model.CashFlowEntity{}, err
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

	// Query existing record
	existingEntity := cash_flow_mapper.INSTANCE.GetCashFlowByObjectId(plainId)
	if existingEntity.IsEmpty() {
		return model.CashFlowEntity{}, errors.New("cash_flow not found")
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
		categoryEntity := category_mapper.INSTANCE.GetCategoryByName(categoryName)
		if categoryEntity.IsEmpty() {
			return model.CashFlowEntity{}, errors.New("category does not exist")
		}
		existingEntity.CategoryId = categoryEntity.Id
	}

	if amount != 0 {
		// Round to 2 decimal places
		amount, _ = decimal.NewFromFloat(amount).Round(2).Float64()
		existingEntity.Amount = amount
	}

	if description != "" {
		existingEntity.Description = description
	}

	// Update modify time
	existingEntity.ModifyTime = time.Now().UTC() // Store in UTC

	// Call mapper to update the record
	updatedEntity := cash_flow_mapper.INSTANCE.UpdateCashFlowByEntity(plainId, existingEntity)
	if updatedEntity.IsEmpty() {
		return model.CashFlowEntity{}, errors.New("failed to update cash_flow")
	}

	return updatedEntity, nil
}
