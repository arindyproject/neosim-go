package pegawai

import (
	"neosim_go/internal/modules/kepegawaian/pegawai/handlers"
	authMiddlewares "neosim_go/internal/modules/auth/middlewares"
	"neosim_go/internal/shared/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, h *handlers.KepegawaianPegawaiHandler, jwtManager *utils.JWTManager, db *gorm.DB) {
	jwt := authMiddlewares.JWTMiddleware(jwtManager, db)
	g := e.Group("/api/v1/kepegawaian/pegawai", jwt)
	g.GET("", h.ListPegawai)
	g.GET("/:id", h.GetPegawaiByID)
	g.POST("", h.CreatePegawai)
	g.PUT("/:id", h.UpdatePegawai)
	g.DELETE("/:id", h.DeletePegawai)
	// GEN:ITEM_ROUTES
}
