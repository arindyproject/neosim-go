package handlers

import (
	"io"
	"net/http"

	"neosim_go/internal/modules/kepegawaian/pendidikan/dto"
	"neosim_go/internal/shared/binding"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"

	"github.com/labstack/echo/v5"
)

// ─── ListPendidikan ─────────────────────────────────────────────────────
//
//	@Summary		Get list of KepegawaianPendidikan
//	@Description	Get paginated list of KepegawaianPendidikan
//	@Tags			kepegawaian/pendidikan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.KepegawaianPendidikanResponse}
//	@Router			/kepegawaian/pendidikan [get]
func (h *KepegawaianPendidikanHandler) ListPendidikan(c *echo.Context) error {

	filter := dto.FilterKepegawaianPendidikanRequest{
		NamaInstitusi: c.QueryParam("nama_institusi"),
	}
	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListPendidikan(c.Request().Context(), page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetPendidikanByID ───────────────────────────────────────────────────
//
//	@Summary		Get KepegawaianPendidikan
//	@Description	Get KepegawaianPendidikan by :id
//	@Tags			kepegawaian/pendidikan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"KepegawaianPendidikan ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.KepegawaianPendidikanResponse}
//	@Router			/kepegawaian/pendidikan/{id} [get]
func (h *KepegawaianPendidikanHandler) GetPendidikanByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetPendidikanByID(c.Request().Context(), id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreatePendidikan ────────────────────────────────────────────────────
//
//	@Summary		Create KepegawaianPendidikan
//	@Description	Create New KepegawaianPendidikan
//	@Tags			kepegawaian/pendidikan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateKepegawaianPendidikanRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.KepegawaianPendidikanResponse}
//	@Router			/kepegawaian/pendidikan [post]
func (h *KepegawaianPendidikanHandler) CreatePendidikan(c *echo.Context) error {
	var req dto.CreateKepegawaianPendidikanRequest
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
	item, err := h.service.CreatePendidikan(c.Request().Context(), &req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdatePendidikan ────────────────────────────────────────────────────
//
//	@Summary		Update KepegawaianPendidikan
//	@Description	Update KepegawaianPendidikan by :id
//	@Tags			kepegawaian/pendidikan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"KepegawaianPendidikan ID"
//	@Param			body	body		dto.UpdateKepegawaianPendidikanRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.KepegawaianPendidikanResponse}
//	@Router			/kepegawaian/pendidikan/{id} [put]
func (h *KepegawaianPendidikanHandler) UpdatePendidikan(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	var req dto.UpdateKepegawaianPendidikanRequest
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
	item, err := h.service.UpdatePendidikan(c.Request().Context(), id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "KepegawaianPendidikan tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeletePendidikan ────────────────────────────────────────────────────
//
//	@Summary		Delete KepegawaianPendidikan
//	@Description	Delete KepegawaianPendidikan by :id
//	@Tags			kepegawaian/pendidikan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"KepegawaianPendidikan ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/kepegawaian/pendidikan/{id} [delete]
func (h *KepegawaianPendidikanHandler) DeletePendidikan(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeletePendidikan(c.Request().Context(), id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "KepegawaianPendidikan tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
