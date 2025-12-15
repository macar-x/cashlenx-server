package model

import (
	"reflect"
	"strconv"
	"time"

	"github.com/macar-x/cashlenx-server/util"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CashFlowEntity struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	CategoryId  primitive.ObjectID `json:"category_id" bson:"category_id"`
	BelongsDate time.Time          `json:"belongs_date" bson:"belongs_date"`
	FlowType    string             `json:"flow_type" bson:"flow_type"`
	Amount      float64            `json:"amount" bson:"amount"`
	Description string             `json:"description" bson:"description"`
	Remark      string             `json:"remark" bson:"remark"`
	CreateTime  time.Time          `json:"create_time" bson:"create_time"`
	ModifyTime  time.Time          `json:"modify_time" bson:"modify_time"`
}

func (entity CashFlowEntity) IsEmpty() bool {
	return reflect.DeepEqual(entity, CashFlowEntity{})
}

func (entity CashFlowEntity) ToString() string {
	// todo: category query from cache like redis would be better.
	return "[ " +
		"Id: " + entity.Id.Hex() +
		// fixme: using mapper to fetch category name?
		// ", Category: " + mapper.GetCategoryMapper().GetCategoryByObjectId(entity.CategoryId.Hex()).Name +
		", Date: " + util.FormatDateToStringWithoutDash(entity.BelongsDate) +
		", FlowType: " + entity.FlowType +
		", Amount: " + strconv.FormatFloat(entity.Amount, 'f', 2, 64) +
		", Description: " + entity.Description +
		" ]"
}

func (entity CashFlowEntity) Build(fieldMap map[string]string) CashFlowEntity {
	newEntity := entity
	for key, value := range fieldMap {
		switch key {
		case "Id":
			objectId, err := primitive.ObjectIDFromHex(value)
			if err != nil {
				util.Logger.Warnln("build cash failed with err: " + err.Error())
			}
			newEntity.Id = objectId
		case "CategoryId":
			newEntity.CategoryId = util.Convert2ObjectId(value)
		case "BelongsDate":
			newEntity.BelongsDate = util.FormatDateFromStringWithoutDash(value)
		case "FlowType":
			// todo: use enum to check if value available
			newEntity.FlowType = value
		case "Amount":
			amount, err := strconv.ParseFloat(value, 64)
			if err != nil {
				util.Logger.Warnln("build cash failed with err: " + err.Error())
			}
			newEntity.Amount = amount
		case "Description":
			newEntity.Description = value
		case "Remark":
			newEntity.Remark = value
		}
	}
	return newEntity
}
