package pendidikan

import (
	"neosim_go/internal/modules/kepegawaian/pendidikan/handlers"
	authMiddlewares "neosim_go/internal/modules/auth/middlewares"
	"neosim_go/internal/shared/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, h *handlers.KepegawaianPendidikanHandler, jwtManager *utils.JWTManager, db *gorm.DB) {
	jwt := authMiddlewares.JWTMiddleware(jwtManager, db)
	g := e.Group("/api/v1/kepegawaian/pendidikan", jwt)
	g.GET("", h.ListPendidikan)
	g.GET("/:id", h.GetPendidikanByID)
	g.POST("", h.CreatePendidikan)
	g.PUT("/:id", h.UpdatePendidikan)
	g.DELETE("/:id", h.DeletePendidikan)
	gJenjang := e.Group("/api/v1/kepegawaian/pendidikan/jenjangs", jwt)
	gJenjang.GET("", h.ListJenjang)
	gJenjang.GET("/:id", h.GetJenjangByID)
	gJenjang.POST("", h.CreateJenjang)
	gJenjang.PUT("/:id", h.UpdateJenjang)
	gJenjang.DELETE("/:id", h.DeleteJenjang)
	// GEN:ITEM_ROUTES
}
