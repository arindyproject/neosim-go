package pegawai

import (
	authMiddlewares "neosim_go/internal/modules/auth/middlewares"
	"neosim_go/internal/modules/kepegawaian/pegawai/handlers"
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

	gStatus := e.Group("/api/v1/kepegawaian/pegawai/statuss", jwt)
	gStatus.GET("", h.ListStatus)
	gStatus.GET("/select", h.ListSelectStatus)
	gStatus.GET("/:id", h.GetStatusByID)
	gStatus.POST("", h.CreateStatus)
	gStatus.PUT("/:id", h.UpdateStatus)
	gStatus.DELETE("/:id", h.DeleteStatus)

	gJenis := e.Group("/api/v1/kepegawaian/pegawai/jeniss", jwt)
	gJenis.GET("", h.ListJenis)
	gJenis.GET("/select", h.ListSelectJenis)
	gJenis.GET("/:id", h.GetJenisByID)
	gJenis.POST("", h.CreateJenis)
	gJenis.PUT("/:id", h.UpdateJenis)
	gJenis.DELETE("/:id", h.DeleteJenis)
	// GEN:ITEM_ROUTES
}
