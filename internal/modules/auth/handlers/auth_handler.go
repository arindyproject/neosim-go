package handlers

import (
	"net/http"

	"neosim_go/internal/modules/auth/contracts"
	"neosim_go/internal/modules/auth/dto"
	"neosim_go/internal/modules/auth/services"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"

	"github.com/labstack/echo/v5"
)

// AuthHandler menangani HTTP request untuk auth
type AuthHandler struct {
	service contracts.AuthService
}

// NewAuthHandler membuat instance handler baru
func NewAuthHandler(service contracts.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// ─── Login ─────────────────────────────────────────────────────────────────────

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(c *echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	ip := c.RealIP()
	userAgent := c.Request().UserAgent()

	result, err := h.service.Login(&req, ip, userAgent)
	if err != nil {
		return handleAuthError(c, err)
	}

	return response.Response(c, http.StatusOK, true, "Login berhasil", result, nil)
}

// ─── Register ──────────────────────────────────────────────────────────────────

// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(c *echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	result, err := h.service.Register(&req)
	if err != nil {
		return handleAuthError(c, err)
	}

	return response.Response(c, http.StatusCreated, true, "Registrasi berhasil", result, nil)
}

// ─── Refresh Token ─────────────────────────────────────────────────────────────

// RefreshToken handles POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *echo.Context) error {
	var req dto.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	result, err := h.service.RefreshToken(&req)
	if err != nil {
		return handleAuthError(c, err)
	}

	return response.Response(c, http.StatusOK, true, "Token berhasil diperbarui", result, nil)
}

// ─── Forgot Password ───────────────────────────────────────────────────────────

// ForgotPassword handles POST /api/v1/auth/forgot-password
func (h *AuthHandler) ForgotPassword(c *echo.Context) error {
	var req dto.ForgotPasswordRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	// Selalu return 200 — jangan bocorkan apakah identifier terdaftar
	h.service.ForgotPassword(&req)

	return response.Response(c, http.StatusOK, true,
		"Jika email terdaftar, kami telah mengirimkan instruksi reset password.", nil, nil)
}

// ─── Reset Password ────────────────────────────────────────────────────────────

// ResetPassword handles POST /api/v1/auth/reset-password
func (h *AuthHandler) ResetPassword(c *echo.Context) error {
	var req dto.ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	if err := h.service.ResetPassword(&req); err != nil {
		return handleAuthError(c, err)
	}

	return response.Response(c, http.StatusOK, true,
		"Password berhasil direset. Silakan login dengan password baru.", nil, nil)
}

// ─── Helper ────────────────────────────────────────────────────────────────────

// handleAuthError mengubah AuthError menjadi response yang sesuai
func handleAuthError(c *echo.Context, err error) error {
	if authErr, ok := err.(*services.AuthError); ok {
		return response.Response(c, authErr.Code, false, authErr.Message, nil, nil)
	}
	return response.Response(c, http.StatusInternalServerError, false, "Terjadi kesalahan sistem.", nil, nil)
}
