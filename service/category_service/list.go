package category_service

import (
	"github.com/macar-x/cashlenx-server/mapper/category_mapper"
	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ListAllService lists all categories with pagination, filtered by user ID and type
func ListAllService(userId, categoryType string, limit, offset int) ([]model.CategoryEntity, int64, error) {
	// Validate user ID
	userObjectId := util.Convert2ObjectId(userId)
	if userObjectId == primitive.NilObjectID {
		return nil, 0, util.NewError(util.ErrInvalidUserId, "invalid user ID")
	}

	// Get total count with filters
	totalCount, err := category_mapper.INSTANCE.CountCategoriesByUserAndType(userObjectId, categoryType)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results with filters
	categories, err := category_mapper.INSTANCE.GetCategoriesByUserAndType(userObjectId, categoryType, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return categories, totalCount, nil
}