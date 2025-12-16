package cash_flow_service

import (
	"errors"
	"reflect"
	"time"

	"github.com/macar-x/cashlenx-server/mapper/cash_flow_mapper"
	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/util"
)

func IsQueryFieldsConflicted(plainId, belongsDate, exactDescription, fuzzyDescription string) bool {
	// check if already one semi-optional field is filled
	semiOptionalFieldFilledFlag := false

	// plain_id is not empty
	if plainId != "" {
		semiOptionalFieldFilledFlag = true
	}

	// belongs_date is not empty
	if belongsDate != "" {
		if semiOptionalFieldFilledFlag {
			return true
		}
		semiOptionalFieldFilledFlag = true
	}

	// exact_description is not empty
	if exactDescription != "" {
		if semiOptionalFieldFilledFlag {
			return true
		}
		semiOptionalFieldFilledFlag = true
	}

	// fuzzy_description is not empty
	if fuzzyDescription != "" {
		if semiOptionalFieldFilledFlag {
			return true
		}
		semiOptionalFieldFilledFlag = true
	}

	// should have one and only one field filled
	return !semiOptionalFieldFilledFlag
}

func QueryById(plainId string) (model.CashFlowEntity, error) {
	cashFlowEntity := cash_flow_mapper.INSTANCE.GetCashFlowByObjectId(plainId)
	if cashFlowEntity.IsEmpty() {
		return model.CashFlowEntity{}, errors.New("cash_flow not found")
	}
	return cashFlowEntity, nil
}

func QueryByDate(belongsDate string) ([]model.CashFlowEntity, error) {
	queryDate := util.FormatDateFromStringWithoutDash(belongsDate)
	if reflect.DeepEqual(queryDate, time.Time{}) {
		return []model.CashFlowEntity{}, errors.New("belongs_date error, try format like 19700101")
	}

	// Use date range query instead of exact match to handle dates with time components
	// Set start to beginning of the day and end to end of the day
	startOfDay := time.Date(queryDate.Year(), queryDate.Month(), queryDate.Day(), 0, 0, 0, 0, queryDate.Location())
	endOfDay := time.Date(queryDate.Year(), queryDate.Month(), queryDate.Day(), 23, 59, 59, 999999999, queryDate.Location())

	matchedCashFlowList := cash_flow_mapper.INSTANCE.GetCashFlowsByDateRange(startOfDay, endOfDay)
	return matchedCashFlowList, nil
}

func QueryByExactDescription(exactDescription string) ([]model.CashFlowEntity, error) {
	matchedCashFlowList := cash_flow_mapper.INSTANCE.GetCashFlowsByExactDesc(exactDescription)
	return matchedCashFlowList, nil
}

func QueryByFuzzyDescription(fuzzyDescription string) ([]model.CashFlowEntity, error) {
	matchedCashFlowList := cash_flow_mapper.INSTANCE.GetCashFlowsByFuzzyDesc(fuzzyDescription)
	return matchedCashFlowList, nil
}
