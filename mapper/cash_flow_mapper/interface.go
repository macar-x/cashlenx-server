package cash_flow_mapper

import (
	"time"

	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/util"
)

var INSTANCE CashFlowMapper

type CashFlowMapper interface {
	GetCashFlowByObjectId(userId, plainId string) model.CashFlowEntity
	GetCashFlowsByObjectIdArray(userId string, plainIdList []string) []model.CashFlowEntity
	GetCashFlowsByBelongsDate(userId string, belongsDate time.Time) []model.CashFlowEntity
	GetCashFlowsByDateRange(userId string, from, to time.Time) []model.CashFlowEntity
	GetCashFlowsByCategoryId(userId, categoryPlainId string) []model.CashFlowEntity
	GetCashFlowsByExactDesc(userId, description string) []model.CashFlowEntity
	GetCashFlowsByFuzzyDesc(userId, description string) []model.CashFlowEntity
	CountCashFLowsByCategoryId(userId, categoryPlainId string) int64
	InsertCashFlowByEntity(newEntity model.CashFlowEntity) string
	BulkInsertCashFlows(entities []model.CashFlowEntity) ([]string, error)
	UpdateCashFlowByEntity(userId, plainId string, updatedEntity model.CashFlowEntity) model.CashFlowEntity
	GetAllCashFlows(userId string, limit, offset int) []model.CashFlowEntity
	CountAllCashFlows(userId string) int64
	DeleteCashFlowByObjectId(userId, plainId string) model.CashFlowEntity
	DeleteCashFlowByBelongsDate(userId string, belongsDate time.Time) []model.CashFlowEntity
	TruncateCashFlows() error
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
