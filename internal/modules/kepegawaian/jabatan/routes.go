package jabatan

import (
	authMiddlewares "neosim_go/internal/modules/auth/middlewares"
	"neosim_go/internal/modules/kepegawaian/jabatan/handlers"
	"neosim_go/internal/shared/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, h *handlers.KepegawaianJabatanHandler, jwtManager *utils.JWTManager, db *gorm.DB) {
	jwt := authMiddlewares.JWTMiddleware(jwtManager, db)
	g := e.Group("/api/v1/kepegawaian/jabatan", jwt)
	g.GET("", h.ListJabatan)
	g.GET("/:id", h.GetJabatanByID)
	g.POST("", h.CreateJabatan)
	g.PUT("/:id", h.UpdateJabatan)
	g.DELETE("/:id", h.DeleteJabatan)
	//--------------------------------------------------------------------
	gPosition := e.Group("/api/v1/kepegawaian/jabatan/positions", jwt)
	gPosition.GET("", h.ListPosition)
	gPosition.GET("/:id", h.GetPositionByID)
	gPosition.GET("/tree", h.GetPositionTree)
	gPosition.POST("", h.CreatePosition)
	gPosition.PUT("/:id", h.UpdatePosition)
	gPosition.DELETE("/:id", h.DeletePosition)
	//--------------------------------------
	gPositionKategori := e.Group("/api/v1/kepegawaian/jabatan/position_kategoris", jwt)
	gPositionKategori.GET("", h.ListPositionKategori)
	gPositionKategori.GET("/:id", h.GetPositionKategoriByID)
	gPositionKategori.GET("/select", h.ListSelectPositionKategori)
	gPositionKategori.POST("", h.CreatePositionKategori)
	gPositionKategori.PUT("/:id", h.UpdatePositionKategori)
	gPositionKategori.DELETE("/:id", h.DeletePositionKategori)
	//--------------------------------------------------------------------
	gJobTitle := e.Group("/api/v1/kepegawaian/jabatan/job_titles", jwt)
	gJobTitle.GET("", h.ListJobTitle)
	gJobTitle.GET("/:id", h.GetJobTitleByID)
	gJobTitle.POST("", h.CreateJobTitle)
	gJobTitle.PUT("/:id", h.UpdateJobTitle)
	gJobTitle.DELETE("/:id", h.DeleteJobTitle)
	//--------------------------------------------------------------------
	gSpecialization := e.Group("/api/v1/kepegawaian/jabatan/specializations", jwt)
	gSpecialization.GET("", h.ListSpecialization)
	gSpecialization.GET("/:id", h.GetSpecializationByID)
	gSpecialization.POST("", h.CreateSpecialization)
	gSpecialization.PUT("/:id", h.UpdateSpecialization)
	gSpecialization.DELETE("/:id", h.DeleteSpecialization)
	//--------------------------------------------------------------------

	// GEN:ITEM_ROUTES
}
