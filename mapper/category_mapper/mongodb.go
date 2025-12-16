package category_mapper

import (
	"context"
	"time"

	"github.com/macar-x/cashlenx-server/cache"
	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/util"
	"github.com/macar-x/cashlenx-server/util/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CategoryMongoDbMapper struct{}

func (CategoryMongoDbMapper) GetCategoryByObjectId(plainId string) model.CategoryEntity {
	objectId := util.Convert2ObjectId(plainId)
	if plainId == "" || objectId == primitive.NilObjectID {
		util.Logger.Warnln("category's id is not acceptable")
		return model.CategoryEntity{}
	}

	filter := bson.D{
		primitive.E{Key: "_id", Value: objectId},
	}

	database.OpenMongoDbConnection(database.CategoryTableName)
	defer database.CloseMongoDbConnection()
	return convertBsonM2CategoryEntity(database.GetOneInMongoDB(filter))
}

func (CategoryMongoDbMapper) GetCategoryByName(categoryName string) model.CategoryEntity {
	// Check cache first
	categoryCache := cache.GetCategoryCache()
	if cached, ok := categoryCache.GetByName(categoryName); ok {
		return *cached
	}

	// Cache miss - query database
	filter := bson.D{
		primitive.E{Key: "name", Value: categoryName},
	}

	database.OpenMongoDbConnection(database.CategoryTableName)
	defer database.CloseMongoDbConnection()
	entity := convertBsonM2CategoryEntity(database.GetOneInMongoDB(filter))

	// Store in cache if found
	if !entity.IsEmpty() {
		categoryCache.Set(&entity)
	}

	return entity
}

func (CategoryMongoDbMapper) GetCategoryByParentId(parentPlainId string) []model.CategoryEntity {
	// Convert parentPlainId to ObjectID
	parentObjectId := util.Convert2ObjectId(parentPlainId)
	filter := bson.D{
		primitive.E{Key: "parent_id", Value: parentObjectId},
	}

	database.OpenMongoDbConnection(database.CategoryTableName)
	defer database.CloseMongoDbConnection()

	var targetEntityList []model.CategoryEntity
	queryResultList := database.GetManyInMongoDB(filter)
	for _, queryResult := range queryResultList {
		targetEntityList = append(targetEntityList, convertBsonM2CategoryEntity(queryResult))
	}

	return targetEntityList
}

func (CategoryMongoDbMapper) InsertCategoryByEntity(newEntity model.CategoryEntity) string {
	operatingTime := time.Now()
	newEntity.CreateTime = operatingTime
	newEntity.ModifyTime = operatingTime

	database.OpenMongoDbConnection(database.CategoryTableName)
	defer database.CloseMongoDbConnection()

	newCategoryId := database.InsertOneInMongoDB(convertCategoryEntity2BsonD(newEntity))

	// Invalidate cache on insert
	cache.GetCategoryCache().Clear()

	return newCategoryId.Hex()
}

func (CategoryMongoDbMapper) UpdateCategoryByEntity(plainId string, updatedEntity model.CategoryEntity) model.CategoryEntity {
	objectId := util.Convert2ObjectId(plainId)
	if plainId == "" || objectId == primitive.NilObjectID {
		util.Logger.Warnln("category's id is not acceptable")
		return model.CategoryEntity{}
	}

	filter := bson.D{
		primitive.E{Key: "_id", Value: objectId},
	}

	database.OpenMongoDbConnection(database.CategoryTableName)
	defer database.CloseMongoDbConnection()

	targetEntity := convertBsonM2CategoryEntity(database.GetOneInMongoDB(filter))
	if targetEntity.IsEmpty() {
		util.Logger.Infoln("category is not exist")
		return model.CategoryEntity{}
	}

	// Update fields from updatedEntity while preserving ID and CreateTime
	updatedEntity.Id = targetEntity.Id
	updatedEntity.CreateTime = targetEntity.CreateTime
	updatedEntity.ModifyTime = time.Now()

	rowsAffected := database.UpdateManyInMongoDB(filter, convertCategoryEntity2BsonD(updatedEntity))
	if rowsAffected != 1 {
		// fixme: maybe we should have a rollback here.
		util.Logger.Errorw("update failed", "rows_affected", rowsAffected)
		return model.CategoryEntity{}
	}

	// Invalidate cache on update
	cache.GetCategoryCache().Clear()

	return updatedEntity
}

func (CategoryMongoDbMapper) DeleteCategoryByObjectId(plainId string) model.CategoryEntity {
	objectId := util.Convert2ObjectId(plainId)
	if plainId == "" || objectId == primitive.NilObjectID {
		util.Logger.Warnln("category's id is not acceptable")
		return model.CategoryEntity{}
	}

	filter := bson.D{
		primitive.E{Key: "_id", Value: objectId},
	}

	database.OpenMongoDbConnection(database.CategoryTableName)
	defer database.CloseMongoDbConnection()

	targetEntity := convertBsonM2CategoryEntity(database.GetOneInMongoDB(filter))
	if targetEntity.IsEmpty() {
		util.Logger.Infoln("category is not exist")
		return model.CategoryEntity{}
	}

	// can not delete a category that has referred child-categories.
	if len(INSTANCE.GetCategoryByParentId(plainId)) != 0 {
		util.Logger.Infoln("can not delete a category which has child-categories refer to")
		return model.CategoryEntity{}
	}

	rowsAffected := database.DeleteManyInMongoDB(filter)
	if rowsAffected != 1 {
		// fixme: maybe we should have a rollback here.
		util.Logger.Errorw("delete failed", "rows_affected", rowsAffected)
		return model.CategoryEntity{}
	}

	// Invalidate cache on delete
	cache.GetCategoryCache().Clear()

	return targetEntity
}

func (CategoryMongoDbMapper) GetAllCategories(limit, offset int) []model.CategoryEntity {
	database.OpenMongoDbConnection(database.CategoryTableName)
	defer database.CloseMongoDbConnection()

	collection := database.GetMongoCollection(database.CategoryTableName)

	// Empty filter to get all documents, with pagination
	filter := bson.D{}

	ctx := context.TODO()
	findOptions := options.Find()
	if limit > 0 {
		findOptions.SetLimit(int64(limit))
	}
	if offset > 0 {
		findOptions.SetSkip(int64(offset))
	}
	// Sort by name ascending
	findOptions.SetSort(bson.D{primitive.E{Key: "name", Value: 1}})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		util.Logger.Errorw("query all categories failed", "error", err)
		return []model.CategoryEntity{}
	}
	defer cursor.Close(ctx)

	var targetEntityList []model.CategoryEntity
	for cursor.Next(ctx) {
		var bsonM bson.M
		if err := cursor.Decode(&bsonM); err != nil {
			util.Logger.Errorw("decode failed", "error", err)
			continue
		}
		targetEntityList = append(targetEntityList, convertBsonM2CategoryEntity(bsonM))
	}

	return targetEntityList
}

func (CategoryMongoDbMapper) CountAllCategories() int64 {
	filter := bson.D{}

	database.OpenMongoDbConnection(database.CategoryTableName)
	defer database.CloseMongoDbConnection()

	return database.CountInMongoDB(filter)
}

func convertCategoryEntity2BsonD(entity model.CategoryEntity) bson.D {
	// Generate a new Id automatically if it's empty
	if entity.Id == primitive.NilObjectID {
		entity.Id = primitive.NewObjectID()
	}

	return bson.D{
		primitive.E{Key: "_id", Value: entity.Id},
		primitive.E{Key: "parent_id", Value: entity.ParentId},
		primitive.E{Key: "name", Value: entity.Name},
		primitive.E{Key: "remark", Value: entity.Remark},
		primitive.E{Key: "create_time", Value: entity.CreateTime},
		primitive.E{Key: "modify_time", Value: entity.ModifyTime},
	}
}

func convertBsonM2CategoryEntity(bsonM bson.M) model.CategoryEntity {
	var newEntity model.CategoryEntity
	bsonBytes, _ := bson.Marshal(bsonM)
	err := bson.Unmarshal(bsonBytes, &newEntity)
	if err != nil {
		panic(err)
	}
	return newEntity
}
