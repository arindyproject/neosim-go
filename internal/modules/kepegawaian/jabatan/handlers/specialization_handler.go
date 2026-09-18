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
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.SpecializationResponse}
//	@Router			/kepegawaian/jabatan/specializations [get]
func (h *KepegawaianJabatanHandler) ListSpecialization(c *echo.Context) error {
	filter := dto.FilterSpecializationRequest{Name: c.QueryParam("name")}
	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListSpecialization(c.Request().Context(), page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, err.Error(), nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
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
//	@Param			id		path		int							true	"Specialization ID"
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
