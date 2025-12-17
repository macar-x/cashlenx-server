package category_mapper

import (
	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/util"
)

var INSTANCE CategoryMapper

type CategoryMapper interface {
	GetCategoryByObjectId(plainId string) model.CategoryEntity
	GetCategoryByName(categoryName string) model.CategoryEntity
	GetCategoryByParentId(parentPlainId string) []model.CategoryEntity
	InsertCategoryByEntity(newEntity model.CategoryEntity) string
	UpdateCategoryByEntity(plainId string, updatedEntity model.CategoryEntity) model.CategoryEntity
	GetAllCategories(limit, offset int) []model.CategoryEntity
	CountAllCategories() int64
	DeleteCategoryByObjectId(plainId string) model.CategoryEntity
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
