package handlers

import (
	"io"
	"net/http"

	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"
	"neosim_go/internal/shared/binding"
	he "neosim_go/internal/shared/httputil"

	"github.com/labstack/echo/v5"
)

// ─── ListJabatan ─────────────────────────────────────────────────────
//
//	@Summary		Get list of KepegawaianJabatan
//	@Description	Get paginated list of KepegawaianJabatan
//	@Tags			kepegawaian/jabatan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.KepegawaianJabatanResponse}
//	@Router			/kepegawaian/jabatan [get]
func (h *KepegawaianJabatanHandler) ListJabatan(c *echo.Context) error {

	filter := dto.FilterKepegawaianJabatanRequest{
		Name: c.QueryParam("name"),
	}
	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListJabatan(c.Request().Context(),page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, err.Error(), nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetJabatanByID ───────────────────────────────────────────────────
//
//	@Summary		Get KepegawaianJabatan
//	@Description	Get KepegawaianJabatan by :id
//	@Tags			kepegawaian/jabatan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"KepegawaianJabatan ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.KepegawaianJabatanResponse}
//	@Router			/kepegawaian/jabatan/{id} [get]
func (h *KepegawaianJabatanHandler) GetJabatanByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetJabatanByID(c.Request().Context(),id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreateJabatan ────────────────────────────────────────────────────
//
//	@Summary		Create KepegawaianJabatan
//	@Description	Create New KepegawaianJabatan
//	@Tags			kepegawaian/jabatan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateKepegawaianJabatanRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.KepegawaianJabatanResponse}
//	@Router			/kepegawaian/jabatan [post]
func (h *KepegawaianJabatanHandler) CreateJabatan(c *echo.Context) error {
	var req dto.CreateKepegawaianJabatanRequest
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
	item, err := h.service.CreateJabatan(c.Request().Context(),&req,  actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdateJabatan ────────────────────────────────────────────────────
//
//	@Summary		Update KepegawaianJabatan
//	@Description	Update KepegawaianJabatan by :id
//	@Tags			kepegawaian/jabatan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"KepegawaianJabatan ID"
//	@Param			body	body		dto.UpdateKepegawaianJabatanRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.KepegawaianJabatanResponse}
//	@Router			/kepegawaian/jabatan/{id} [put]
func (h *KepegawaianJabatanHandler) UpdateJabatan(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	var req dto.UpdateKepegawaianJabatanRequest
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
	item, err := h.service.UpdateJabatan(c.Request().Context(),id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "KepegawaianJabatan tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeleteJabatan ────────────────────────────────────────────────────
//
//	@Summary		Delete KepegawaianJabatan
//	@Description	Delete KepegawaianJabatan by :id
//	@Tags			kepegawaian/jabatan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"KepegawaianJabatan ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/kepegawaian/jabatan/{id} [delete]
func (h *KepegawaianJabatanHandler) DeleteJabatan(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteJabatan(c.Request().Context(),id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "KepegawaianJabatan tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
