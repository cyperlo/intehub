package field

import (
	model "intehub/internal/app/models"
	"intehub/internal/app/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	fieldService service.FieldService
}

func NewHandler(fieldService service.FieldService) *Handler {
	return &Handler{fieldService: fieldService}
}

func (h *Handler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var uid uint
	if role != "admin" {
		uid = userID.(uint)
	}

	fields, err := h.fieldService.GetAll(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fields)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	field, err := h.fieldService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "字段不存在"})
		return
	}
	c.JSON(http.StatusOK, field)
}

func (h *Handler) Create(c *gin.Context) {
	var field model.FieldSchema
	if err := c.ShouldBindJSON(&field); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	field.UserID = userID.(uint)

	if err := h.fieldService.Create(&field); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, field)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var field model.FieldSchema
	if err := c.ShouldBindJSON(&field); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	field.ID = uint(id)
	if err := h.fieldService.Update(&field); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, field)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.fieldService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
