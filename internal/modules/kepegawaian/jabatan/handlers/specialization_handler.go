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
// diberi suffix Specialization agar tidak bentrok dengan method entitas utama
// pada struct handler yang sama.

// ─── ListSpecialization ──────────────────────────────────────────────────────
//
//	@Summary		Get list of Specialization
//	@Description	Get paginated list of Specialization
//	@Tags			kepegawaian/jabatan/specializations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			code			query		string	false	"Filter by code (partial match)"
//	@Param			label			query		string	false	"Filter by label (partial match)"
//	@Param			job_title_id	query		int		false	"Filter by job title ID"
//	@Param			kategori_id		query		int		false	"Filter by kategori ID"
//	@Param			is_aktif		query		bool	false	"Filter by status aktif"
//	@Param			sort_by			query		string	false	"Sort column (code,label,job_title_id,kategori_id,gelar,lama_pendidikan_tahun,fhir_code, created_at, updated_at)"	Enums(code,label,job_title_id,kategori_id,gelar,lama_pendidikan_tahun,fhir_code, created_at, updated_at)
//	@Param			sort_order		query		string	false	"Sort direction"	Enums(asc, desc)
//	@Param			page			query		int		false	"Page number"
//	@Param			page_size		query		int		false	"Page size"
//	@Success		200				{object}	response.MyGoResponse{data=[]dto.SpecializationResponse}
//	@Router			/kepegawaian/jabatan/specializations [get]
func (h *KepegawaianJabatanHandler) ListSpecialization(c *echo.Context) error {
	filter := dto.FilterSpecializationRequest{
		Code:  c.QueryParam("code"),
		Label: c.QueryParam("label"),

		// Sorting ---------------------------
		SortBy:    c.QueryParam("sort_by"),
		SortOrder: c.QueryParam("sort_order"),
	}

	if v := c.QueryParam("job_title_id"); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n <= 0 {
			return response.Response(c, http.StatusBadRequest, false, "job_title_id tidak valid", nil, nil)
		}
		filter.JobTitleID = &n
	}
	if v := c.QueryParam("kategori_id"); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n <= 0 {
			return response.Response(c, http.StatusBadRequest, false, "kategori_id tidak valid", nil, nil)
		}
		filter.KategoriID = &n
	}
	if v := c.QueryParam("is_aktif"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return response.Response(c, http.StatusBadRequest, false, "is_aktif tidak valid", nil, nil)
		}
		filter.IsAktif = &b
	}

	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListSpecialization(c.Request().Context(), page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, err.Error(), nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── ListSelectSpecialization ────────────────────────────────────────────────
//
//	@Summary		Get select list of Specialization
//	@Description	Get list Specialization aktif untuk dropdown/select (tanpa pagination)
//	@Tags			kepegawaian/jabatan/specializations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			search	query		string	false	"Search by label, code, atau gelar"
//	@Success		200		{object}	response.MyGoResponse{data=[]dto.SpecializationSelectResponse}
//	@Router			/kepegawaian/jabatan/specializations/select [get]
func (h *KepegawaianJabatanHandler) ListSelectSpecialization(c *echo.Context) error {
	actor := he.BuildAuthContext(c)
	items, err := h.service.ListSelectSpecialization(c.Request().Context(), c.QueryParam("search"), actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", items, nil)
}

// ─── GetSpecializationByID ───────────────────────────────────────────────────
//
//	@Summary		Get Specialization
//	@Description	Get Specialization by :id
//	@Tags			kepegawaian/jabatan/specializations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Specialization ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.SpecializationResponse}
//	@Router			/kepegawaian/jabatan/specializations/{id} [get]
func (h *KepegawaianJabatanHandler) GetSpecializationByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetSpecializationByID(c.Request().Context(), id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreateSpecialization ────────────────────────────────────────────────────
//
//	@Summary		Create Specialization
//	@Description	Create New Specialization
//	@Tags			kepegawaian/jabatan/specializations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateSpecializationRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.SpecializationResponse}
//	@Router			/kepegawaian/jabatan/specializations [post]
func (h *KepegawaianJabatanHandler) CreateSpecialization(c *echo.Context) error {
	var req dto.CreateSpecializationRequest
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
	item, err := h.service.CreateSpecialization(c.Request().Context(), &req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdateSpecialization ────────────────────────────────────────────────────
//
//	@Summary		Update Specialization
//	@Description	Update Specialization by :id
//	@Tags			kepegawaian/jabatan/specializations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int								true	"Specialization ID"
//	@Param			body	body		dto.UpdateSpecializationRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.SpecializationResponse}
//	@Router			/kepegawaian/jabatan/specializations/{id} [put]
func (h *KepegawaianJabatanHandler) UpdateSpecialization(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateSpecializationRequest
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
	item, err := h.service.UpdateSpecialization(c.Request().Context(), id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "Specialization tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeleteSpecialization ────────────────────────────────────────────────────
//
//	@Summary		Delete Specialization
//	@Description	Delete Specialization by :id
//	@Tags			kepegawaian/jabatan/specializations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Specialization ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/kepegawaian/jabatan/specializations/{id} [delete]
func (h *KepegawaianJabatanHandler) DeleteSpecialization(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteSpecialization(c.Request().Context(), id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "Specialization tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
