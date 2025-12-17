package controllers

import (
	"encoding/json"
	"fmt"
	"intehub/config"
	"intehub/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

var cronScheduler *cron.Cron

// InitScheduler 初始化定时任务调度器
func InitScheduler() {
	cronScheduler = cron.New()
	cronScheduler.Start()

	// 加载已启用的任务
	db := config.GetDB()
	var tasks []models.ScheduleTask
	db.Where("enabled = ?", true).Find(&tasks)

	for _, task := range tasks {
		addTaskToCron(task)
	}
}

// addTaskToCron 添加任务到调度器
func addTaskToCron(task models.ScheduleTask) error {
	_, err := cronScheduler.AddFunc(task.CronExpr, func() {
		executeTask(task)
	})
	return err
}

// executeTask 执行任务
func executeTask(task models.ScheduleTask) {
	startTime := time.Now()
	db := config.GetDB()

	log := models.ScheduleLog{
		TaskID:    task.ID,
		TaskName:  task.Name,
		StartedAt: startTime,
	}

	var err error
	var message string

	// 根据任务类型执行
	switch task.TaskType {
	case "push":
		message, err = executePushTask(task.ConfigID, task.FieldData)
	default:
		message = "未知的任务类型"
		err = nil
	}

	finishTime := time.Now()
	log.FinishedAt = finishTime
	log.Duration = finishTime.Sub(startTime).Milliseconds()

	if err != nil {
		log.Status = "failed"
		log.Message = err.Error()
	} else {
		log.Status = "success"
		log.Message = message
	}

	db.Create(&log)

	// 更新任务的最后运行时间
	now := time.Now()
	db.Model(&task).Updates(map[string]interface{}{
		"last_run_at": now,
	})
}

// executePushTask 执行推送任务
func executePushTask(configID uint, fieldDataJSON string) (string, error) {
	db := config.GetDB()
	var pushConfig models.PushConfig

	if err := db.First(&pushConfig, configID).Error; err != nil {
		return "", err
	}

	// 解析字段数据
	var fieldData map[string]interface{}
	if fieldDataJSON != "" {
		if err := json.Unmarshal([]byte(fieldDataJSON), &fieldData); err != nil {
			return "", err
		}
	} else {
		fieldData = make(map[string]interface{})
	}

	// 转换为 map[string]string 供 executePush 使用
	dataStr := make(map[string]string)
	for k, v := range fieldData {
		dataStr[k] = fmt.Sprintf("%v", v)
	}

	// 调用实际的推送函数
	history := ExecutePushInternal(&pushConfig, dataStr)

	// 保存历史记录（系统用户ID为0）
	history.UserID = 0
	db.Create(&history)

	if !history.Success {
		return "", fmt.Errorf("集成失败: %s", history.Error)
	}

	return fmt.Sprintf("集成成功，状态码: %d", history.StatusCode), nil
}

// GetScheduleTasks 获取定时任务列表
func GetScheduleTasks(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	db := config.GetDB()
	var tasks []models.ScheduleTask

	if role == "admin" {
		db.Order("created_at desc").Find(&tasks)
	} else {
		db.Where("user_id = ?", userID).Order("created_at desc").Find(&tasks)
	}

	c.JSON(http.StatusOK, tasks)
}

// GetScheduleTask 获取单个定时任务
func GetScheduleTask(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	db := config.GetDB()
	var task models.ScheduleTask

	query := db.Where("id = ?", id)
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// CreateScheduleTask 创建定时任务
func CreateScheduleTask(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var task models.ScheduleTask
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	task.ID = 0
	task.UserID = userID.(uint)

	// 验证 Cron 表达式
	if _, err := cron.ParseStandard(task.CronExpr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cron表达式格式错误"})
		return
	}

	db := config.GetDB()
	if err := db.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	// 如果启用，添加到调度器
	if task.Enabled {
		addTaskToCron(task)
	}

	c.JSON(http.StatusCreated, task)
}

// UpdateScheduleTask 更新定时任务
func UpdateScheduleTask(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	db := config.GetDB()
	var task models.ScheduleTask

	query := db.Where("id = ?", id)
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	var updateData models.ScheduleTask
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 验证 Cron 表达式
	if updateData.CronExpr != "" {
		if _, err := cron.ParseStandard(updateData.CronExpr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cron表达式格式错误"})
			return
		}
	}

	if err := db.Model(&task).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	// 重新加载调度器（简单实现：停止并重启）
	cronScheduler.Stop()
	InitScheduler()

	c.JSON(http.StatusOK, task)
}

// DeleteScheduleTask 删除定时任务
func DeleteScheduleTask(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	db := config.GetDB()
	var task models.ScheduleTask

	query := db.Where("id = ?", id)
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	if err := db.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	// 重新加载调度器
	cronScheduler.Stop()
	InitScheduler()

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetScheduleLogs 获取定时任务执行日志
func GetScheduleLogs(c *gin.Context) {
	taskID := c.Query("task_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	db := config.GetDB()
	var logs []models.ScheduleLog
	var total int64

	query := db.Model(&models.ScheduleLog{})
	if taskID != "" {
		query = query.Where("task_id = ?", taskID)
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

// ToggleScheduleTask 启用/禁用定时任务
func ToggleScheduleTask(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	db := config.GetDB()
	var task models.ScheduleTask

	query := db.Where("id = ?", id)
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	task.Enabled = !task.Enabled
	if err := db.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}

	// 重新加载调度器
	cronScheduler.Stop()
	InitScheduler()

	c.JSON(http.StatusOK, task)
}
