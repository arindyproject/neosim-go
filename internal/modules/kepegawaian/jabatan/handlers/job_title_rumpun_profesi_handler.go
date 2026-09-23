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
// diberi suffix JobTitleRumpunProfesi agar tidak bentrok dengan method entitas utama
// pada struct handler yang sama.

// ─── ListSelectJobTitleRumpunProfesi ────────────────────────────────────────────
//
//	@Summary		Get list of JobTitleRumpunProfesi for select
//	@Description	Get list of JobTitleRumpunProfesi for select with optional search
//	@Tags			kepegawaian/jabatan/jabatan/job_title_rumpun_profesi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			search	query	string	false	"Search by label or code"
//	@Success		200		{object}	response.MyGoResponse{data=[]dto.JobTitleRumpunProfesiSimpelResponse}
//	@Router			/kepegawaian/jabatan/job_title_rumpun_profesis/select [get]
func (h *KepegawaianJabatanHandler) ListSelectListJobTitleRumpunProfesi(c *echo.Context) error {
	search := c.QueryParam("search")

	actor := he.BuildAuthContext(c)
	items, err := h.service.ListSelectJobTitleKategori(c.Request().Context(), search, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data : "+err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", items, nil)
}

// ─── ListJobTitleRumpunProfesi ──────────────────────────────────────────────────────
//
//	@Summary		Get list of JobTitleRumpunProfesi
//	@Description	Get paginated list of JobTitleRumpunProfesi
//	@Tags			kepegawaian/jabatan/jabatan/job_title_rumpun_profesi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.JobTitleRumpunProfesiResponse}
//	@Router			/kepegawaian/jabatan/job_title_rumpun_profesis [get]
func (h *KepegawaianJabatanHandler) ListJobTitleRumpunProfesi(c *echo.Context) error {
	filter := dto.FilterJobTitleRumpunProfesiRequest{Code: c.QueryParam("code"), Label: c.QueryParam("label")}
	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListJobTitleRumpunProfesi(c.Request().Context(), page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, err.Error(), nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetJobTitleRumpunProfesiByID ───────────────────────────────────────────────────
//
//	@Summary		Get JobTitleRumpunProfesi
//	@Description	Get JobTitleRumpunProfesi by :id
//	@Tags			kepegawaian/jabatan/jabatan/job_title_rumpun_profesi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"JobTitleRumpunProfesi ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.JobTitleRumpunProfesiResponse}
//	@Router			/kepegawaian/jabatan/job_title_rumpun_profesis/{id} [get]
func (h *KepegawaianJabatanHandler) GetJobTitleRumpunProfesiByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetJobTitleRumpunProfesiByID(c.Request().Context(), id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreateJobTitleRumpunProfesi ────────────────────────────────────────────────────
//
//	@Summary		Create JobTitleRumpunProfesi
//	@Description	Create New JobTitleRumpunProfesi
//	@Tags			kepegawaian/jabatan/jabatan/job_title_rumpun_profesi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateJobTitleRumpunProfesiRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.JobTitleRumpunProfesiResponse}
//	@Router			/kepegawaian/jabatan/job_title_rumpun_profesis [post]
func (h *KepegawaianJabatanHandler) CreateJobTitleRumpunProfesi(c *echo.Context) error {
	var req dto.CreateJobTitleRumpunProfesiRequest
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
	item, err := h.service.CreateJobTitleRumpunProfesi(c.Request().Context(), &req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdateJobTitleRumpunProfesi ────────────────────────────────────────────────────
//
//	@Summary		Update JobTitleRumpunProfesi
//	@Description	Update JobTitleRumpunProfesi by :id
//	@Tags			kepegawaian/jabatan/jabatan/job_title_rumpun_profesi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"JobTitleRumpunProfesi ID"
//	@Param			body	body		dto.UpdateJobTitleRumpunProfesiRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.JobTitleRumpunProfesiResponse}
//	@Router			/kepegawaian/jabatan/job_title_rumpun_profesis/{id} [put]
func (h *KepegawaianJabatanHandler) UpdateJobTitleRumpunProfesi(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateJobTitleRumpunProfesiRequest
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
	item, err := h.service.UpdateJobTitleRumpunProfesi(c.Request().Context(), id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "JobTitleRumpunProfesi tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeleteJobTitleRumpunProfesi ────────────────────────────────────────────────────
//
//	@Summary		Delete JobTitleRumpunProfesi
//	@Description	Delete JobTitleRumpunProfesi by :id
//	@Tags			kepegawaian/jabatan/jabatan/job_title_rumpun_profesi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"JobTitleRumpunProfesi ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/kepegawaian/jabatan/job_title_rumpun_profesis/{id} [delete]
func (h *KepegawaianJabatanHandler) DeleteJobTitleRumpunProfesi(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteJobTitleRumpunProfesi(c.Request().Context(), id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "JobTitleRumpunProfesi tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
