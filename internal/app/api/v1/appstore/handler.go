package appstore

import (
	"errors"
	"intehub/internal/app/service/appstore"
	httputil "intehub/internal/utils/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *appstore.Service
}

func New(service *appstore.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandleAppStoreAPI(rg *gin.RouterGroup) {
	j := httputil.NewJSONHandler(rg)

	// 应用模板管理
	j.GET("/templates", h.ListTemplates)
	j.GET("/templates/:id", h.GetTemplate)
	j.POST("/templates", h.CreateTemplate)
	j.PUT("/templates/:id", h.UpdateTemplate)
	j.DELETE("/templates/:id", h.DeleteTemplate)

	// 应用配置管理
	j.GET("/configs/:appId", h.GetAppConfigs)
	j.POST("/configs", h.CreateConfig)
	j.PUT("/configs/:id", h.UpdateConfig)
	j.DELETE("/configs/:id", h.DeleteConfig)

	// 从模板安装应用
	j.POST("/install", h.InstallFromTemplate)
}

// ListTemplates 列出应用模板
func (h *Handler) ListTemplates(c *gin.Context) (interface{}, error) {
	category := c.Query("category")

	templates, err := h.service.ListTemplates(category)
	if err != nil {
		return nil, err
	}

	return templates, nil
}

// GetTemplate 获取应用模板详情
func (h *Handler) GetTemplate(c *gin.Context) (interface{}, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	template, err := h.service.GetTemplate(uint(id))
	if err != nil {
		return nil, err
	}

	return template, nil
}

// CreateTemplate 创建应用模板
func (h *Handler) CreateTemplate(c *gin.Context) (interface{}, error) {
	var req appstore.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("参数错误")
	}

	template, err := h.service.CreateTemplate(&req)
	if err != nil {
		return nil, err
	}

	return template, nil
}

// UpdateTemplate 更新应用模板
func (h *Handler) UpdateTemplate(c *gin.Context) (interface{}, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	var req appstore.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("参数错误")
	}

	template, err := h.service.UpdateTemplate(uint(id), &req)
	if err != nil {
		return nil, err
	}

	return template, nil
}

// DeleteTemplate 删除应用模板
func (h *Handler) DeleteTemplate(c *gin.Context) (interface{}, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	if err := h.service.DeleteTemplate(uint(id)); err != nil {
		return nil, err
	}

	return gin.H{"message": "template deleted"}, nil
}

// GetAppConfigs 获取应用配置
func (h *Handler) GetAppConfigs(c *gin.Context) (interface{}, error) {
	appID, err := strconv.ParseUint(c.Param("appId"), 10, 32)
	if err != nil {
		return nil, errors.New("invalid app id")
	}

	configs, err := h.service.GetAppConfigs(uint(appID))
	if err != nil {
		return nil, err
	}

	return configs, nil
}

// CreateConfig 创建应用配置
func (h *Handler) CreateConfig(c *gin.Context) (interface{}, error) {
	var req appstore.CreateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("参数错误")
	}

	config, err := h.service.CreateConfig(&req)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// UpdateConfig 更新应用配置
func (h *Handler) UpdateConfig(c *gin.Context) (interface{}, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	var req appstore.UpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("参数错误")
	}

	config, err := h.service.UpdateConfig(uint(id), &req)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// DeleteConfig 删除应用配置
func (h *Handler) DeleteConfig(c *gin.Context) (interface{}, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	if err := h.service.DeleteConfig(uint(id)); err != nil {
		return nil, err
	}

	return gin.H{"message": "config deleted"}, nil
}

// InstallFromTemplate 从模板安装应用
func (h *Handler) InstallFromTemplate(c *gin.Context) (interface{}, error) {
	var req appstore.InstallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("参数错误")
	}

	app, err := h.service.InstallFromTemplate(&req)
	if err != nil {
		return nil, err
	}

	return app, nil
}
