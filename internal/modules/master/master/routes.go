package master

import (
	authMiddlewares "neosim_go/internal/modules/auth/middlewares"
	"neosim_go/internal/modules/master/master/handlers"
	"neosim_go/internal/shared/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// RegisterRoutes mendaftarkan semua routes untuk module master/master
func RegisterRoutes(e *echo.Echo, h *handlers.MasterHandler, jwtManager *utils.JWTManager, db *gorm.DB) {
	jwt := authMiddlewares.JWTMiddleware(jwtManager, db)
	g := e.Group("/api/v1/master", jwt)

	// Pekerjaan------------------------------------------------------
	g.GET("/pekerjaan", h.ListPekerjaan)
	g.GET("/pekerjaan/:id", h.GetByIDPekerjaan)
	g.POST("/pekerjaan", h.CreatePekerjaan)
	g.PUT("/pekerjaan/:id", h.UpdatePekerjaan)
	g.DELETE("/pekerjaan/:id", h.DeletePekerjaan)

	// Pendidikan------------------------------------------------------
	g.GET("/pendidikan", h.ListPendidikan)
	g.GET("/pendidikan/:id", h.GetByIDPendidikan)
	g.POST("/pendidikan", h.CreatePendidikan)
	g.PUT("/pendidikan/:id", h.UpdatePendidikan)
	g.DELETE("/pendidikan/:id", h.DeletePendidikan)

	// Agama------------------------------------------------------
	g.GET("/agama", h.ListAgama)
	g.GET("/agama/:id", h.GetByIDAgama)
	g.POST("/agama", h.CreateAgama)
	g.PUT("/agama/:id", h.UpdateAgama)
	g.DELETE("/agama/:id", h.DeleteAgama)

	// Status Pernikahan------------------------------------------------------
	g.GET("/status_pernikahan", h.ListStatusPernikahan)
	g.GET("/status_pernikahan/:id", h.GetByIDStatusPernikahan)
	g.POST("/status_pernikahan", h.CreateStatusPernikahan)
	g.PUT("/status_pernikahan/:id", h.UpdateStatusPernikahan)
	g.DELETE("/status_pernikahan/:id", h.DeleteStatusPernikahan)

	// Suku------------------------------------------------------
	g.GET("/suku", h.ListSuku)
	g.GET("/suku/:id", h.GetByIDSuku)
	g.POST("/suku", h.CreateSuku)
	g.PUT("/suku/:id", h.UpdateSuku)
	g.DELETE("/suku/:id", h.DeleteSuku)

	// Golongan Darah------------------------------------------------------
	g.GET("/golongan_darah", h.ListGolonganDarah)
	g.GET("/golongan_darah/:id", h.GetByIDGolonganDarah)
	g.POST("/golongan_darah", h.CreateGolonganDarah)
	g.PUT("/golongan_darah/:id", h.UpdateGolonganDarah)
	g.DELETE("/golongan_darah/:id", h.DeleteGolonganDarah)

	// Jenis Kelamin------------------------------------------------------
	g.GET("/jenis_kelamin", h.ListJenisKelamin)
	g.GET("/jenis_kelamin/:id", h.GetByIDJenisKelamin)
	g.POST("/jenis_kelamin", h.CreateJenisKelamin)
	g.PUT("/jenis_kelamin/:id", h.UpdateJenisKelamin)
	g.DELETE("/jenis_kelamin/:id", h.DeleteJenisKelamin)
}
