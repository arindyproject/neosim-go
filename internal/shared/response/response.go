package response

import (
	"github.com/labstack/echo/v5"
)

// ─── Response Structs ──────────────────────────────────────────────────────────

type MyGoResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type Pagination struct {
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	NextPage    *int  `json:"next_page"`
	PrevPage    *int  `json:"prev_page"`
}

type PaginatedData struct {
	Items      interface{} `json:"items"`
	Pagination Pagination  `json:"pagination"`
}

type MyGoPaginatedResponse struct {
	Status  bool          `json:"status"`
	Message string        `json:"message"`
	Data    PaginatedData `json:"data"`
}

// ─── Helpers ───────────────────────────────────────────────────────────────────

func buildPagination(total int64, page, perPage int) Pagination {
	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	var nextPage *int
	if page < totalPages {
		n := page + 1
		nextPage = &n
	}

	var prevPage *int
	if page > 1 {
		p := page - 1
		prevPage = &p
	}

	return Pagination{
		TotalItems:  total,
		TotalPages:  totalPages,
		CurrentPage: page,
		PerPage:     perPage,
		NextPage:    nextPage,
		PrevPage:    prevPage,
	}
}

// ─── Response Functions ────────────────────────────────────────────────────────

// Response mengirim standard JSON response
// status http ditentukan dari handler
func Response(c *echo.Context, httpStatus int, status bool, message string, data interface{}, errors interface{}) error {
	return c.JSON(httpStatus, MyGoResponse{
		Status:  status,
		Message: message,
		Data:    data,
		Errors:  errors,
	})
}

// Paginated mengirim paginated JSON response
func Paginated(c *echo.Context, httpStatus int, status bool, message string, items interface{}, total int64, page, perPage int) error {
	return c.JSON(httpStatus, MyGoPaginatedResponse{
		Status:  status,
		Message: message,
		Data: PaginatedData{
			Items:      items,
			Pagination: buildPagination(total, page, perPage),
		},
	})
}
