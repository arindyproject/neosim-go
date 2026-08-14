package handlers

import (
	"net/http"

	"neosim_go/internal/modules/kepegawaian/kontak/dto"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"
	he "neosim_go/internal/shared/httputil"

	"github.com/labstack/echo/v5"
)

// ─── ListKontak ─────────────────────────────────────────────────────
//
//	@Summary		Get list of KepegawaianKontak
//	@Description	Get paginated list of KepegawaianKontak
//	@Tags			kepegawaian/kontak
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.KepegawaianKontakResponse}
//	@Router			/kepegawaian/kontak [get]
func (h *KepegawaianKontakHandler) ListKontak(c *echo.Context) error {

	filter := dto.FilterKepegawaianKontakRequest{
		Name: c.QueryParam("name"),
	}
	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListKontak(page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetKontakByID ───────────────────────────────────────────────────
//
//	@Summary		Get KepegawaianKontak
//	@Description	Get KepegawaianKontak by :id
//	@Tags			kepegawaian/kontak
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"KepegawaianKontak ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.KepegawaianKontakResponse}
//	@Router			/kepegawaian/kontak/{id} [get]
func (h *KepegawaianKontakHandler) GetKontakByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetKontakByID(id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreateKontak ────────────────────────────────────────────────────
//
//	@Summary		Create KepegawaianKontak
//	@Description	Create New KepegawaianKontak
//	@Tags			kepegawaian/kontak
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateKepegawaianKontakRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.KepegawaianKontakResponse}
//	@Router			/kepegawaian/kontak [post]
func (h *KepegawaianKontakHandler) CreateKontak(c *echo.Context) error {
	var req dto.CreateKepegawaianKontakRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.CreateKontak(&req,  actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdateKontak ────────────────────────────────────────────────────
//
//	@Summary		Update KepegawaianKontak
//	@Description	Update KepegawaianKontak by :id
//	@Tags			kepegawaian/kontak
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"KepegawaianKontak ID"
//	@Param			body	body		dto.UpdateKepegawaianKontakRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.KepegawaianKontakResponse}
//	@Router			/kepegawaian/kontak/{id} [put]
func (h *KepegawaianKontakHandler) UpdateKontak(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	var req dto.UpdateKepegawaianKontakRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.UpdateKontak(id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "KepegawaianKontak tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeleteKontak ────────────────────────────────────────────────────
//
//	@Summary		Delete KepegawaianKontak
//	@Description	Delete KepegawaianKontak by :id
//	@Tags			kepegawaian/kontak
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"KepegawaianKontak ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/kepegawaian/kontak/{id} [delete]
func (h *KepegawaianKontakHandler) DeleteKontak(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteKontak(id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "KepegawaianKontak tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
