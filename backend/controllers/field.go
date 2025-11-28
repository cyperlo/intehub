package controllers

import (
	"intehub/config"
	"intehub/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetFieldSchemas 获取字段定义列表
func GetFieldSchemas(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	db := config.GetDB()
	var fields []models.FieldSchema

	// 普通用户只能看到自己的字段，管理员可以看到所有
	if role == "admin" {
		db.Order("created_at desc").Find(&fields)
	} else {
		db.Where("user_id = ?", userID).Order("created_at desc").Find(&fields)
	}

	c.JSON(http.StatusOK, fields)
}

// GetFieldSchema 获取单个字段定义
func GetFieldSchema(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	db := config.GetDB()
	var field models.FieldSchema

	query := db.Where("id = ?", id)
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&field).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "字段不存在"})
		return
	}

	c.JSON(http.StatusOK, field)
}

// CreateFieldSchema 创建字段定义
func CreateFieldSchema(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var field models.FieldSchema
	if err := c.ShouldBindJSON(&field); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	field.UserID = userID.(uint)

	db := config.GetDB()
	if err := db.Create(&field).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusCreated, field)
}

// UpdateFieldSchema 更新字段定义
func UpdateFieldSchema(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	db := config.GetDB()
	var existingField models.FieldSchema

	query := db.Where("id = ?", id)
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&existingField).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "字段不存在"})
		return
	}

	var updateData models.FieldSchema
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 保留原有的用户ID
	updateData.UserID = existingField.UserID

	if err := db.Model(&existingField).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, existingField)
}

// DeleteFieldSchema 删除字段定义
func DeleteFieldSchema(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	db := config.GetDB()
	var field models.FieldSchema

	query := db.Where("id = ?", id)
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&field).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "字段不存在"})
		return
	}

	if err := db.Delete(&field).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetConfigFields 获取配置关联的字段
func GetConfigFields(c *gin.Context) {
	configID := c.Param("id")

	db := config.GetDB()
	var relations []models.ConfigFieldRelation

	if err := db.Where("config_id = ?", configID).Order("\"order\" asc").Find(&relations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	// 获取字段详情
	var fieldIDs []uint
	for _, rel := range relations {
		fieldIDs = append(fieldIDs, rel.FieldID)
	}

	var fields []models.FieldSchema
	if len(fieldIDs) > 0 {
		db.Where("id IN ?", fieldIDs).Find(&fields)
	}

	// 构建返回数据，保持顺序
	fieldMap := make(map[uint]models.FieldSchema)
	for _, field := range fields {
		fieldMap[field.ID] = field
	}

	result := make([]models.FieldSchema, 0, len(relations))
	for _, rel := range relations {
		if field, ok := fieldMap[rel.FieldID]; ok {
			result = append(result, field)
		}
	}

	c.JSON(http.StatusOK, result)
}

// UpdateConfigFields 更新配置关联的字段
func UpdateConfigFields(c *gin.Context) {
	configID := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	// 验证配置权限
	db := config.GetDB()
	var pushConfig models.PushConfig
	query := db.Where("id = ?", configID)
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&pushConfig).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "配置不存在"})
		return
	}

	// 接收字段ID列表
	var req struct {
		FieldIDs []uint `json:"field_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 删除旧的关联
	db.Where("config_id = ?", configID).Delete(&models.ConfigFieldRelation{})

	// 创建新的关联
	for i, fieldID := range req.FieldIDs {
		relation := models.ConfigFieldRelation{
			ConfigID: pushConfig.ID,
			FieldID:  fieldID,
			Order:    i,
		}
		if err := db.Create(&relation).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}
