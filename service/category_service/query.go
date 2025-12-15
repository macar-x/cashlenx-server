package category_service

import (
	"errors"
	"fmt"

	"github.com/macar-x/cashlenx-server/mapper/category_mapper"
	"github.com/macar-x/cashlenx-server/model"
)

func QueryService(plainId, parentPlainId, categoryName string) ([]model.CategoryEntity, error) {
	if isQueryFieldsConflicted(plainId, parentPlainId, categoryName) {
		return nil, errors.New("should have one and only one query type")
	}

	if plainId != "" {
		return queryById(plainId), nil
	}

	if parentPlainId != "" {
		return queryByParentId(parentPlainId), nil
	}

	if categoryName != "" {
		return queryByName(categoryName), nil
	}

	return nil, errors.New("not supported query type")
}

func isQueryFieldsConflicted(plainId, parentPlainId, name string) bool {
	// check if already one semi-optional field is filled
	semiOptionalFieldFilledFlag := false

	// plain_id is not empty
	if plainId != "" {
		semiOptionalFieldFilledFlag = true
	}

	// parent_plain_id is not empty
	if parentPlainId != "" {
		if semiOptionalFieldFilledFlag {
			return true
		}
		semiOptionalFieldFilledFlag = true
	}

	// category name is not empty
	if name != "" {
		if semiOptionalFieldFilledFlag {
			return true
		}
		semiOptionalFieldFilledFlag = true
	}

	// should have one and only one field filled
	return !semiOptionalFieldFilledFlag
}

func queryById(plainId string) []model.CategoryEntity {
	categoryEntity := category_mapper.INSTANCE.GetCategoryByObjectId(plainId)
	if categoryEntity.IsEmpty() {
		fmt.Println("category not found")
		return []model.CategoryEntity{}
	}
	fmt.Println("category ", 0, ": ", categoryEntity.ToString())
	return []model.CategoryEntity{categoryEntity}
}

func queryByParentId(plainParentId string) []model.CategoryEntity {
	matchedCategoryList := category_mapper.INSTANCE.GetCategoryByParentId(plainParentId)
	if len(matchedCategoryList) == 0 {
		fmt.Println("no matched categories")
		return []model.CategoryEntity{}
	}

	for index, categoryEntity := range matchedCategoryList {
		fmt.Println("category ", index, ": ", categoryEntity.ToString())
	}
	return matchedCategoryList
}

func queryByName(categoryName string) []model.CategoryEntity {
	categoryEntity := category_mapper.INSTANCE.GetCategoryByName(categoryName)
	if categoryEntity.IsEmpty() {
		fmt.Println("category not found")
		return []model.CategoryEntity{}
	}
	fmt.Println("category ", 0, ": ", categoryEntity.ToString())
	return []model.CategoryEntity{categoryEntity}
}
