package handlers

import (
	"io"
	"net/http"
	"strconv"

	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/shared/binding"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"

	"github.com/labstack/echo/v5"
)

// ─── ListIdentifier ───────────────────────────────────────────────────── v
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

	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListIdentifier(page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetIdentifierByID ─────────────────────────────────────────────────── v
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
	item, err := h.service.GetIdentifierByID(id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreateIdentifier ──────────────────────────────────────────────────── v
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
	item, err := h.service.CreateIdentifier(&req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdateIdentifier ────────────────────────────────────────────────────
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
	item, err := h.service.UpdateIdentifier(id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "KepegawaianIdentifier tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeleteIdentifier ────────────────────────────────────────────────────
//
//	@Summary        Delete KepegawaianIdentifier
//	@Description    Delete KepegawaianIdentifier by :id
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
	if err := h.service.DeleteIdentifier(id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "KepegawaianIdentifier tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
