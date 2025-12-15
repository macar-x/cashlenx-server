package category_service

import (
	"errors"
	"fmt"

	"github.com/macar-x/cashlenx-server/mapper/cash_flow_mapper"
	"github.com/macar-x/cashlenx-server/mapper/category_mapper"
	"github.com/macar-x/cashlenx-server/validation"
)

func DeleteService(plainId, categoryName string) error {
	if isDeleteFieldsConflicted(plainId, categoryName) {
		return errors.New("should have one and only one delete type")
	}

	if plainId != "" {
		return deleteById(plainId)
	}

	if categoryName != "" {
		return deleteByName(categoryName)
	}

	return errors.New("not supported delete type")
}

func isDeleteFieldsConflicted(plainId, categoryName string) bool {
	// check if already one semi-optional field is filled
	semiOptionalFieldFilledFlag := false

	// plain_id is not empty
	if plainId != "" {
		semiOptionalFieldFilledFlag = true
	}

	// category_name is not empty
	if categoryName != "" {
		if semiOptionalFieldFilledFlag {
			return true
		}
		semiOptionalFieldFilledFlag = true
	}

	// should have one and only one field filled
	return !semiOptionalFieldFilledFlag
}

func deleteById(plainId string) error {
	// Validate ID
	if err := validation.ValidateID(plainId); err != nil {
		return err
	}

	existCategoryEntity := category_mapper.INSTANCE.GetCategoryByObjectId(plainId)
	if existCategoryEntity.IsEmpty() {
		fmt.Println("category not found")
		return nil
	}

	if cash_flow_mapper.INSTANCE.CountCashFLowsByCategoryId(existCategoryEntity.Id.Hex()) != 0 {
		return errors.New("can not delete a category which has cash_flows refer to")
	}

	existCategoryEntity = category_mapper.INSTANCE.DeleteCategoryByObjectId(plainId)
	if existCategoryEntity.IsEmpty() {
		return errors.New("category delete failed")
	}
	fmt.Println("category ", 0, ": ", existCategoryEntity.ToString())
	return nil
}

func deleteByName(categoryName string) error {
	// Validate category name
	if err := validation.ValidateCategoryName(categoryName); err != nil {
		return err
	}

	existCategoryEntity := category_mapper.INSTANCE.GetCategoryByName(categoryName)
	if existCategoryEntity.IsEmpty() {
		fmt.Println("category not found")
		return nil
	}

	if cash_flow_mapper.INSTANCE.CountCashFLowsByCategoryId(existCategoryEntity.Id.Hex()) != 0 {
		return errors.New("can not delete a category which has cash_flows refer to")
	}

	existCategoryEntity = category_mapper.INSTANCE.DeleteCategoryByObjectId(existCategoryEntity.Id.Hex())
	if existCategoryEntity.IsEmpty() {
		return errors.New("category delete failed")
	}
	fmt.Println("category ", 0, ": ", existCategoryEntity.ToString())
	return nil
}
