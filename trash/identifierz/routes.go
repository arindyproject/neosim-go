package identifier

import (
	authMiddlewares "neosim_go/internal/modules/auth/middlewares"
	"neosim_go/internal/modules/kepegawaian/identifier/handlers"
	"neosim_go/internal/shared/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, h *handlers.KepegawaianIdentifierHandler, jwtManager *utils.JWTManager, db *gorm.DB) {
	jwt := authMiddlewares.JWTMiddleware(jwtManager, db)
	g := e.Group("/api/v1/kepegawaian/identifier", jwt)
	g.GET("", h.List)
	g.POST("", h.Create)
	g.GET("/expiring-soon", h.GetExpiringSoon)
	g.GET("/types", h.GetIdentifierTypes)
	g.GET("/:pegawai_id/pegawai", h.ListByPegawai)
	g.GET("/:id", h.GetByID)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)

}
