package cash_flow_service

import (
	"strings"
	"time"

	"github.com/macar-x/cashlenx-server/mapper/cash_flow_mapper"
	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/util"
)

// QueryAll queries all cash flows with optional filtering and pagination
func QueryAll(
	cashType string,
	categoryId string,
	description string,
	exactDescription string,
	fromDateStr string,
	toDateStr string,
	limit int,
	offset int,
) ([]*model.CashFlowEntity, int64, error) {
	// Get total count (we'll filter this later if needed)
	totalCount := cash_flow_mapper.INSTANCE.CountAllCashFlows()

	// Get paginated results
	cashFlows := cash_flow_mapper.INSTANCE.GetAllCashFlows(limit, offset)

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
