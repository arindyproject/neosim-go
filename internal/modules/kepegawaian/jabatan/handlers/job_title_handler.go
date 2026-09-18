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
// diberi suffix JobTitle agar tidak bentrok dengan method entitas utama
// pada struct handler yang sama.

// ─── ListJobTitle ──────────────────────────────────────────────────────
//
//	@Summary		Get list of JobTitle
//	@Description	Get paginated list of JobTitle
//	@Tags			kepegawaian/jabatan/job_titles
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.JobTitleResponse}
//	@Router			/kepegawaian/jabatan/job_titles [get]
func (h *KepegawaianJabatanHandler) ListJobTitle(c *echo.Context) error {
	filter := dto.FilterJobTitleRequest{Name: c.QueryParam("name")}
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
