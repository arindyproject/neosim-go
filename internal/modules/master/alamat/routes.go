package alamat

import (
	authMiddlewares "neosim_go/internal/modules/auth/middlewares"
	"neosim_go/internal/modules/master/alamat/handlers"
	"neosim_go/internal/shared/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// RegisterRoutes mendaftarkan semua routes untuk module master/alamat
func RegisterRoutes(e *echo.Echo, h *handlers.MasterAlamatHandler, jwtManager *utils.JWTManager, db *gorm.DB) {
	jwt := authMiddlewares.JWTMiddleware(jwtManager, db)
	g := e.Group("/api/v1/master/alamat", jwt)
	// Negara ---------------------------------------------------------
	negara := g.Group("/negara")
	negara.GET("", h.ListNegara)
	negara.GET("/:id", h.GetByIDNegara)
	negara.POST("", h.CreateNegara)
	negara.PUT("/:id", h.UpdateNegara)
	negara.DELETE("/:id", h.DeleteNegara)

	// Provinsi ---------------------------------------------------------
	provinsi := g.Group("/provinsi")
	provinsi.GET("", h.ListProvinsi)
	provinsi.GET("/:id", h.GetByIDProvinsi)
	provinsi.POST("", h.CreateProvinsi)
	provinsi.PUT("/:id", h.UpdateProvinsi)
	provinsi.DELETE("/:id", h.DeleteProvinsi)

	// Kota/Kabupaten -----------------------------------------------------
	kota := g.Group("/kota")
	kota.GET("", h.ListKotaKabupaten)
	kota.GET("/:id", h.GetByIDKotaKabupaten)
	kota.POST("", h.CreateKotaKabupaten)
	kota.PUT("/:id", h.UpdateKotaKabupaten)
	kota.DELETE("/:id", h.DeleteKotaKabupaten)

	// Kecamatan -----------------------------------------------------
	kecamatan := g.Group("/kecamatan")
	kecamatan.GET("", h.ListKecamatan)
	kecamatan.GET("/:id", h.GetByIDKecamatan)
	kecamatan.POST("", h.CreateKecamatan)
	kecamatan.PUT("/:id", h.UpdateKecamatan)
	kecamatan.DELETE("/:id", h.DeleteKecamatan)

	// Kelurahan/Desa -----------------------------------------------------
	desa := g.Group("/desa")
	desa.GET("", h.ListKelurahanDesa)
	desa.GET("/:id", h.GetByIDKelurahanDesa)
	desa.POST("", h.CreateKelurahanDesa)
	desa.PUT("/:id", h.UpdateKelurahanDesa)
	desa.DELETE("/:id", h.DeleteKelurahanDesa)

}
