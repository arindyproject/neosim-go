package kontak

import (
	"neosim_go/internal/modules/kepegawaian/kontak/handlers"
	authMiddlewares "neosim_go/internal/modules/auth/middlewares"
	"neosim_go/internal/shared/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, h *handlers.KepegawaianKontakHandler, jwtManager *utils.JWTManager, db *gorm.DB) {
	jwt := authMiddlewares.JWTMiddleware(jwtManager, db)
	g := e.Group("/api/v1/kepegawaian/kontak", jwt)
	g.GET("", h.ListKontak)
	g.GET("/:id", h.GetKontakByID)
	g.POST("", h.CreateKontak)
	g.PUT("/:id", h.UpdateKontak)
	g.DELETE("/:id", h.DeleteKontak)
	// GEN:ITEM_ROUTES
}
