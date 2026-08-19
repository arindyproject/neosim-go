package departemen

import (
	"neosim_go/internal/modules/master/departemen/handlers"
	authMiddlewares "neosim_go/internal/modules/auth/middlewares"
	"neosim_go/internal/shared/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, h *handlers.MasterDepartemenHandler, jwtManager *utils.JWTManager, db *gorm.DB) {
	jwt := authMiddlewares.JWTMiddleware(jwtManager, db)
	g := e.Group("/api/v1/master/departemen", jwt)
	g.GET("", h.ListDepartemen)
	g.GET("/:id", h.GetDepartemenByID)
	g.POST("", h.CreateDepartemen)
	g.PUT("/:id", h.UpdateDepartemen)
	g.DELETE("/:id", h.DeleteDepartemen)
	// GEN:ITEM_ROUTES
}
