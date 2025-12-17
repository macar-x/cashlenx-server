package category_mapper

import (
	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/util"
)

var INSTANCE CategoryMapper

type CategoryMapper interface {
	GetCategoryByObjectId(userId, plainId string) model.CategoryEntity
	GetCategoryByName(userId, categoryName string) model.CategoryEntity
	GetCategoryByParentId(userId, parentPlainId string) []model.CategoryEntity
	InsertCategoryByEntity(newEntity model.CategoryEntity) string
	UpdateCategoryByEntity(userId, plainId string, updatedEntity model.CategoryEntity) model.CategoryEntity
	GetAllCategories(userId string, limit, offset int) []model.CategoryEntity
	CountAllCategories(userId string) int64
	DeleteCategoryByObjectId(userId, plainId string) model.CategoryEntity
	TruncateCategories() error
}

func init() {
	switch util.GetConfigByKey("db.type") {
	case "mongodb":
		INSTANCE = CategoryMongoDbMapper{}
	case "mysql":
		INSTANCE = CategoryMySqlMapper{}
	default:
		panic("database type not supported")
	}
}
