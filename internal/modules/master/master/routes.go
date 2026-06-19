package master

import (
	"neosim_go/internal/modules/master/master/handlers"
	"neosim_go/internal/shared/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// RegisterRoutes mendaftarkan semua routes untuk module master/master
func RegisterRoutes(e *echo.Echo, h *handlers.MasterHandler, jwtManager *utils.JWTManager, db *gorm.DB) {
	//jwt := authMiddlewares.JWTMiddleware(jwtManager, db)
	//g := e.Group("/api/v1/master", jwt)
	//g.GET("", h.List)
	//g.GET("/:id", h.GetByID)
	//g.POST("", h.Create)
	//g.PUT("/:id", h.Update)
	//g.DELETE("/:id", h.Delete)
}
