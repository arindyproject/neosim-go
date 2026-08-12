package artikel

import (
	"neosim_go/internal/modules/artikel/artikel/handlers"
	authMiddlewares "neosim_go/internal/modules/auth/middlewares"
	"neosim_go/internal/shared/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, h *handlers.ArtikelHandler, jwtManager *utils.JWTManager, db *gorm.DB) {
	jwt := authMiddlewares.JWTMiddleware(jwtManager, db)
	g := e.Group("/api/v1/artikel", jwt)
	g.GET("", h.ListArtikel)
	g.GET("/:id", h.GetArtikelByID)
	g.POST("", h.CreateArtikel)
	g.PUT("/:id", h.UpdateArtikel)
	g.DELETE("/:id", h.DeleteArtikel)
	// GEN:ITEM_ROUTES
}
