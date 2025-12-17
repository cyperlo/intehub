package controllers

import (
	"context"
	"fmt"
	"intehub/config"
	"intehub/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// GetApps 获取应用列表
func GetApps(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	db := config.GetDB()
	var apps []models.App

	if role == "admin" {
		db.Order("created_at desc").Find(&apps)
	} else {
		db.Where("user_id = ?", userID).Order("created_at desc").Find(&apps)
	}

	c.JSON(http.StatusOK, apps)
}

// GetApp 获取单个应用
func GetApp(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	db := config.GetDB()
	var app models.App

	query := db.Where("id = ?", id)
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "应用不存在"})
		return
	}

	c.JSON(http.StatusOK, app)
}

// CreateApp 创建应用
func CreateApp(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var app models.App
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	app.ID = 0
	app.UserID = userID.(uint)

	db := config.GetDB()
	if err := db.Create(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusCreated, app)
}

// UpdateApp 更新应用
func UpdateApp(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	db := config.GetDB()
	var app models.App

	query := db.Where("id = ?", id)
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "应用不存在"})
		return
	}

	var updateData models.App
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if err := db.Model(&app).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, app)
}

// DeleteApp 删除应用
func DeleteApp(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	db := config.GetDB()
	var app models.App

	query := db.Where("id = ?", id)
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "应用不存在"})
		return
	}

	if err := db.Delete(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func code2Runnable(code string) {}

// RunApp 运行应用
func RunApp(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	db := config.GetDB()
	var app models.App

	query := db.Where("id = ?", id)
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "应用不存在"})
		return
	}

	if !app.Enabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "应用已禁用"})
		return
	}

	// 执行应用
	log := executeApp(&app)
	db.Create(&log)

	if log.Status == "success" {
		c.JSON(http.StatusOK, gin.H{
			"message": "执行成功",
			"log":     log,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "执行失败",
			"log":     log,
		})
	}
}

// executeApp 执行应用代码
func executeApp(app *models.App) models.AppLog {
	startTime := time.Now()
	log := models.AppLog{
		AppID:     app.ID,
		AppName:   app.Name,
		StartedAt: startTime,
	}

	// 使用 yaegi 解释执行 Go 代码
	i := interp.New(interp.Options{})
	if err := i.Use(stdlib.Symbols); err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("初始化解释器失败: %v", err)
		log.FinishedAt = time.Now()
		log.Duration = time.Since(startTime).Milliseconds()
		return log
	}

	// 设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 在 goroutine 中执行代码
	done := make(chan struct{})
	var output string
	var execErr error

	go func() {
		defer func() {
			if r := recover(); r != nil {
				execErr = fmt.Errorf("panic: %v", r)
			}
			close(done)
		}()

		// 执行代码
		_, execErr = i.Eval(app.Code)
		if execErr == nil {
			output = "执行成功"
		}
	}()

	// 等待执行完成或超时
	select {
	case <-done:
		if execErr != nil {
			log.Status = "failed"
			log.Error = execErr.Error()
		} else {
			log.Status = "success"
			log.Output = output
		}
	case <-ctx.Done():
		log.Status = "failed"
		log.Error = "执行超时（30秒）"
	}

	log.FinishedAt = time.Now()
	log.Duration = time.Since(startTime).Milliseconds()
	return log
}

// GetAppLogs 获取应用执行日志
func GetAppLogs(c *gin.Context) {
	appID := c.Query("app_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	db := config.GetDB()
	var logs []models.AppLog
	var total int64

	query := db.Model(&models.AppLog{})
	if appID != "" {
		query = query.Where("app_id = ?", appID)
	}

	query.Count(&total)
	query.Order("created_at desc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"list":  logs,
		"total": total,
		"page":  page,
	})
}

// RunAppInternal 内部调用运行应用（供定时任务使用）
func RunAppInternal(appID uint) (models.AppLog, error) {
	db := config.GetDB()
	var app models.App

	if err := db.First(&app, appID).Error; err != nil {
		return models.AppLog{}, err
	}

	if !app.Enabled {
		return models.AppLog{}, fmt.Errorf("应用已禁用")
	}

	log := executeApp(&app)
	db.Create(&log)

	if log.Status != "success" {
		return log, fmt.Errorf("%s", log.Error)
	}

	return log, nil
}
