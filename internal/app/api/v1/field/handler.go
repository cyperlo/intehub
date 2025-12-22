package field

import (
	"errors"
	fieldModel "intehub/internal/app/models/field"
	fieldService "intehub/internal/app/service/field"
	httputil "intehub/internal/utils/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	fieldService fieldService.Service
}

func NewHandler(fieldService fieldService.Service) *Handler {
	return &Handler{fieldService: fieldService}
}

func (h *Handler) List(c *gin.Context) (interface{}, error) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var uid uint
	if role != "admin" && userID != nil {
		uid = userID.(uint)
	}

	fields, err := h.fieldService.GetAll(uid)
	if err != nil {
		return nil, err
	}
	return fields, nil
}

func (h *Handler) Get(c *gin.Context) (interface{}, error) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	field, err := h.fieldService.GetByID(uint(id))
	if err != nil {
		return nil, errors.New("字段不存在")
	}
	return field, nil
}

func (h *Handler) Create(c *gin.Context) (interface{}, error) {
	var field fieldModel.FieldSchema
	if err := c.ShouldBindJSON(&field); err != nil {
		return nil, errors.New("参数错误")
	}

	userID, exists := c.Get("user_id")
	if exists && userID != nil {
		field.UserID = userID.(uint)
	}

	if err := h.fieldService.Create(&field); err != nil {
		return nil, err
	}
	return field, nil
}

func (h *Handler) Update(c *gin.Context) (interface{}, error) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var field fieldModel.FieldSchema
	if err := c.ShouldBindJSON(&field); err != nil {
		return nil, errors.New("参数错误")
	}

	field.ID = uint(id)
	if err := h.fieldService.Update(&field); err != nil {
		return nil, err
	}
	return field, nil
}

func (h *Handler) Delete(c *gin.Context) (interface{}, error) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.fieldService.Delete(uint(id)); err != nil {
		return nil, err
	}
	return gin.H{"message": "删除成功"}, nil
}

// HandleFieldAPI 注册路由
func (h *Handler) HandleFieldAPI(r *gin.RouterGroup) {
	j := httputil.NewJSONHandler(r)
	j.GET("", h.List)
	j.GET("/:id", h.Get)
	j.POST("", h.Create)
	j.PUT("/:id", h.Update)
	j.DELETE("/:id", h.Delete)
}
