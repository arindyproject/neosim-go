package kualifikasi

import (
	authMiddlewares "neosim_go/internal/modules/auth/middlewares"
	"neosim_go/internal/modules/kepegawaian/kualifikasi/handlers"
	"neosim_go/internal/shared/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, h *handlers.KepegawaianKualifikasiHandler, jwtManager *utils.JWTManager, db *gorm.DB) {
	jwt := authMiddlewares.JWTMiddleware(jwtManager, db)

	g := e.Group("/api/v1/kepegawaian/kualifikasi", jwt)
	g.GET("", h.ListKualifikasi)
	g.GET("/:id", h.GetKualifikasiByID)
	g.GET("/expiring_soon", h.GetExpiringSoonKualifikasi)
	g.GET("/expired", h.GetExpiredKualifikasi)
	g.GET("/:pegawai_id/pegawai", h.ListKualifikasiByPegawai)
	g.POST("", h.CreateKualifikasi)
	g.PUT("/:id", h.UpdateKualifikasi)
	g.DELETE("/:id", h.DeleteKualifikasi)

	gTipe := e.Group("/api/v1/kepegawaian/kualifikasi/tipes", jwt)
	gTipe.GET("", h.ListTipe)
	gTipe.GET("/:id", h.GetTipeByID)
	gTipe.POST("", h.CreateTipe)
	gTipe.PUT("/:id", h.UpdateTipe)
	gTipe.DELETE("/:id", h.DeleteTipe)
	// GEN:ITEM_ROUTES
}
