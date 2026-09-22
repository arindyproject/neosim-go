package handlers

import (
	"io"
	"net/http"

	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/shared/binding"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"

	"github.com/labstack/echo/v5"
)

// Method di bawah ini ditempelkan ke struct KepegawaianJabatanHandler yang
// sama dengan handler entitas utama (lihat handlers/handler.go). Nama method
// diberi suffix PositionKategori agar tidak bentrok dengan method entitas utama
// pada struct handler yang sama.

// ─── ListPositionKategori ──────────────────────────────────────────────────────
//
//	@Summary		Get list of PositionKategori
//	@Description	Get paginated list of PositionKategori
//	@Tags			kepegawaian/jabatan/position/kategoris
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.PositionKategoriResponse}
//	@Router			/kepegawaian/jabatan/position_kategoris [get]
func (h *KepegawaianJabatanHandler) ListPositionKategori(c *echo.Context) error {
	filter := dto.FilterPositionKategoriRequest{
		Search: c.QueryParam("search"),
		Code:   c.QueryParam("code"),
		Label:  c.QueryParam("label"),
	}
	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListPositionKategori(c.Request().Context(), page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, err.Error(), nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── ListSelectPositionKategori ────────────────────────────────────
//
//	@Summary		Get list of Tipe for select
//	@Description	Get list of Tipe for select with optional search
//	@Tags			kepegawaian/jabatan/position/kategoris
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			search	query	string	false	"Search by label or code"
//	@Success		200		{object}	response.MyGoResponse{data=[]dto.PositionKategoriSelectResponse}
//	@Router			/kepegawaian/jabatan/position_kategoris/select [get]
func (h *KepegawaianJabatanHandler) ListSelectPositionKategori(c *echo.Context) error {
	search := c.QueryParam("search")

	actor := he.BuildAuthContext(c)
	items, err := h.service.ListSelectPositionKategori(c.Request().Context(), search, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data : "+err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", items, nil)
}

// ─── GetPositionKategoriByID ───────────────────────────────────────────────────
//
//	@Summary		Get PositionKategori
//	@Description	Get PositionKategori by :id
//	@Tags			kepegawaian/jabatan/position/kategoris
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"PositionKategori ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.PositionKategoriResponse}
//	@Router			/kepegawaian/jabatan/position_kategoris/{id} [get]
func (h *KepegawaianJabatanHandler) GetPositionKategoriByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetPositionKategoriByID(c.Request().Context(), id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreatePositionKategori ────────────────────────────────────────────────────
//
//	@Summary		Create PositionKategori
//	@Description	Create New PositionKategori
//	@Tags			kepegawaian/jabatan/position/kategoris
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreatePositionKategoriRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.PositionKategoriResponse}
//	@Router			/kepegawaian/jabatan/position_kategoris [post]
func (h *KepegawaianJabatanHandler) CreatePositionKategori(c *echo.Context) error {
	var req dto.CreatePositionKategoriRequest
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
	item, err := h.service.CreatePositionKategori(c.Request().Context(), &req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdatePositionKategori ────────────────────────────────────────────────────
//
//	@Summary		Update PositionKategori
//	@Description	Update PositionKategori by :id
//	@Tags			kepegawaian/jabatan/position/kategoris
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"PositionKategori ID"
//	@Param			body	body		dto.UpdatePositionKategoriRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.PositionKategoriResponse}
//	@Router			/kepegawaian/jabatan/position_kategoris/{id} [put]
func (h *KepegawaianJabatanHandler) UpdatePositionKategori(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdatePositionKategoriRequest
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
	item, err := h.service.UpdatePositionKategori(c.Request().Context(), id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "PositionKategori tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeletePositionKategori ────────────────────────────────────────────────────
//
//	@Summary		Delete PositionKategori
//	@Description	Delete PositionKategori by :id
//	@Tags			kepegawaian/jabatan/position/kategoris
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"PositionKategori ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/kepegawaian/jabatan/position_kategoris/{id} [delete]
func (h *KepegawaianJabatanHandler) DeletePositionKategori(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeletePositionKategori(c.Request().Context(), id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "PositionKategori tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
