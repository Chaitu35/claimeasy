package clinics

import "github.com/gin-gonic/gin"

// RegisterClinicRoutes mounts clinic routes to the router
func RegisterClinicRoutes(rg *gin.RouterGroup, h *Handler) {
	
	clinics := rg.Group("/clinics")

	clinics.POST("", h.CreateClinicHandler)
	clinics.PUT("/:id", h.UpdateClinicHandler)
	clinics.GET("", h.ListClinicsHandler)
	clinics.POST("/:id/assign", h.AssignUserHandler)
	clinics.PATCH("/:id/status", h.UpdateClinicStatusHandler)
}
