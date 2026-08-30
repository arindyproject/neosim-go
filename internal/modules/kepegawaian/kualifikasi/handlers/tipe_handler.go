package handlers

import (
	"io"
	"net/http"

	"neosim_go/internal/modules/kepegawaian/kualifikasi/dto"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"
	"neosim_go/internal/shared/binding"

	"github.com/labstack/echo/v5"
)

// Method di bawah ini ditempelkan ke struct KepegawaianKualifikasiHandler yang
// sama dengan handler entitas utama (lihat handlers/handler.go). Nama method
// diberi suffix Tipe agar tidak bentrok dengan method entitas utama
// pada struct handler yang sama.

// ─── ListTipe ──────────────────────────────────────────────────────
//
//	@Summary		Get list of Tipe
//	@Description	Get paginated list of Tipe
//	@Tags			kepegawaian/kualifikasi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.TipeResponse}
//	@Router			/kepegawaian/kualifikasi/tipes [get]
func (h *KepegawaianKualifikasiHandler) ListTipe(c *echo.Context) error {
	filter := dto.FilterTipeRequest{Name: c.QueryParam("name")}
	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListTipe(c.Request().Context(),page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetTipeByID ───────────────────────────────────────────────────
//
//	@Summary		Get Tipe
//	@Description	Get Tipe by :id
//	@Tags			kepegawaian/kualifikasi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Tipe ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.TipeResponse}
//	@Router			/kepegawaian/kualifikasi/tipes/{id} [get]
func (h *KepegawaianKualifikasiHandler) GetTipeByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetTipeByID(c.Request().Context(),id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreateTipe ────────────────────────────────────────────────────
//
//	@Summary		Create Tipe
//	@Description	Create New Tipe
//	@Tags			kepegawaian/kualifikasi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateTipeRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.TipeResponse}
//	@Router			/kepegawaian/kualifikasi/tipes [post]
func (h *KepegawaianKualifikasiHandler) CreateTipe(c *echo.Context) error {
	var req dto.CreateTipeRequest
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Gagal membaca request body", nil, err.Error())
	}

	if errs := binding.BindErrors(body, &req); len(errs) > 0 {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal (binding)", nil, errs)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal (validator)", nil, errs)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.CreateTipe(c.Request().Context(),&req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdateTipe ────────────────────────────────────────────────────
//
//	@Summary		Update Tipe
//	@Description	Update Tipe by :id
//	@Tags			kepegawaian/kualifikasi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"Tipe ID"
//	@Param			body	body		dto.UpdateTipeRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.TipeResponse}
//	@Router			/kepegawaian/kualifikasi/tipes/{id} [put]
func (h *KepegawaianKualifikasiHandler) UpdateTipe(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateTipeRequest
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Gagal membaca request body", nil, err.Error())
	}

	if errs := binding.BindErrors(body, &req); len(errs) > 0 {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal (binding)", nil, errs)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal (validator)", nil, errs)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.UpdateTipe(c.Request().Context(),id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "Tipe tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeleteTipe ────────────────────────────────────────────────────
//
//	@Summary		Delete Tipe
//	@Description	Delete Tipe by :id
//	@Tags			kepegawaian/kualifikasi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Tipe ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/kepegawaian/kualifikasi/tipes/{id} [delete]
func (h *KepegawaianKualifikasiHandler) DeleteTipe(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteTipe(c.Request().Context(),id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "Tipe tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
