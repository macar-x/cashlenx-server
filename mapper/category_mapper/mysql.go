package category_mapper

import (
	"bytes"
	"database/sql"
	"time"

	"github.com/macar-x/cashlenx-server/cache"
	"github.com/macar-x/cashlenx-server/model"
	"github.com/macar-x/cashlenx-server/util"
	"github.com/macar-x/cashlenx-server/util/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryMySqlMapper struct{}

func (CategoryMySqlMapper) GetCategoryByObjectId(plainId string) model.CategoryEntity {
	var sqlString bytes.Buffer
	sqlString.WriteString("SELECT ID, PARENT_ID, NAME, TYPE FROM ")
	sqlString.WriteString(database.CategoryTableName)
	sqlString.WriteString(" WHERE ID = ? ")

	connection := database.GetMySqlConnection()
	defer database.CloseMySqlConnection()

	rows, err := connection.Query(sqlString.String(), plainId)
	if err != nil {
		util.Logger.Errorw("query failed", "error", err)
	}

	var categoryEntity model.CategoryEntity
	for rows.Next() {
		categoryEntity = convertRow2CategoryEntity(rows)
		break
	}
	return categoryEntity
}

func (CategoryMySqlMapper) GetCategoryByName(categoryName string) model.CategoryEntity {
	// Check cache first
	categoryCache := cache.GetCategoryCache()
	if cached, ok := categoryCache.GetByName(categoryName); ok {
		return *cached
	}

	// Cache miss - query database
	var sqlString bytes.Buffer
	sqlString.WriteString("SELECT ID, PARENT_ID, NAME, TYPE FROM ")
	sqlString.WriteString(database.CategoryTableName)
	sqlString.WriteString(" WHERE NAME = ? ")

	connection := database.GetMySqlConnection()
	defer database.CloseMySqlConnection()

	rows, err := connection.Query(sqlString.String(), categoryName)
	if err != nil {
		util.Logger.Errorw("query failed", "error", err)
	}

	var categoryEntity model.CategoryEntity
	for rows.Next() {
		categoryEntity = convertRow2CategoryEntity(rows)
		break
	}

	// Store in cache if found
	if !categoryEntity.IsEmpty() {
		categoryCache.Set(&categoryEntity)
	}

	return categoryEntity
}

func (CategoryMySqlMapper) GetCategoryByParentId(parentPlainId string) []model.CategoryEntity {
	var sqlString bytes.Buffer
	sqlString.WriteString("SELECT ID, PARENT_ID, NAME, TYPE FROM ")
	sqlString.WriteString(database.CategoryTableName)
	sqlString.WriteString(" WHERE PARENT_ID = ? ")

	connection := database.GetMySqlConnection()
	defer database.CloseMySqlConnection()

	rows, err := connection.Query(sqlString.String(), parentPlainId)
	if err != nil {
		util.Logger.Errorw("query failed", "error", err)
	}

	var targetEntityList []model.CategoryEntity
	for rows.Next() {
		targetEntityList = append(targetEntityList, convertRow2CategoryEntity(rows))
	}
	return targetEntityList
}

func (CategoryMySqlMapper) InsertCategoryByEntity(newEntity model.CategoryEntity) string {
	operatingTime := time.Now().UTC() // Store in UTC
	newEntity.CreateTime = operatingTime
	newEntity.ModifyTime = operatingTime

	var sqlString bytes.Buffer
	sqlString.WriteString("INSERT ")
	sqlString.WriteString(database.CategoryTableName)
	sqlString.WriteString(" SET ID = ?, ")
	sqlString.WriteString(" PARENT_ID = ?, ")
	sqlString.WriteString(" NAME = ?, ")
	sqlString.WriteString(" TYPE = ?, ")
	sqlString.WriteString(" REMARK = ?, ")
	sqlString.WriteString(" CREATE_TIME = ?, ")
	sqlString.WriteString(" MODIFY_TIME = ? ")

	connection := database.GetMySqlConnection()
	defer database.CloseMySqlConnection()

	statement, err := connection.Prepare(sqlString.String())
	if err != nil {
		util.Logger.Errorw("insert failed", "error", err)
	}

	newPlainId := primitive.NewObjectID().Hex()
	result, err := statement.Exec(newPlainId, newEntity.ParentId.Hex(), newEntity.Name,
		newEntity.Type, newEntity.Remark, operatingTime, operatingTime)
	if err != nil {
		util.Logger.Errorw("insert failed", "error", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected != 1 {
		// fixme: maybe we should have a rollback here.
		util.Logger.Errorw("insert failed", "error", err, "rows_affected", rowsAffected)
	}

	// Invalidate cache on insert
	cache.GetCategoryCache().Clear()

	return newPlainId
}

func (CategoryMySqlMapper) UpdateCategoryByEntity(plainId string, updatedEntity model.CategoryEntity) model.CategoryEntity {
	targetEntity := INSTANCE.GetCategoryByObjectId(plainId)
	if targetEntity.IsEmpty() {
		util.Logger.Infoln("category is not exist")
		return model.CategoryEntity{}
	}

	// Update fields from updatedEntity while preserving ID and CreateTime
	updatedEntity.Id = targetEntity.Id
	updatedEntity.CreateTime = targetEntity.CreateTime
	updatedEntity.ModifyTime = time.Now().UTC() // Store in UTC

	var sqlString bytes.Buffer
	sqlString.WriteString("UPDATE ")
	sqlString.WriteString(database.CategoryTableName)
	sqlString.WriteString(" SET PARENT_ID = ?, ")
	sqlString.WriteString(" NAME = ?, ")
	sqlString.WriteString(" TYPE = ?, ")
	sqlString.WriteString(" REMARK = ?, ")
	sqlString.WriteString(" MODIFY_TIME = ? ")
	sqlString.WriteString(" WHERE ID = ? ")

	connection := database.GetMySqlConnection()
	defer database.CloseMySqlConnection()

	statement, err := connection.Prepare(sqlString.String())
	if err != nil {
		util.Logger.Errorw("update failed", "error", err)
	}

	result, err := statement.Exec(updatedEntity.ParentId.Hex(), updatedEntity.Name, updatedEntity.Type, updatedEntity.Remark,
		updatedEntity.ModifyTime, updatedEntity.Id)
	if err != nil {
		util.Logger.Errorw("update failed", "error", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected != 1 {
		// fixme: maybe we should have a rollback here.
		util.Logger.Errorw("update failed", "error", err, "rows_affected", rowsAffected)
	}

	// Invalidate cache on update
	cache.GetCategoryCache().Clear()

	return updatedEntity
}

func (CategoryMySqlMapper) DeleteCategoryByObjectId(plainId string) model.CategoryEntity {
	targetEntity := INSTANCE.GetCategoryByObjectId(plainId)
	if targetEntity.IsEmpty() {
		util.Logger.Infoln("category is not exist")
		return model.CategoryEntity{}
	}

	// can not delete a category that has referred child-categories.
	if len(INSTANCE.GetCategoryByParentId(plainId)) != 0 {
		util.Logger.Infoln("can not delete a category which has child-categories refer to")
		return model.CategoryEntity{}
	}

	var sqlString bytes.Buffer
	sqlString.WriteString("DELETE FROM ")
	sqlString.WriteString(database.CategoryTableName)
	sqlString.WriteString(" WHERE ID = ? ")

	connection := database.GetMySqlConnection()
	defer database.CloseMySqlConnection()

	statement, err := connection.Prepare(sqlString.String())
	if err != nil {
		util.Logger.Errorw("delete failed", "error", err)
	}

	result, err := statement.Exec(plainId)
	if err != nil {
		util.Logger.Errorw("delete failed", "error", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected != 1 {
		// fixme: maybe we should have a rollback here.
		util.Logger.Errorw("delete failed", "error", err, "rows_affected", rowsAffected)
	}

	// Invalidate cache on delete
	cache.GetCategoryCache().Clear()

	return targetEntity
}

func (CategoryMySqlMapper) GetAllCategories(limit, offset int) []model.CategoryEntity {
	var sqlString bytes.Buffer
	sqlString.WriteString("SELECT ID, PARENT_ID, NAME FROM ")
	sqlString.WriteString(database.CategoryTableName)
	sqlString.WriteString(" ORDER BY NAME ASC ")

	if limit > 0 {
		sqlString.WriteString(" LIMIT ? OFFSET ? ")
	}

	connection := database.GetMySqlConnection()
	defer database.CloseMySqlConnection()

	var rows *sql.Rows
	var err error

	if limit > 0 {
		rows, err = connection.Query(sqlString.String(), limit, offset)
	} else {
		rows, err = connection.Query(sqlString.String())
	}

	if err != nil {
		util.Logger.Errorw("query all categories failed", "error", err)
		return []model.CategoryEntity{}
	}

	var targetEntityList []model.CategoryEntity
	for rows.Next() {
		targetEntityList = append(targetEntityList, convertRow2CategoryEntity(rows))
	}
	return targetEntityList
}

func (CategoryMySqlMapper) CountAllCategories() int64 {
	var sqlString bytes.Buffer
	sqlString.WriteString("SELECT COUNT(1) FROM ")
	sqlString.WriteString(database.CategoryTableName)

	connection := database.GetMySqlConnection()
	defer database.CloseMySqlConnection()

	rows, err := connection.Query(sqlString.String())
	if err != nil {
		util.Logger.Errorw("count all categories failed", "error", err)
		return 0
	}

	var count int64
	rows.Next()
	if err = rows.Scan(&count); err != nil {
		util.Logger.Errorw("parse count failed", "error", err)
		return 0
	}
	return count
}

func (CategoryMySqlMapper) TruncateCategories() error {
	var sqlString bytes.Buffer
	sqlString.WriteString("TRUNCATE TABLE ")
	sqlString.WriteString(database.CategoryTableName)

	connection := database.GetMySqlConnection()
	defer database.CloseMySqlConnection()

	_, err := connection.Exec(sqlString.String())
	if err != nil {
		util.Logger.Errorw("truncate categories failed", "error", err)
		return err
	}

	// Clear cache after truncate
	cache.GetCategoryCache().Clear()

	util.Logger.Infow("Categories truncated successfully")
	return nil
}

func convertRow2CategoryEntity(rows *sql.Rows) model.CategoryEntity {
	var id string
	var parentId string
	var name string
	var categoryType string

	err := rows.Scan(&id, &parentId, &name, &categoryType)
	if err != nil {
		util.Logger.Errorw("covert into entity failed", "error", err)
	}

	return model.CategoryEntity{
		Id:       util.Convert2ObjectId(id),
		ParentId: util.Convert2ObjectId(parentId),
		Name:     name,
		Type:     categoryType,
	}
}
