package permission

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
	db      *gorm.DB
}

func NewHandler(service Service, db *gorm.DB) *Handler {
	return &Handler{
		service: service,
		db:      db,
	}
}

// -----------------------------
// Permission CRUD Endpoints
// -----------------------------

func (h *Handler) CreatePermission(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.service.CreatePermission(c.Request.Context(), h.db, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) UpdatePermission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission ID"})
		return
	}
	var req UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}	
	resp,err:= h.service.UpdatePermission(c.Request.Context(), h.db, uint(id),req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}	
	c.JSON(http.StatusOK, resp)
}


func (h *Handler) GetPermissionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission ID"})
		return
	}
	resp, err := h.service.GetPermissionByID(c.Request.Context(), h.db, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ListPermissions(c *gin.Context) {
	resp, err := h.service.ListPermissions(c.Request.Context(), h.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// -----------------------------
// Role ↔ Permission Endpoints
// -----------------------------

func (h *Handler) AssignPermissionToRole(c *gin.Context) {
	var req AssignRolePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.service.AssignPermissionToRole(c.Request.Context(), h.db, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) RevokePermissionFromRole(c *gin.Context) {
	roleID, err1 := strconv.Atoi(c.Param("roleID"))
	permID, err2 := strconv.Atoi(c.Param("permID"))
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid roleID or permissionID"})
		return
	}

	var body struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason is required"})
		return
	}

	if err := h.service.RevokePermissionFromRole(c.Request.Context(), h.db, uint(roleID), uint(permID), body.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "permission revoked from role"})
}


func (h *Handler) ListRolePermissions(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("roleID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}
	resp, err := h.service.ListRolePermissions(c.Request.Context(), h.db, uint(roleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// -----------------------------
// User ↔ Permission Endpoints
// -----------------------------

func (h *Handler) AssignPermissionToUser(c *gin.Context) {
	var req AssignUserPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.service.AssignPermissionToUser(c.Request.Context(), h.db, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) RevokePermissionFromUser(c *gin.Context) {
	userID, err1 := strconv.Atoi(c.Param("userID"))
	permID, err2 := strconv.Atoi(c.Param("permID"))
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID or permissionID"})
		return
	}

	var body struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason is required"})
		return
	}

	if err := h.service.RevokePermissionFromUser(c.Request.Context(), h.db, uint(userID), uint(permID), body.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "permission revoked from user"})
}

func (h *Handler) ListUserPermissions(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	resp, err := h.service.ListUserPermissions(c.Request.Context(), h.db, uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
