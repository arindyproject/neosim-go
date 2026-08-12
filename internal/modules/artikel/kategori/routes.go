package kategori

import (
	"neosim_go/internal/modules/artikel/kategori/handlers"
	authMiddlewares "neosim_go/internal/modules/auth/middlewares"
	"neosim_go/internal/shared/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, h *handlers.ArtikelKategoriHandler, jwtManager *utils.JWTManager, db *gorm.DB) {
	jwt := authMiddlewares.JWTMiddleware(jwtManager, db)
	g := e.Group("/api/v1/artikel/kategori", jwt)
	g.GET("", h.ListKategori)
	g.GET("/:id", h.GetKategoriByID)
	g.POST("", h.CreateKategori)
	g.PUT("/:id", h.UpdateKategori)
	g.DELETE("/:id", h.DeleteKategori)
	gTag := e.Group("/api/v1/artikel/kategori/tags", jwt)
	gTag.GET("", h.ListTag)
	gTag.GET("/:id", h.GetTagByID)
	gTag.POST("", h.CreateTag)
	gTag.PUT("/:id", h.UpdateTag)
	gTag.DELETE("/:id", h.DeleteTag)
	// GEN:ITEM_ROUTES
}
