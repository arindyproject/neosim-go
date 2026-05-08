package main

import (
	"log"
	"net/http"

	"neosim_go/config"
	"neosim_go/internal/apps"

	// =====================================================================
	// import module di sini
	// =====================================================================
	_ "neosim_go/internal/modules/users"

	// =====================================================================

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"neosim_go/internal/shared/response"
)

func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig("DEV") // Ganti ke "PROD" untuk production

	// 2. Echo Instance
	e := echo.New()

	// 3. Global Middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(cfg.CORS))

	// 4. Health Check Route
	e.GET("/", func(c *echo.Context) error {
		return response.Response(c, http.StatusOK, true, "Hello, World!", nil, nil)
	})

	// 5. Database Connection
	db, err := cfg.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer config.CloseDB(db)

	// 6. Register Semua Module (routes, handler, service, repo)
	apps.RegisterModules(db, e)

	// Debug: uncomment untuk verifikasi semua routes terdaftar
	for _, r := range e.Router().Routes() {
		log.Printf("%-7s %s", r.Method, r.Path)
	}

	// 7. Start Server
	if err := e.Start(":" + cfg.ServerPort); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
