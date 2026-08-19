package handlers

import (
	"io"
	"net/http"

	"neosim_go/internal/modules/kepegawaian/pegawai/dto"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"
	"neosim_go/internal/shared/binding"
	he "neosim_go/internal/shared/httputil"

	"github.com/labstack/echo/v5"
)

// ─── ListPegawai ─────────────────────────────────────────────────────
//
//	@Summary		Get list of KepegawaianPegawai
//	@Description	Get paginated list of KepegawaianPegawai
//	@Tags			kepegawaian/pegawai
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.KepegawaianPegawaiResponse}
//	@Router			/kepegawaian/pegawai [get]
func (h *KepegawaianPegawaiHandler) ListPegawai(c *echo.Context) error {

	filter := dto.FilterKepegawaianPegawaiRequest{
		Name: c.QueryParam("name"),
	}
	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListPegawai(c.Request().Context(),page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetPegawaiByID ───────────────────────────────────────────────────
//
//	@Summary		Get KepegawaianPegawai
//	@Description	Get KepegawaianPegawai by :id
//	@Tags			kepegawaian/pegawai
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"KepegawaianPegawai ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.KepegawaianPegawaiResponse}
//	@Router			/kepegawaian/pegawai/{id} [get]
func (h *KepegawaianPegawaiHandler) GetPegawaiByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetPegawaiByID(c.Request().Context(),id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreatePegawai ────────────────────────────────────────────────────
//
//	@Summary		Create KepegawaianPegawai
//	@Description	Create New KepegawaianPegawai
//	@Tags			kepegawaian/pegawai
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateKepegawaianPegawaiRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.KepegawaianPegawaiResponse}
//	@Router			/kepegawaian/pegawai [post]
func (h *KepegawaianPegawaiHandler) CreatePegawai(c *echo.Context) error {
	var req dto.CreateKepegawaianPegawaiRequest
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
	item, err := h.service.CreatePegawai(c.Request().Context(),&req,  actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdatePegawai ────────────────────────────────────────────────────
//
//	@Summary		Update KepegawaianPegawai
//	@Description	Update KepegawaianPegawai by :id
//	@Tags			kepegawaian/pegawai
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"KepegawaianPegawai ID"
//	@Param			body	body		dto.UpdateKepegawaianPegawaiRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.KepegawaianPegawaiResponse}
//	@Router			/kepegawaian/pegawai/{id} [put]
func (h *KepegawaianPegawaiHandler) UpdatePegawai(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	var req dto.UpdateKepegawaianPegawaiRequest
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
	item, err := h.service.UpdatePegawai(c.Request().Context(),id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "KepegawaianPegawai tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeletePegawai ────────────────────────────────────────────────────
//
//	@Summary		Delete KepegawaianPegawai
//	@Description	Delete KepegawaianPegawai by :id
//	@Tags			kepegawaian/pegawai
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"KepegawaianPegawai ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/kepegawaian/pegawai/{id} [delete]
func (h *KepegawaianPegawaiHandler) DeletePegawai(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeletePegawai(c.Request().Context(),id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "KepegawaianPegawai tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
