package handlers

import (
	"io"
	"net/http"

	"neosim_go/internal/modules/kepegawaian/alamat/dto"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"
	"neosim_go/internal/shared/binding"
	he "neosim_go/internal/shared/httputil"

	"github.com/labstack/echo/v5"
)

// ─── ListAlamat ─────────────────────────────────────────────────────
//
//	@Summary		Get list of KepegawaianAlamat
//	@Description	Get paginated list of KepegawaianAlamat
//	@Tags			kepegawaian/alamat
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.KepegawaianAlamatResponse}
//	@Router			/kepegawaian/alamat [get]
func (h *KepegawaianAlamatHandler) ListAlamat(c *echo.Context) error {

	filter := dto.FilterKepegawaianAlamatRequest{
		Name: c.QueryParam("name"),
	}
	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListAlamat(c.Request().Context(),page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetAlamatByID ───────────────────────────────────────────────────
//
//	@Summary		Get KepegawaianAlamat
//	@Description	Get KepegawaianAlamat by :id
//	@Tags			kepegawaian/alamat
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"KepegawaianAlamat ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.KepegawaianAlamatResponse}
//	@Router			/kepegawaian/alamat/{id} [get]
func (h *KepegawaianAlamatHandler) GetAlamatByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetAlamatByID(c.Request().Context(),id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreateAlamat ────────────────────────────────────────────────────
//
//	@Summary		Create KepegawaianAlamat
//	@Description	Create New KepegawaianAlamat
//	@Tags			kepegawaian/alamat
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateKepegawaianAlamatRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.KepegawaianAlamatResponse}
//	@Router			/kepegawaian/alamat [post]
func (h *KepegawaianAlamatHandler) CreateAlamat(c *echo.Context) error {
	var req dto.CreateKepegawaianAlamatRequest
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
	item, err := h.service.CreateAlamat(c.Request().Context(),&req,  actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdateAlamat ────────────────────────────────────────────────────
//
//	@Summary		Update KepegawaianAlamat
//	@Description	Update KepegawaianAlamat by :id
//	@Tags			kepegawaian/alamat
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"KepegawaianAlamat ID"
//	@Param			body	body		dto.UpdateKepegawaianAlamatRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.KepegawaianAlamatResponse}
//	@Router			/kepegawaian/alamat/{id} [put]
func (h *KepegawaianAlamatHandler) UpdateAlamat(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	var req dto.UpdateKepegawaianAlamatRequest
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
	item, err := h.service.UpdateAlamat(c.Request().Context(),id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "KepegawaianAlamat tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeleteAlamat ────────────────────────────────────────────────────
//
//	@Summary		Delete KepegawaianAlamat
//	@Description	Delete KepegawaianAlamat by :id
//	@Tags			kepegawaian/alamat
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"KepegawaianAlamat ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/kepegawaian/alamat/{id} [delete]
func (h *KepegawaianAlamatHandler) DeleteAlamat(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteAlamat(c.Request().Context(),id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "KepegawaianAlamat tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
