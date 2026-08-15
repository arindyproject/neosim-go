package handlers

import (
	"io"
	"neosim_go/internal/modules/master/alamat/dto"
	"neosim_go/internal/shared/binding"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"
	"net/http"

	"github.com/labstack/echo/v5"
)

// Kecamatan ===========================================================================
// ─────────────── List ─────────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Get list of Kecamatan
//	@Description	Get paginated list of Kecamatan
//	@Tags			master/alamat/kecamatan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			kota_kabupaten_id	query		int		false	"Filter by kota_kabupaten_id"
//	@Param			code				query		string	false	"Filter by code"
//	@Param			name				query		string	false	"Filter by name (partial match)"
//	@Param			page				query		int		false	"Page number"
//	@Param			page_size			query		int		false	"Page size"
//	@Success		200					{object}	response.MyGoResponse{data=[]dto.KecamatanResponse}
//	@Router			/master/alamat/kecamatan [get]
//
// ListKecamatan handles GET /api/v1/master/alamat/kecamatan
func (h *MasterAlamatHandler) ListKecamatan(c *echo.Context) error {
	page, pageSize := he.ParsePagination(c, h.cfg)
	kotaKabupatenID := he.ParseOptionalInt64Query(c, "kota_kabupaten_id")
	filter := dto.FilterKecamatanRequest{
		Code: c.QueryParam("code"),
		Name: c.QueryParam("name"),
	}

	items, total, err := h.service.ListKecamatan(page, pageSize, kotaKabupatenID, &filter)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── GetByID ────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Get Kecamatan
//	@Description	Get Kecamatan detail by :id, termasuk statistik turunan
//	@Tags			master/alamat/kecamatan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Kecamatan ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.KecamatanDetailResponse}
//	@Router			/master/alamat/kecamatan/{id} [get]
//
// GetByIDKecamatan handles GET /api/v1/master/alamat/kecamatan/:id
func (h *MasterAlamatHandler) GetByIDKecamatan(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	item, err := h.service.GetByIDKecamatan(id)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Create Kecamatan
//	@Description	Create New Kecamatan
//	@Tags			master/alamat/kecamatan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateKecamatanRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.KecamatanResponse}
//	@Router			/master/alamat/kecamatan [post]
//
// CreateKecamatan handles POST /api/v1/master/alamat/kecamatan
func (h *MasterAlamatHandler) CreateKecamatan(c *echo.Context) error {
	var req dto.CreateKecamatanRequest
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Gagal membaca request body", nil, err.Error())
	}

	if errs := binding.BindErrors(body, &req); len(errs) > 0 {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.CreateKecamatan(&req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Update Kecamatan
//	@Description	Update Kecamatan by :id
//	@Tags			master/alamat/kecamatan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"Kecamatan ID"
//	@Param			body	body		dto.UpdateKecamatanRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.KecamatanResponse}
//	@Router			/master/alamat/kecamatan/{id} [put]
//
// UpdateKecamatan handles PUT /api/v1/master/alamat/kecamatan/:id
func (h *MasterAlamatHandler) UpdateKecamatan(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateKecamatanRequest
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Gagal membaca request body", nil, err.Error())
	}

	if errs := binding.BindErrors(body, &req); len(errs) > 0 {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.UpdateKecamatan(id, &req, actor)
	if err != nil {
		return response.Response(c, he.NotFoundStatus(err, "Kecamatan tidak ditemukan"), false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Delete Kecamatan
//	@Description	Delete Kecamatan by :id
//	@Tags			master/alamat/kecamatan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Kecamatan ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/master/alamat/kecamatan/{id} [delete]
//
// DeleteKecamatan handles DELETE /api/v1/master/alamat/kecamatan/:id
func (h *MasterAlamatHandler) DeleteKecamatan(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteKecamatan(id, actor); err != nil {
		return response.Response(c, he.NotFoundStatus(err, "Kecamatan tidak ditemukan"), false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
} // ───────────── Delete ─────────────────────────────────────────────────────────
// Kecamatan ===========================================================================
