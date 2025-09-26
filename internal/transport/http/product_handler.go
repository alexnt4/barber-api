package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/alexnt4/barber-api/internal/domain"
	"github.com/alexnt4/barber-api/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc}
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required, gt=0"`
	Description string  `json:"description"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required, gt=0"`
	Description string  `json:"description"`
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := &domain.Product{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
	}

	if err := h.svc.Create(c.Request.Context(), product); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) List(c *gin.Context) {
	products, err := h.svc.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    len(products),
	})
}

func (h *ProductHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	appt, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == domain.ErrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "producto no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appt)
}

func (h *AppointmentHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	var req UpdateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parseo de fechas
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formato de fecha invalido para start_time, use RFC3339"})
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formato de fecha invalido para end_time, use RFC3339"})
		return
	}

	// Mapear IDs a domain.Product
	products := make([]domain.Product, len(req.Products))
	for i, id := range req.Products {
		products[i] = domain.Product{ID: id}
	}

	appt := &domain.Appointment{
		ID:          uint(id),
		ClienteName: req.ClienteName,
		StartTime:   startTime,
		EndTime:     endTime,
		Products:    products,
	}

	if err := h.svc.Update(c.Request.Context(), uint(id), appt); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appt)
}

func (h *AppointmentHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	if err := h.svc.Cancel(c.Request.Context(), uint(id)); err != nil {
		if err == domain.ErrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "cita no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cita cancelada exitosamente"})
}

func (h *AppointmentHandler) GetTotal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	total, err := h.svc.GetTotalPrice(c.Request.Context(), uint(id))
	if err != nil {
		if err == domain.ErrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "cita no encontrada"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"appoinmet_id": id,
		"total":        total,
	})
}
