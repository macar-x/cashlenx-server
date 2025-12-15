package category_service

import (
	"errors"
	"fmt"

	"github.com/macar-x/cashlenx-server/mapper/category_mapper"
)

func QueryService(plainId, parentPlainId, categoryName string) error {
	if isQueryFieldsConflicted(plainId, parentPlainId, categoryName) {
		return errors.New("should have one and only one query type")
	}

	if plainId != "" {
		queryById(plainId)
		return nil
	}

	if parentPlainId != "" {
		queryByParentId(parentPlainId)
		return nil
	}

	if categoryName != "" {
		queryByName(categoryName)
		return nil
	}

	return errors.New("not supported query type")
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

func queryById(plainId string) {
	categoryEntity := category_mapper.INSTANCE.GetCategoryByObjectId(plainId)
	if categoryEntity.IsEmpty() {
		fmt.Println("category not found")
		return
	}
	fmt.Println("category ", 0, ": ", categoryEntity.ToString())
}

func queryByParentId(plainParentId string) {
	matchedCategoryList := category_mapper.INSTANCE.GetCategoryByParentId(plainParentId)
	if len(matchedCategoryList) == 0 {
		fmt.Println("no matched categories")
		return
	}

	for index, categoryEntity := range matchedCategoryList {
		fmt.Println("category ", index, ": ", categoryEntity.ToString())
	}
}

func queryByName(categoryName string) {
	categoryEntity := category_mapper.INSTANCE.GetCategoryByName(categoryName)
	if categoryEntity.IsEmpty() {
		fmt.Println("category not found")
		return
	}
	fmt.Println("category ", 0, ": ", categoryEntity.ToString())
}
