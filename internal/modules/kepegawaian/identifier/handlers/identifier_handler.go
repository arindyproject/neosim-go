package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"

	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/modules/kepegawaian/identifier/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"
)

// ─── List ──────────────────────────────────────────────────────────────────────

// List godoc
//
//	@Summary		Daftar identifier pegawai
//	@Description	Menampilkan daftar identifier pegawai dengan filter dan pagination
//	@Tags			kepegawaian/identifier
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			kepegawaian_id	query		int		false	"Filter by ID pegawai"
//	@Param			tipe			query		string	false	"Filter by tipe (NIK, STR, SIP, dll)"
//	@Param			is_aktif		query		bool	false	"Filter by status aktif"
//	@Param			is_expired		query		bool	false	"Filter identifier yang sudah expired"
//	@Param			is_primary		query		bool	false	"Filter identifier primary"
//	@Param			page			query		int		false	"Halaman (default: 1)"
//	@Param			page_size		query		int		false	"Jumlah per halaman (default: 10, max: 100)"
//	@Success		200				{object}	response.MyGoResponse{data=[]dto.KepegawaianIdentifierResponse}
//	@Router			/kepegawaian/identifier [get]
func (h *KepegawaianIdentifierHandler) List(c *echo.Context) error {
	actor := he.BuildAuthContext(c)

	page, pageSize := he.ParsePagination(c)
	filter := parseIdentifierFilter(c)
	items, total, err := h.service.List(c.Request().Context(), page, pageSize, filter, actor)
	if err != nil {
		return response.Response(c, he.NotFoundStatus(err, "Data tidak ditemukan"), false, err.Error(), nil, nil)
	}

	if len(items) == 0 {
		return response.Response(c, http.StatusNotFound, false, "Data tidak ditemukan", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Data berhasil diambil", items, total, page, pageSize)
}

// ─── GetByID ──────────────────────────────────────────────────────────────────

// GetByID godoc
//
//	@Summary		Detail identifier pegawai
//	@Description	Menampilkan detail satu identifier pegawai berdasarkan ID
//	@Tags			kepegawaian/identifier
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"ID identifier"
//	@Success		200	{object}	response.MyGoResponse{data=dto.KepegawaianIdentifierResponse}
//	@Router			/kepegawaian/identifier/{id} [get]
func (h *KepegawaianIdentifierHandler) GetByID(c *echo.Context) error {
	actor := he.BuildAuthContext(c)

	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	item, err := h.service.GetByID(c.Request().Context(), id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diambil", item, nil)
}

// ─── ListByPegawai ────────────────────────────────────────────────────────────

// ListByPegawai godoc
//
//	@Summary		Daftar identifier milik satu pegawai
//	@Description	Menampilkan semua identifier yang dimiliki oleh pegawai tertentu
//	@Tags			kepegawaian/identifier
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			pegawai_id	path		int	true	"ID pegawai"
//	@Success		200				{object}	response.MyGoResponse{data=[]dto.KepegawaianIdentifierResponse}
//	@Router			/kepegawaian/identifier/{pegawai_id}/pegawai [get]
func (h *KepegawaianIdentifierHandler) ListByPegawai(c *echo.Context) error {
	actor := he.BuildAuthContext(c)

	kepegawaianID, err := parseKepegawaianID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	items, err := h.service.ListByPegawai(c.Request().Context(), kepegawaianID, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	if len(items) == 0 {
		return response.Response(c, http.StatusNotFound, false, "Data tidak ditemukan", nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diambil", items, nil)
}

// ─── Create ───────────────────────────────────────────────────────────────────

// Create godoc
//
//	@Summary		Tambah identifier pegawai
//	@Description	Menambahkan identifier baru untuk pegawai (NIK, STR, SIP, dll)
//	@Tags			kepegawaian/identifier
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateKepegawaianIdentifierRequest	true	"Data identifier"
//	@Success		201		{object}	response.MyGoResponse{data=dto.KepegawaianIdentifierResponse}
//	@Router			/kepegawaian/identifier [post]
func (h *KepegawaianIdentifierHandler) Create(c *echo.Context) error {
	actor := he.BuildAuthContext(c)

	var req dto.CreateKepegawaianIdentifierRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	item, err := h.service.Create(c.Request().Context(), &req, actor)
	if err != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── Update ───────────────────────────────────────────────────────────────────

// Update godoc
//
//	@Summary		Update identifier pegawai
//	@Description	Memperbarui data identifier pegawai berdasarkan ID
//	@Tags			kepegawaian/identifier
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int										true	"ID identifier"
//	@Param			body	body		dto.UpdateKepegawaianIdentifierRequest	true	"Data yang diperbarui"
//	@Success		200		{object}	response.MyGoResponse{data=dto.KepegawaianIdentifierResponse}
//	@Router			/kepegawaian/identifier/{id} [put]
func (h *KepegawaianIdentifierHandler) Update(c *echo.Context) error {
	actor := he.BuildAuthContext(c)

	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	var req dto.UpdateKepegawaianIdentifierRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	item, err := h.service.Update(c.Request().Context(), id, &req, actor)
	if err != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diperbarui", item, nil)
}

// ─── Delete ───────────────────────────────────────────────────────────────────

// Delete godoc
//
//	@Summary		Hapus identifier pegawai
//	@Description	Menghapus identifier pegawai berdasarkan ID (soft delete)
//	@Tags			kepegawaian/identifier
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"ID identifier"
//	@Success		200	{object}	response.MyGoResponse
//	@Router			/kepegawaian/identifier/{id} [delete]
func (h *KepegawaianIdentifierHandler) Delete(c *echo.Context) error {
	actor := he.BuildAuthContext(c)

	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	if err := h.service.Delete(c.Request().Context(), id, actor); err != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}

// ─── GetExpiringSoon ──────────────────────────────────────────────────────────

// GetExpiringSoon godoc
//
//	@Summary		Identifier yang akan segera expired
//	@Description	Menampilkan daftar identifier (STR/SIP) yang akan expired dalam N hari ke depan
//	@Tags			kepegawaian/identifier
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			days	query		int	false	"Jumlah hari ke depan (default: 30)"
//	@Success		200		{object}	response.MyGoResponse{data=[]dto.KepegawaianIdentifierResponse}
//	@Router			/kepegawaian/identifier/expiring-soon [get]
func (h *KepegawaianIdentifierHandler) GetExpiringSoon(c *echo.Context) error {
	actor := he.BuildAuthContext(c)

	days := 30
	if d := c.QueryParam("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	items, err := h.service.GetExpiringSoon(c.Request().Context(), days, actor)
	if err != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, err.Error(), nil, nil)
	}

	if len(items) == 0 {
		return response.Response(c, http.StatusOK, true, "Data tidak ditemukan", nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diambil", items, nil)
}

// ─── GetIdentifierTypes ───────────────────────────────────────────────────────

// GetIdentifierTypes godoc
//
//	@Summary		Daftar tipe identifier
//	@Description	Menampilkan semua tipe identifier yang tersedia untuk dropdown UI
//	@Tags			kepegawaian/identifier
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.MyGoResponse{data=[]dto.IdentifierMetaResponse}
//	@Failure		401	{object}	response.MyGoResponse
//	@Router			/kepegawaian/identifier/types [get]
func (h *KepegawaianIdentifierHandler) GetIdentifierTypes(c *echo.Context) error {

	return response.Response(c, http.StatusOK, true, "Data berhasil diambil", h.service.GetIdentifierTypes(), nil)
}

// ─── Helper privat ────────────────────────────────────────────────────────────

func parseKepegawaianID(c *echo.Context) (int64, error) {
	return strconv.ParseInt(c.Param("pegawai_id"), 10, 64)
}

func parseIdentifierFilter(c *echo.Context) dto.FilterKepegawaianIdentifierRequest {
	filter := dto.FilterKepegawaianIdentifierRequest{}

	if v := c.QueryParam("pegawai_id"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			filter.PegawaiID = &parsed
		}
	}
	if v := c.QueryParam("tipe"); v != "" {
		tipe := models.IdentifierType(v)
		filter.Tipe = &tipe
	}
	if v := c.QueryParam("is_aktif"); v != "" {
		b := v == "true"
		filter.IsAktif = &b
	}
	if v := c.QueryParam("is_expired"); v != "" {
		b := v == "true"
		filter.IsExpired = &b
	}
	if v := c.QueryParam("is_primary"); v != "" {
		b := v == "true"
		filter.IsPrimary = &b
	}

	return filter
}
