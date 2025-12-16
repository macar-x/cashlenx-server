package model

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/macar-x/cashlenx-server/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryEntity struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	ParentId   primitive.ObjectID `json:"parent_id" bson:"parent_id"`
	Name       string             `json:"name" bson:"name"`
	Type       string             `json:"type" bson:"type"`
	Remark     string             `json:"remark" bson:"remark"`
	CreateTime time.Time          `json:"create_time" bson:"create_time"`
	ModifyTime time.Time          `json:"modify_time" bson:"modify_time"`
}

func (entity CategoryEntity) IsEmpty() bool {
	return reflect.DeepEqual(entity, CategoryEntity{})
}

func (entity CategoryEntity) ToString() string {
	return "[ " +
		"Id: " + entity.Id.Hex() +
		", Name: " + entity.Name +
		", Type: " + entity.Type +
		" ]"
}

// MarshalJSON customizes JSON marshaling for CategoryEntity
// Converts timestamps to the configured timezone for display
func (entity CategoryEntity) MarshalJSON() ([]byte, error) {
	// Convert timestamps to configured timezone for display
	localCreateTime := util.ToTimezone(entity.CreateTime)
	localModifyTime := util.ToTimezone(entity.ModifyTime)

	// Create a temporary struct with local timezone timestamps
	type Alias CategoryEntity
	return json.Marshal(&struct {
		CreateTime time.Time `json:"create_time"`
		ModifyTime time.Time `json:"modify_time"`
		*Alias
	}{
		CreateTime: localCreateTime,
		ModifyTime: localModifyTime,
		Alias:      (*Alias)(&entity),
	})
}
