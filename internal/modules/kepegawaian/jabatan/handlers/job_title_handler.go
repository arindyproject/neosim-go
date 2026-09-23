package handlers

import (
	"io"
	"net/http"
	"strconv"

	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/shared/binding"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"

	"github.com/labstack/echo/v5"
)

// Method di bawah ini ditempelkan ke struct KepegawaianJabatanHandler yang
// sama dengan handler entitas utama (lihat handlers/handler.go). Nama method
// diberi suffix JobTitle agar tidak bentrok dengan method entitas utama
// pada struct handler yang sama.

// ─── ListSelectJobTitle ───────────────────────────────────────────────
//
//	@Summary		Get list of JobTitle for select
//	@Description	Get list of JobTitle for select with optional search
//	@Tags			kepegawaian/jabatan/job_titles
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			search	query	string	false	"Search by label or code"
//	@Success		200		{object}	response.MyGoResponse{data=[]dto.JobTitleSimpelResponse}
//	@Router			/kepegawaian/jabatan/job_titles/select [get]
func (h *KepegawaianJabatanHandler) ListSelectJobTitle(c *echo.Context) error {
	search := c.QueryParam("search")

	actor := he.BuildAuthContext(c)
	items, err := h.service.ListSelectJobTitle(c.Request().Context(), search, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data : "+err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", items, nil)
}

// ─── ListJobTitle ──────────────────────────────────────────────────────
//
//	@Summary		Get list of JobTitle
//	@Description	Get paginated list of JobTitle
//	@Tags			kepegawaian/jabatan/job_titles
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			label				query		string	false	"Filter by label (partial match)"
//	@Param			code				query		string	false	"Filter by code (partial match)"
//	@Param			kategori_id			query		int		false	"Filter by kategori ID"
//	@Param			rumpun_profesi_id	query		int		false	"Filter by rumpun profesi ID"
//	@Param			memerlukan_str		query		bool	false	"Filter by wajib STR"
//	@Param			memerlukan_sip		query		bool	false	"Filter by wajib SIP"
//	@Param			is_aktif			query		bool	false	"Filter by status aktif"
//	@Param			page				query		int		false	"Page number"
//	@Param			page_size			query		int		false	"Page size"
//	@Success		200					{object}	response.MyGoResponse{data=[]dto.JobTitleResponse}
//	@Router			/kepegawaian/jabatan/job_titles [get]
func (h *KepegawaianJabatanHandler) ListJobTitle(c *echo.Context) error {
	filter := dto.FilterJobTitleRequest{
		Label: c.QueryParam("label"),
		Code:  c.QueryParam("code"),
	}

	if v := c.QueryParam("kategori_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			filter.KategoriID = &id
		}
	}
	if v := c.QueryParam("rumpun_profesi_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			filter.RumpunProfesiID = &id
		}
	}
	if v := c.QueryParam("memerlukan_str"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			filter.MemerlukanSTR = &b
		}
	}
	if v := c.QueryParam("memerlukan_sip"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			filter.MemerlukanSIP = &b
		}
	}
	if v := c.QueryParam("is_aktif"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			filter.IsAktif = &b
		}
	}

	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListJobTitle(c.Request().Context(), page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, err.Error(), nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetJobTitleByID ───────────────────────────────────────────────────
//
//	@Summary		Get JobTitle
//	@Description	Get JobTitle by :id
//	@Tags			kepegawaian/jabatan/job_titles
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"JobTitle ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.JobTitleResponse}
//	@Router			/kepegawaian/jabatan/job_titles/{id} [get]
func (h *KepegawaianJabatanHandler) GetJobTitleByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetJobTitleByID(c.Request().Context(), id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreateJobTitle ────────────────────────────────────────────────────
//
//	@Summary		Create JobTitle
//	@Description	Create New JobTitle
//	@Tags			kepegawaian/jabatan/job_titles
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateJobTitleRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.JobTitleResponse}
//	@Router			/kepegawaian/jabatan/job_titles [post]
func (h *KepegawaianJabatanHandler) CreateJobTitle(c *echo.Context) error {
	var req dto.CreateJobTitleRequest
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
	item, err := h.service.CreateJobTitle(c.Request().Context(), &req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdateJobTitle ────────────────────────────────────────────────────
//
//	@Summary		Update JobTitle
//	@Description	Update JobTitle by :id
//	@Tags			kepegawaian/jabatan/job_titles
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"JobTitle ID"
//	@Param			body	body		dto.UpdateJobTitleRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.JobTitleResponse}
//	@Router			/kepegawaian/jabatan/job_titles/{id} [put]
func (h *KepegawaianJabatanHandler) UpdateJobTitle(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateJobTitleRequest
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
	item, err := h.service.UpdateJobTitle(c.Request().Context(), id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "JobTitle tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeleteJobTitle ────────────────────────────────────────────────────
//
//	@Summary		Delete JobTitle
//	@Description	Delete JobTitle by :id
//	@Tags			kepegawaian/jabatan/job_titles
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"JobTitle ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/kepegawaian/jabatan/job_titles/{id} [delete]
func (h *KepegawaianJabatanHandler) DeleteJobTitle(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteJobTitle(c.Request().Context(), id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "JobTitle tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
