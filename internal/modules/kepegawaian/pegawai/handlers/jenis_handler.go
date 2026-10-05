package handlers

import (
	"io"
	"net/http"

	"neosim_go/internal/modules/kepegawaian/pegawai/dto"
	"neosim_go/internal/shared/binding"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"

	"github.com/labstack/echo/v5"
)

// Method di bawah ini ditempelkan ke struct KepegawaianPegawaiHandler yang
// sama dengan handler entitas utama (lihat handlers/handler.go). Nama method
// diberi suffix Jenis agar tidak bentrok dengan method entitas utama
// pada struct handler yang sama.

// ─── ListSelectJenis ────────────────────────────────────────────────
//
//	@Summary		Get list of Jenis for select
//	@Description	Get list of Jenis for select with optional search
//	@Tags			kepegawaian/pegawai/jeniss
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			search	query	string	false	"Search by label or code"
//	@Success		200		{object}	response.MyGoResponse{data=[]dto.JenisSimpelResponse}
//	@Router			/kepegawaian/alamat/pegawai/jeniss [get]
func (h *KepegawaianPegawaiHandler) ListSelectJenis(c *echo.Context) error {
	search := c.QueryParam("search")

	actor := he.BuildAuthContext(c)
	items, err := h.service.ListSelectJenis(c.Request().Context(), search, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data : "+err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", items, nil)
}

// ─── ListJenis ──────────────────────────────────────────────────────
//
//	@Summary		Get list of Jenis
//	@Description	Get paginated list of Jenis
//	@Tags			kepegawaian/pegawai/jeniss
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			label		query		string	false	"Filter by label (partial match)"
//	@Param			code		query		string	false	"Filter by code (partial match)"
//	@Param			sort_by		query		string	false	"Sort column (code, label,fhir_code, created_at, updated_at)"	Enums(label, code,fhir_code, created_at, updated_at)
//	@Param			sort_order	query		string	false	"Sort direction"	Enums(asc, desc)
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.JenisResponse}
//	@Router			/kepegawaian/pegawai/jeniss [get]
func (h *KepegawaianPegawaiHandler) ListJenis(c *echo.Context) error {
	filter := dto.FilterJenisRequest{
		Code:  c.QueryParam("code"),
		Label: c.QueryParam("label"),
		// Sorting ---------------------------
		SortBy:    c.QueryParam("sort_by"),
		SortOrder: c.QueryParam("sort_order"),
	}
	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListJenis(c.Request().Context(), page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, err.Error(), nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetJenisByID ───────────────────────────────────────────────────
//
//	@Summary		Get Jenis
//	@Description	Get Jenis by :id
//	@Tags			kepegawaian/pegawai/jeniss
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Jenis ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.JenisResponse}
//	@Router			/kepegawaian/pegawai/jeniss/{id} [get]
func (h *KepegawaianPegawaiHandler) GetJenisByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetJenisByID(c.Request().Context(), id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreateJenis ────────────────────────────────────────────────────
//
//	@Summary		Create Jenis
//	@Description	Create New Jenis
//	@Tags			kepegawaian/pegawai/jeniss
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateJenisRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.JenisResponse}
//	@Router			/kepegawaian/pegawai/jeniss [post]
func (h *KepegawaianPegawaiHandler) CreateJenis(c *echo.Context) error {
	var req dto.CreateJenisRequest
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
	item, err := h.service.CreateJenis(c.Request().Context(), &req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdateJenis ────────────────────────────────────────────────────
//
//	@Summary		Update Jenis
//	@Description	Update Jenis by :id
//	@Tags			kepegawaian/pegawai/jeniss
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"Jenis ID"
//	@Param			body	body		dto.UpdateJenisRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.JenisResponse}
//	@Router			/kepegawaian/pegawai/jeniss/{id} [put]
func (h *KepegawaianPegawaiHandler) UpdateJenis(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateJenisRequest
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
	item, err := h.service.UpdateJenis(c.Request().Context(), id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "Jenis tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeleteJenis ────────────────────────────────────────────────────
//
//	@Summary		Delete Jenis
//	@Description	Delete Jenis by :id
//	@Tags			kepegawaian/pegawai/jeniss
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Jenis ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/kepegawaian/pegawai/jeniss/{id} [delete]
func (h *KepegawaianPegawaiHandler) DeleteJenis(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteJenis(c.Request().Context(), id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "Jenis tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
