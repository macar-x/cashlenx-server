package cash_flow_service

import (
	"github.com/macar-x/cashlenx-server/mapper/cash_flow_mapper"
	"github.com/macar-x/cashlenx-server/model"
)

// QueryAll queries all cash flows with optional filtering and pagination
func QueryAll(cashType string, limit, offset int) ([]*model.CashFlowEntity, int64, error) {
	// Get total count
	totalCount := cash_flow_mapper.INSTANCE.CountAllCashFlows()

	// Get paginated results
	cashFlows := cash_flow_mapper.INSTANCE.GetAllCashFlows(limit, offset)

	// Filter by cash type if specified
	var filteredResults []*model.CashFlowEntity
	if cashType != "" {
		for i := range cashFlows {
			if cashFlows[i].FlowType == cashType {
				filteredResults = append(filteredResults, &cashFlows[i])
			}
		}
	} else {
		// Convert to pointer slice
		for i := range cashFlows {
			filteredResults = append(filteredResults, &cashFlows[i])
		}
	}

	return filteredResults, totalCount, nil
}
