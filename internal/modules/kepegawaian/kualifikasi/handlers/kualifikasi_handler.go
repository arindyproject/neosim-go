package handlers

import (
	"io"
	"net/http"

	"neosim_go/internal/modules/kepegawaian/kualifikasi/dto"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"
	"neosim_go/internal/shared/binding"
	he "neosim_go/internal/shared/httputil"

	"github.com/labstack/echo/v5"
)

// ─── ListKualifikasi ─────────────────────────────────────────────────────
//
//	@Summary		Get list of KepegawaianKualifikasi
//	@Description	Get paginated list of KepegawaianKualifikasi
//	@Tags			kepegawaian/kualifikasi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.KepegawaianKualifikasiResponse}
//	@Router			/kepegawaian/kualifikasi [get]
func (h *KepegawaianKualifikasiHandler) ListKualifikasi(c *echo.Context) error {

	filter := dto.FilterKepegawaianKualifikasiRequest{
		Name: c.QueryParam("name"),
	}
	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListKualifikasi(c.Request().Context(),page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetKualifikasiByID ───────────────────────────────────────────────────
//
//	@Summary		Get KepegawaianKualifikasi
//	@Description	Get KepegawaianKualifikasi by :id
//	@Tags			kepegawaian/kualifikasi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"KepegawaianKualifikasi ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.KepegawaianKualifikasiResponse}
//	@Router			/kepegawaian/kualifikasi/{id} [get]
func (h *KepegawaianKualifikasiHandler) GetKualifikasiByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetKualifikasiByID(c.Request().Context(),id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreateKualifikasi ────────────────────────────────────────────────────
//
//	@Summary		Create KepegawaianKualifikasi
//	@Description	Create New KepegawaianKualifikasi
//	@Tags			kepegawaian/kualifikasi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateKepegawaianKualifikasiRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.KepegawaianKualifikasiResponse}
//	@Router			/kepegawaian/kualifikasi [post]
func (h *KepegawaianKualifikasiHandler) CreateKualifikasi(c *echo.Context) error {
	var req dto.CreateKepegawaianKualifikasiRequest
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
	item, err := h.service.CreateKualifikasi(c.Request().Context(),&req,  actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdateKualifikasi ────────────────────────────────────────────────────
//
//	@Summary		Update KepegawaianKualifikasi
//	@Description	Update KepegawaianKualifikasi by :id
//	@Tags			kepegawaian/kualifikasi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"KepegawaianKualifikasi ID"
//	@Param			body	body		dto.UpdateKepegawaianKualifikasiRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.KepegawaianKualifikasiResponse}
//	@Router			/kepegawaian/kualifikasi/{id} [put]
func (h *KepegawaianKualifikasiHandler) UpdateKualifikasi(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	var req dto.UpdateKepegawaianKualifikasiRequest
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
	item, err := h.service.UpdateKualifikasi(c.Request().Context(),id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "KepegawaianKualifikasi tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeleteKualifikasi ────────────────────────────────────────────────────
//
//	@Summary		Delete KepegawaianKualifikasi
//	@Description	Delete KepegawaianKualifikasi by :id
//	@Tags			kepegawaian/kualifikasi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"KepegawaianKualifikasi ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/kepegawaian/kualifikasi/{id} [delete]
func (h *KepegawaianKualifikasiHandler) DeleteKualifikasi(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteKualifikasi(c.Request().Context(),id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "KepegawaianKualifikasi tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
