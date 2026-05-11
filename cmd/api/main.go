package main

import (
	"log"
	"net/http"

	"neosim_go/config"
	"neosim_go/internal/apps"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	// =============== Modules ===============
	// input module baru di sini saat ada, misal:
	_ "neosim_go/internal/modules/auth"
	_ "neosim_go/internal/modules/users"
	// =============== Modules ===============
)

func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig("DEV")

	// 2. Echo Instance
	e := echo.New()

	// 3. Global Middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(cfg.CORS))

	// 4. Health Check
	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"message": "Neosim API is running",
		})
	})

	// 5. Database
	db, err := cfg.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer config.CloseDB(db)

	// 6. Redis
	redisClient, err := cfg.ConnectRedis()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer config.CloseRedis(redisClient)

	// 7. Register semua module — inject cfg, db, redis sekaligus
	apps.RegisterModules(cfg, db, redisClient, e)

	// Debug: uncomment untuk cek routes terdaftar
	for _, r := range e.Router().Routes() {
		log.Printf("-------------------------------------")
		log.Printf("%-7s %s", r.Method, r.Path)
	}

	// 8. Start Server
	if err := e.Start(":" + cfg.ServerPort); err != nil {

		e.Logger.Error("failed to start server", "error", err)
	}
}
