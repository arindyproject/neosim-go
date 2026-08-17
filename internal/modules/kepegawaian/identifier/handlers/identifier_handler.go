package handlers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"

	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/shared/binding"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"
)

// ─── ListIdentifier ──────────────────────────────────────────────────────────
//
//	@Summary        Get list of KepegawaianIdentifier
//	@Description    Get paginated list of KepegawaianIdentifier
//	@Tags           kepegawaian/identifier
//	@Accept         json
//	@Produce        json
//	@Security       BearerAuth
//	@Param          pegawai_id  query       int     false   "Filter by Pegawai ID"
//	@Param          tipe_id     query       int     false   "Filter by Tipe ID"
//	@Param          nilai       query       string  false   "Filter by Nilai / Nomor Identifier (partial match)"
//	@Param          is_primary  query       boolean false   "Filter by Is Primary Status"
//	@Param          is_aktif    query       boolean false   "Filter by Is Aktif Status"
//	@Param          is_expired  query       boolean false   "Filter identifier yang sudah expired"
//	@Param          page        query       int     false   "Page number"
//	@Param          page_size   query       int     false   "Page size"
//	@Success        200         {object}    response.MyGoResponse{data=[]dto.KepegawaianIdentifierResponse}
//	@Router         /kepegawaian/identifier [get]
func (h *KepegawaianIdentifierHandler) ListIdentifier(c *echo.Context) error {
	filter := dto.FilterKepegawaianIdentifierRequest{
		Nilai: c.QueryParam("nilai"),
	}

	if pegawaiIDStr := c.QueryParam("pegawai_id"); pegawaiIDStr != "" {
		if val, err := strconv.ParseInt(pegawaiIDStr, 10, 64); err == nil {
			filter.PegawaiID = &val
		}
	}

	if tipeIDStr := c.QueryParam("tipe_id"); tipeIDStr != "" {
		if val, err := strconv.ParseInt(tipeIDStr, 10, 64); err == nil {
			filter.TipeID = &val
		}
	}

	if isPrimaryStr := c.QueryParam("is_primary"); isPrimaryStr != "" {
		if val, err := strconv.ParseBool(isPrimaryStr); err == nil {
			filter.IsPrimary = &val
		}
	}

	if isAktifStr := c.QueryParam("is_aktif"); isAktifStr != "" {
		if val, err := strconv.ParseBool(isAktifStr); err == nil {
			filter.IsAktif = &val
		}
	}

	if isExpiredStr := c.QueryParam("is_expired"); isExpiredStr != "" {
		if val, err := strconv.ParseBool(isExpiredStr); err == nil {
			filter.IsExpired = &val
		}
	}

	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListIdentifier(c.Request().Context(), page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, he.NotFoundStatus(err, "Gagal mengambil data"), false, err.Error(), nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetIdentifierByID ────────────────────────────────────────────────────────
//
//	@Summary        Get KepegawaianIdentifier
//	@Description    Get KepegawaianIdentifier by :id
//	@Tags           kepegawaian/identifier
//	@Accept         json
//	@Produce        json
//	@Security       BearerAuth
//	@Param          id  path        int true    "KepegawaianIdentifier ID"
//	@Success        200 {object}    response.MyGoResponse{data=dto.KepegawaianIdentifierResponse}
//	@Router         /kepegawaian/identifier/{id} [get]
func (h *KepegawaianIdentifierHandler) GetIdentifierByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetIdentifierByID(c.Request().Context(), id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── ListByPegawai ────────────────────────────────────────────────────────────
//
//	@Summary        Daftar identifier milik satu pegawai
//	@Description    Menampilkan semua identifier yang dimiliki oleh pegawai tertentu
//	@Tags           kepegawaian/identifier
//	@Accept         json
//	@Produce        json
//	@Security       BearerAuth
//	@Param          pegawai_id  path        int true    "ID pegawai"
//	@Param          page        query       int     false   "Page number"
//	@Param          page_size   query       int     false   "Page size"
//	@Success        200         {object}    response.MyGoResponse{data=[]dto.KepegawaianIdentifierResponse}
//	@Router         /kepegawaian/identifier/{pegawai_id}/pegawai [get]
func (h *KepegawaianIdentifierHandler) ListByPegawai(c *echo.Context) error {
	page, pageSize := he.ParsePagination(c, h.cfg)
	actor := he.BuildAuthContext(c)

	pegawaiID, err := parsePegawaiID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	items, total, err := h.service.ListByPegawai(c.Request().Context(), pegawaiID, page, pageSize, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	if len(items) == 0 {
		return response.Response(c, http.StatusNotFound, false, "Data tidak ditemukan", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── CreateIdentifier ──────────────────────────────────────────────────────────
//
//	@Summary        Create KepegawaianIdentifier
//	@Description    Create New KepegawaianIdentifier
//	@Tags           kepegawaian/identifier
//	@Accept         json
//	@Produce        json
//	@Security       BearerAuth
//	@Param          body    body        dto.CreateKepegawaianIdentifierRequest  true    "Create Request"
//	@Success        201     {object}    response.MyGoResponse{data=dto.KepegawaianIdentifierResponse}
//	@Router         /kepegawaian/identifier [post]
func (h *KepegawaianIdentifierHandler) CreateIdentifier(c *echo.Context) error {
	var req dto.CreateKepegawaianIdentifierRequest

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
	item, err := h.service.CreateIdentifier(c.Request().Context(), &req, actor)
	if err != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdateIdentifier ──────────────────────────────────────────────────────────
//
//	@Summary        Update KepegawaianIdentifier
//	@Description    Update KepegawaianIdentifier by :id
//	@Tags           kepegawaian/identifier
//	@Accept         json
//	@Produce        json
//	@Security       BearerAuth
//	@Param          id      path        int                                     true    "KepegawaianIdentifier ID"
//	@Param          body    body        dto.UpdateKepegawaianIdentifierRequest  true    "Update Request"
//	@Success        200     {object}    response.MyGoResponse{data=dto.KepegawaianIdentifierResponse}
//	@Router         /kepegawaian/identifier/{id} [put]
func (h *KepegawaianIdentifierHandler) UpdateIdentifier(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, err.Error())
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Gagal membaca request body", nil, err.Error())
	}

	var req dto.UpdateKepegawaianIdentifierRequest
	if errs := binding.BindErrors(body, &req); len(errs) > 0 {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.UpdateIdentifier(c.Request().Context(), id, &req, actor)
	if err != nil {
		return response.Response(c, he.NotFoundStatus(err, "Data tidak ditemukan"), false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeleteIdentifier ──────────────────────────────────────────────────────────
//
//	@Summary        Delete KepegawaianIdentifier
//	@Description    Delete KepegawaianIdentifier by :id (soft delete)
//	@Tags           kepegawaian/identifier
//	@Accept         json
//	@Produce        json
//	@Security       BearerAuth
//	@Param          id  path        int true    "KepegawaianIdentifier ID"
//	@Success        200 {object}    response.MyGoResponse{}
//	@Router         /kepegawaian/identifier/{id} [delete]
func (h *KepegawaianIdentifierHandler) DeleteIdentifier(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteIdentifier(c.Request().Context(), id, actor); err != nil {
		return response.Response(c, he.NotFoundStatus(err, "Data tidak ditemukan"), false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}

// ─── GetExpiringSoonIdentifier ──────────────────────────────────────────────────

// GetExpiringSoonIdentifier godoc
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
func (h *KepegawaianIdentifierHandler) GetExpiringSoonIdentifier(c *echo.Context) error {
	actor := he.BuildAuthContext(c)

	days := 30
	if d := c.QueryParam("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	items, err := h.service.GetExpiringSoonIdentifier(c.Request().Context(), days, actor)
	if err != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, err.Error(), nil, nil)
	}

	if len(items) == 0 {
		return response.Response(c, http.StatusOK, true, "Data tidak ditemukan", nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diambil", items, nil)
}

// ─── GetExpiredIdentifier ────────────────────────────────────────────────────

// GetExpiredIdentifier godoc
//
//	@Summary		Identifier yang sudah expired
//	@Description	Menampilkan daftar identifier (STR/SIP) yang sudah melewati tanggal expired
//	@Tags			kepegawaian/identifier
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.MyGoResponse{data=[]dto.KepegawaianIdentifierResponse}
//	@Router			/kepegawaian/identifier/expired [get]
func (h *KepegawaianIdentifierHandler) GetExpiredIdentifier(c *echo.Context) error {
	actor := he.BuildAuthContext(c)

	items, err := h.service.GetExpiredIdentifier(c.Request().Context(), actor)
	if err != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, err.Error(), nil, nil)
	}

	if len(items) == 0 {
		return response.Response(c, http.StatusOK, true, "Data tidak ditemukan", nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diambil", items, nil)
}

// ─── Helper privat ──────────────────────────────────────────────────────────

func parsePegawaiID(c *echo.Context) (int64, error) {
	return strconv.ParseInt(c.Param("pegawai_id"), 10, 64)
}
