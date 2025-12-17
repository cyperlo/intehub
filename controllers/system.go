package controllers

import (
	"intehub/config"
	"intehub/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetSystemLogs 获取系统日志列表
func GetSystemLogs(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	level := c.Query("level")
	module := c.Query("module")

	db := config.GetDB()
	var logs []models.SystemLog
	var total int64

	query := db.Model(&models.SystemLog{})
	if level != "" {
		query = query.Where("level = ?", level)
	}
	if module != "" {
		query = query.Where("module = ?", module)
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

// CreateSystemLog 创建系统日志
func CreateSystemLog(c *gin.Context) {
	var log models.SystemLog
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	log.IP = c.ClientIP()
	if userID, exists := c.Get("user_id"); exists {
		log.UserID = userID.(uint)
	}

	db := config.GetDB()
	if err := db.Create(&log).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusCreated, log)
}

// GetMenus 获取菜单列表
func GetMenus(c *gin.Context) {
	role, _ := c.Get("role")

	db := config.GetDB()
	var menus []models.Menu

	query := db.Where("visible = ?", true).Order("sort asc, id asc")

	// 非管理员需要过滤角色权限
	if role != "admin" {
		query = query.Where("roles LIKE ?", "%"+role.(string)+"%")
	}

	query.Find(&menus)
	c.JSON(http.StatusOK, menus)
}

// CreateMenu 创建菜单
func CreateMenu(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限操作"})
		return
	}

	var menu models.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	db := config.GetDB()
	if err := db.Create(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusCreated, menu)
}

// UpdateMenu 更新菜单
func UpdateMenu(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限操作"})
		return
	}

	id := c.Param("id")
	db := config.GetDB()
	var menu models.Menu

	if err := db.First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "菜单不存在"})
		return
	}

	var updateData models.Menu
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if err := db.Model(&menu).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, menu)
}

// DeleteMenu 删除菜单
func DeleteMenu(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限操作"})
		return
	}

	id := c.Param("id")
	db := config.GetDB()
	var menu models.Menu

	if err := db.First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "菜单不存在"})
		return
	}

	if err := db.Delete(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetUsers 获取用户列表
func GetUsers(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问"})
		return
	}

	db := config.GetDB()
	var users []models.User

	db.Select("id, username, nickname, role, created_at, updated_at").Find(&users)
	c.JSON(http.StatusOK, users)
}

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限操作"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 检查用户名是否已存在
	db := config.GetDB()
	var count int64
	db.Model(&models.User{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 密码加密
	if err := user.HashPassword(user.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	// 清除密码字段
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限操作"})
		return
	}

	id := c.Param("id")
	db := config.GetDB()
	var user models.User

	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	var updateData struct {
		Nickname string `json:"nickname"`
		Role     string `json:"role"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 更新字段
	if updateData.Nickname != "" {
		user.Nickname = updateData.Nickname
	}
	if updateData.Role != "" {
		user.Role = updateData.Role
	}
	if updateData.Password != "" {
		if err := user.HashPassword(updateData.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
			return
		}
	}

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限操作"})
		return
	}

	id := c.Param("id")
	currentUserID, _ := c.Get("user_id")

	// 不能删除自己
	if strconv.Itoa(int(currentUserID.(uint))) == id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能删除当前登录用户"})
		return
	}

	db := config.GetDB()
	var user models.User

	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
