package apps

import (
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterModules(db *gorm.DB, e *echo.Echo) {
	InjectDB(db)
	InitAllRoutes(e)
}
