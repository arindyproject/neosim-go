package identifier

import (
	"neosim_go/internal/modules/kepegawaian/identifier/handlers"
	authMiddlewares "neosim_go/internal/modules/auth/middlewares"
	"neosim_go/internal/shared/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, h *handlers.KepegawaianIdentifierHandler, jwtManager *utils.JWTManager, db *gorm.DB) {
	jwt := authMiddlewares.JWTMiddleware(jwtManager, db)
	g := e.Group("/api/v1/kepegawaian/identifier", jwt)
	g.GET("", h.ListIdentifier)
	g.GET("/:id", h.GetIdentifierByID)
	g.POST("", h.CreateIdentifier)
	g.PUT("/:id", h.UpdateIdentifier)
	g.DELETE("/:id", h.DeleteIdentifier)
	gTipe := e.Group("/api/v1/kepegawaian/identifier/tipes", jwt)
	gTipe.GET("", h.ListTipe)
	gTipe.GET("/:id", h.GetTipeByID)
	gTipe.POST("", h.CreateTipe)
	gTipe.PUT("/:id", h.UpdateTipe)
	gTipe.DELETE("/:id", h.DeleteTipe)
	// GEN:ITEM_ROUTES
}
