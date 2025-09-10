package clinics

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler wraps the clinic service
type Handler struct {
	service *Service
}

// NewHandler creates a new clinics handler
func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

// CreateClinicHandler → POST /clinics
func (h *Handler) CreateClinicHandler(c *gin.Context) {
	var dto CreateClinicDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clinic, err := h.service.CreateClinic(c, &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, clinic)
}

// UpdateClinicHandler → PUT /clinics/:id
func (h *Handler) UpdateClinicHandler(c *gin.Context) {
	var dto UpdateClinicDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	clinic, err := h.service.UpdateClinic(c, uint(id), &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, clinic)
}

// ListClinicsHandler → GET /clinics
func (h *Handler) ListClinicsHandler(c *gin.Context) {
	var query ListClinicsQuery
	if err := c.BindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}

	resp, err := h.service.ListClinics(c, &query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// AssignUserHandler → POST /clinics/:id/assign
func (h *Handler) AssignUserHandler(c *gin.Context) {
	var dto AssignUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.service.AssignUserToClinic(c, uint(id), dto.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user assigned"})
}

// UpdateClinicStatusHandler → PATCH /clinics/:id/status
func (h *Handler) UpdateClinicStatusHandler(c *gin.Context) {
	var dto UpdateClinicStatusDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.service.UpdateClinicStatus(c, uint(id), dto.Active); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "clinic status updated"})
}
