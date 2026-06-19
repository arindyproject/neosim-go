package ketegori

import (
	"neosim_go/internal/modules/artikel/ketegori/handlers"
	authMiddlewares "neosim_go/internal/modules/auth/middlewares"
	"neosim_go/internal/shared/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
	
)

// RegisterRoutes mendaftarkan semua routes untuk module artikel/ketegori
func RegisterRoutes(e *echo.Echo, h *handlers.ArtikelKetegoriHandler, jwtManager *utils.JWTManager, db *gorm.DB) {
	jwt := authMiddlewares.JWTMiddleware(jwtManager, db)
	g := e.Group("/api/v1/artikel/ketegori", jwt)
	g.GET("", h.List)
	g.GET("/:id", h.GetByID)
	g.POST("", h.Create)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}
