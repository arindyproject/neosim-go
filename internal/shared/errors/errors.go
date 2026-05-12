package errors

import (
	"errors"
	"net/http"

	"neosim_go/internal/shared/response"

	"github.com/labstack/echo/v5"
)

// ─── App Error ─────────────────────────────────────────────────────────────────

// AppError adalah standard error dengan HTTP status code
type AppError struct {
	Code    int
	Message string
	Err     error // original error (opsional, untuk logging)
}

func (e *AppError) Error() string   { return e.Message }
func (e *AppError) Unwrap() error   { return e.Err }
func (e *AppError) StatusCode() int { return e.Code } // implementasi echo.HTTPStatusCoder

// ─── Constructor Shortcuts ─────────────────────────────────────────────────────

func BadRequest(message string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: message}
}

func Unauthorized(message string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: message}
}

func Forbidden(message string) *AppError {
	return &AppError{Code: http.StatusForbidden, Message: message}
}

func NotFound(message string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: message}
}

func UnprocessableEntity(message string) *AppError {
	return &AppError{Code: http.StatusUnprocessableEntity, Message: message}
}

func TooManyRequests(message string) *AppError {
	return &AppError{Code: http.StatusTooManyRequests, Message: message}
}

func Internal(message string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: message}
}

// Wrap membungkus error asli dengan AppError
func Wrap(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// ─── Global Error Handler ──────────────────────────────────────────────────────

// Handler adalah global error handler untuk Echo v5
// Daftarkan di main.go: e.HTTPErrorHandler = errors.Handler
func Handler(c *echo.Context, err error) {
	// Cek apakah response sudah dikirim
	if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
		if resp.Committed {
			return
		}
	}

	// Ambil HTTP status code dari error chain (menggunakan echo.HTTPStatusCoder)
	code := http.StatusInternalServerError
	var sc echo.HTTPStatusCoder
	if errors.As(err, &sc) {
		if tmp := sc.StatusCode(); tmp != 0 {
			code = tmp
		}
	}

	// Ambil pesan yang sesuai
	message := resolveMessage(err, code)

	// HEAD request tidak boleh ada body
	if c.Request().Method == http.MethodHead {
		if cErr := c.NoContent(code); cErr != nil {
			c.Logger().Error("failed to send no content", "error", errors.Join(err, cErr))
		}
		return
	}

	// Kirim JSON response
	if cErr := response.Response(c, code, false, message, nil, nil); cErr != nil {
		c.Logger().Error("failed to send error response", "error", errors.Join(err, cErr))
	}
}

// ─── Helper ────────────────────────────────────────────────────────────────────

// resolveMessage mengambil pesan error yang sesuai dari error chain
func resolveMessage(err error, code int) string {
	// 1. AppError — pesan dari aplikasi kita
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Message
	}

	// 2. Echo HTTPError — Di v5, Message sudah bertipe string
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		if httpErr.Message != "" {
			return httpErr.Message
		}
	}

	// 3. Fallback berdasarkan status code
	return defaultMessage(code)
}

// defaultMessage mengembalikan pesan default berdasarkan HTTP status code
func defaultMessage(code int) string {
	switch code {
	case http.StatusBadRequest:
		return "Request tidak valid"
	case http.StatusUnauthorized:
		return "Autentikasi diperlukan"
	case http.StatusForbidden:
		return "Akses ditolak"
	case http.StatusNotFound:
		return "Endpoint tidak ditemukan (404)"
	case http.StatusMethodNotAllowed:
		return "Method tidak diizinkan"
	case http.StatusRequestEntityTooLarge:
		return "Ukuran request terlalu besar"
	case http.StatusUnsupportedMediaType:
		return "Tipe konten tidak didukung"
	case http.StatusUnprocessableEntity:
		return "Validasi gagal"
	case http.StatusTooManyRequests:
		return "Terlalu banyak permintaan, coba lagi nanti"
	case http.StatusInternalServerError:
		return "Terjadi kesalahan sistem"
	case http.StatusServiceUnavailable:
		return "Layanan tidak tersedia"
	default:
		return "Terjadi kesalahan"
	}
}
