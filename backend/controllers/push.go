package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"intehub/config"
	"intehub/models"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetPushConfigs 获取推送配置列表
func GetPushConfigs(c *gin.Context) {
	userID, _ := c.Get("user_id")

	db := config.GetDB()
	var configs []models.PushConfig

	// 普通用户只能看到自己的配置，管理员可以看到所有
	role, _ := c.Get("role")
	if role == "admin" {
		db.Order("created_at desc").Find(&configs)
	} else {
		db.Where("user_id = ?", userID).Order("created_at desc").Find(&configs)
	}

	c.JSON(http.StatusOK, configs)
}

// GetPushConfig 获取单个推送配置
func GetPushConfig(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	db := config.GetDB()
	var pushConfig models.PushConfig

	query := db.Where("id = ?", id)
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&pushConfig).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "配置不存在"})
		return
	}

	c.JSON(http.StatusOK, pushConfig)
}

// CreatePushConfig 创建推送配置
func CreatePushConfig(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var pushConfig models.PushConfig
	if err := c.ShouldBindJSON(&pushConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	pushConfig.UserID = userID.(uint)

	db := config.GetDB()
	if err := db.Create(&pushConfig).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusCreated, pushConfig)
}

// UpdatePushConfig 更新推送配置
func UpdatePushConfig(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	db := config.GetDB()
	var existingConfig models.PushConfig

	query := db.Where("id = ?", id)
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&existingConfig).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "配置不存在"})
		return
	}

	var updateData models.PushConfig
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 保留原有的用户ID
	updateData.UserID = existingConfig.UserID

	if err := db.Model(&existingConfig).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, existingConfig)
}

// DeletePushConfig 删除推送配置
func DeletePushConfig(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	db := config.GetDB()
	var pushConfig models.PushConfig

	query := db.Where("id = ?", id)
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&pushConfig).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "配置不存在"})
		return
	}

	if err := db.Delete(&pushConfig).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// SendPush 执行推送
func SendPush(c *gin.Context) {
	var req struct {
		ConfigID uint              `json:"config_id" binding:"required"`
		Data     map[string]string `json:"data"` // 动态数据，用于替换模板变量
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	db := config.GetDB()
	var pushConfig models.PushConfig

	query := db.Where("id = ?", req.ConfigID)
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&pushConfig).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "配置不存在"})
		return
	}

	if !pushConfig.Enabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该配置已禁用"})
		return
	}

	// 执行推送
	history := executePush(&pushConfig, req.Data)
	history.UserID = userID.(uint)

	// 保存历史记录
	db.Create(&history)

	if history.Success {
		c.JSON(http.StatusOK, gin.H{
			"message": "推送成功",
			"history": history,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "推送失败",
			"history": history,
		})
	}
}

// executePush 执行实际的推送操作
func executePush(config *models.PushConfig, data map[string]string) models.PushHistory {
	startTime := time.Now()
	history := models.PushHistory{
		ConfigID:   config.ID,
		ConfigName: config.Name,
		URL:        config.URL,
		Method:     config.Method,
	}

	// 替换模板中的变量
	content := []byte(config.Template)
	for key, value := range data {
		placeholder := fmt.Sprintf("{{%s}}", key)
		content = bytes.ReplaceAll(content, []byte(placeholder), []byte(value))
	}
	history.Content = string(content)

	// 创建HTTP请求
	var req *http.Request
	var err error

	if config.Method == "GET" {
		req, err = http.NewRequest(config.Method, config.URL, nil)
	} else {
		req, err = http.NewRequest(config.Method, config.URL, bytes.NewBuffer(content))
	}

	if err != nil {
		history.Success = false
		history.Error = fmt.Sprintf("创建请求失败: %v", err)
		history.Duration = time.Since(startTime).Milliseconds()
		return history
	}

	// 设置请求头
	req.Header.Set("Content-Type", config.ContentType)
	if config.Headers != "" {
		var headers map[string]string
		if err := json.Unmarshal([]byte(config.Headers), &headers); err == nil {
			for key, value := range headers {
				req.Header.Set(key, value)
			}
		}
	}

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		history.Success = false
		history.Error = fmt.Sprintf("请求失败: %v", err)
		history.Duration = time.Since(startTime).Milliseconds()
		return history
	}
	defer resp.Body.Close()

	// 读取响应
	body, _ := io.ReadAll(resp.Body)
	history.StatusCode = resp.StatusCode
	history.Response = string(body)
	history.Success = resp.StatusCode >= 200 && resp.StatusCode < 300
	history.Duration = time.Since(startTime).Milliseconds()

	if !history.Success {
		history.Error = fmt.Sprintf("HTTP状态码: %d", resp.StatusCode)
	}

	return history
}

// GetPushHistory 获取推送历史
func GetPushHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	configID := c.Query("config_id")

	db := config.GetDB()
	query := db.Model(&models.PushHistory{})

	// 权限过滤
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	// 配置过滤
	if configID != "" {
		query = query.Where("config_id = ?", configID)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var history []models.PushHistory
	offset := (page - 1) * pageSize
	query.Order("created_at desc").Limit(pageSize).Offset(offset).Find(&history)

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"data":      history,
	})
}
