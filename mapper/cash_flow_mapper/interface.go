package cash_flow_mapper

import (
	"time"

	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/util"
)

var INSTANCE CashFlowMapper

type CashFlowMapper interface {
	GetCashFlowByObjectId(plainId string) model.CashFlowEntity
	GetCashFlowsByObjectIdArray(plainIdList []string) []model.CashFlowEntity
	GetCashFlowsByBelongsDate(belongsDate time.Time) []model.CashFlowEntity
	GetCashFlowsByDateRange(from, to time.Time) []model.CashFlowEntity
	GetCashFlowsByCategoryId(categoryPlainId string) []model.CashFlowEntity
	GetCashFlowsByExactDesc(description string) []model.CashFlowEntity
	GetCashFlowsByFuzzyDesc(description string) []model.CashFlowEntity
	CountCashFLowsByCategoryId(categoryPlainId string) int64
	InsertCashFlowByEntity(newEntity model.CashFlowEntity) string
	BulkInsertCashFlows(entities []model.CashFlowEntity) ([]string, error)
	UpdateCashFlowByEntity(plainId string, updatedEntity model.CashFlowEntity) model.CashFlowEntity
	GetAllCashFlows(limit, offset int) []model.CashFlowEntity
	CountAllCashFlows() int64
	DeleteCashFlowByObjectId(plainId string) model.CashFlowEntity
	DeleteCashFlowByBelongsDate(belongsDate time.Time) []model.CashFlowEntity
}

func init() {
	switch util.GetConfigByKey("db.type") {
	case "mongodb":
		INSTANCE = CashFlowMongoDbMapper{}
	case "mysql":
		INSTANCE = CashFlowMySqlMapper{}
	default:
		panic("database type not supported")
	}
}
