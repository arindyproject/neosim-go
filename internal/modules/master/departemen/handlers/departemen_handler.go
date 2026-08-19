package handlers

import (
	"io"
	"net/http"

	"neosim_go/internal/modules/master/departemen/dto"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"
	"neosim_go/internal/shared/binding"
	he "neosim_go/internal/shared/httputil"

	"github.com/labstack/echo/v5"
)

// ─── ListDepartemen ─────────────────────────────────────────────────────
//
//	@Summary		Get list of MasterDepartemen
//	@Description	Get paginated list of MasterDepartemen
//	@Tags			master/departemen
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.MasterDepartemenResponse}
//	@Router			/master/departemen [get]
func (h *MasterDepartemenHandler) ListDepartemen(c *echo.Context) error {

	filter := dto.FilterMasterDepartemenRequest{
		Name: c.QueryParam("name"),
	}
	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListDepartemen(c.Request().Context(),page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetDepartemenByID ───────────────────────────────────────────────────
//
//	@Summary		Get MasterDepartemen
//	@Description	Get MasterDepartemen by :id
//	@Tags			master/departemen
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"MasterDepartemen ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.MasterDepartemenResponse}
//	@Router			/master/departemen/{id} [get]
func (h *MasterDepartemenHandler) GetDepartemenByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetDepartemenByID(c.Request().Context(),id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreateDepartemen ────────────────────────────────────────────────────
//
//	@Summary		Create MasterDepartemen
//	@Description	Create New MasterDepartemen
//	@Tags			master/departemen
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateMasterDepartemenRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.MasterDepartemenResponse}
//	@Router			/master/departemen [post]
func (h *MasterDepartemenHandler) CreateDepartemen(c *echo.Context) error {
	var req dto.CreateMasterDepartemenRequest
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
	item, err := h.service.CreateDepartemen(c.Request().Context(),&req,  actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdateDepartemen ────────────────────────────────────────────────────
//
//	@Summary		Update MasterDepartemen
//	@Description	Update MasterDepartemen by :id
//	@Tags			master/departemen
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"MasterDepartemen ID"
//	@Param			body	body		dto.UpdateMasterDepartemenRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.MasterDepartemenResponse}
//	@Router			/master/departemen/{id} [put]
func (h *MasterDepartemenHandler) UpdateDepartemen(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	var req dto.UpdateMasterDepartemenRequest
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
	item, err := h.service.UpdateDepartemen(c.Request().Context(),id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "MasterDepartemen tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeleteDepartemen ────────────────────────────────────────────────────
//
//	@Summary		Delete MasterDepartemen
//	@Description	Delete MasterDepartemen by :id
//	@Tags			master/departemen
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"MasterDepartemen ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/master/departemen/{id} [delete]
func (h *MasterDepartemenHandler) DeleteDepartemen(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteDepartemen(c.Request().Context(),id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "MasterDepartemen tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
