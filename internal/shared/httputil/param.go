package httputil

import (
	"neosim_go/config"
	"strconv"

	"github.com/labstack/echo/v5"
)

var DefaultPageSizeMax = 100

func ParseID(c *echo.Context) (int64, error) {
	return strconv.ParseInt(c.Param("id"), 10, 64)
}

func ParsePagination(c *echo.Context, cfg *config.Config) (page, pageSize int) {
	// Fallback jika config belum diset
	defaultPage := 1
	defaultPageSize := 10
	maxPage := 1000
	maxPageSize := 100

	if cfg != nil {
		if cfg.DefaultPageSize > 0 {
			defaultPageSize = cfg.DefaultPageSize
		}
		if cfg.DefaultPageSizeMax > 0 {
			maxPageSize = cfg.DefaultPageSizeMax
		}
		if cfg.DefaultPageSizeMax > 0 {
			maxPage = cfg.DefaultPageSizeMax
		}
	}

	page, pageSize = defaultPage, defaultPageSize

	if p := c.QueryParam("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.QueryParam("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	// Clamp ke batas maksimum dari config
	if page > maxPage {
		page = maxPage
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	return page, pageSize
}

func ParseOptionalInt64Query(c *echo.Context, key string) *int64 {
	raw := c.QueryParam(key)
	if raw == "" {
		return nil
	}
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return nil
	}
	return &v
}
