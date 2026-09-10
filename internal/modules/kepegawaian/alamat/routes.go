package alamat

import (
	"neosim_go/internal/modules/kepegawaian/alamat/handlers"
	authMiddlewares "neosim_go/internal/modules/auth/middlewares"
	"neosim_go/internal/shared/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, h *handlers.KepegawaianAlamatHandler, jwtManager *utils.JWTManager, db *gorm.DB) {
	jwt := authMiddlewares.JWTMiddleware(jwtManager, db)
	g := e.Group("/api/v1/kepegawaian/alamat", jwt)
	g.GET("", h.ListAlamat)
	g.GET("/:id", h.GetAlamatByID)
	g.POST("", h.CreateAlamat)
	g.PUT("/:id", h.UpdateAlamat)
	g.DELETE("/:id", h.DeleteAlamat)
	gTipe := e.Group("/api/v1/kepegawaian/alamat/tipes", jwt)
	gTipe.GET("", h.ListTipe)
	gTipe.GET("/:id", h.GetTipeByID)
	gTipe.POST("", h.CreateTipe)
	gTipe.PUT("/:id", h.UpdateTipe)
	gTipe.DELETE("/:id", h.DeleteTipe)
	// GEN:ITEM_ROUTES
}
