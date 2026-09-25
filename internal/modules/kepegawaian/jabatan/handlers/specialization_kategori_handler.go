package handlers

import (
	"io"
	"net/http"

	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"
	"neosim_go/internal/shared/binding"

	"github.com/labstack/echo/v5"
)

// Method di bawah ini ditempelkan ke struct KepegawaianJabatanHandler yang
// sama dengan handler entitas utama (lihat handlers/handler.go). Nama method
// diberi suffix SpecializationKategori agar tidak bentrok dengan method entitas utama
// pada struct handler yang sama.

// ─── ListSpecializationKategori ──────────────────────────────────────────────────────
//
//	@Summary		Get list of SpecializationKategori
//	@Description	Get paginated list of SpecializationKategori
//	@Tags			kepegawaian/jabatan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.SpecializationKategoriResponse}
//	@Router			/kepegawaian/jabatan/specialization_kategoris [get]
func (h *KepegawaianJabatanHandler) ListSpecializationKategori(c *echo.Context) error {
	filter := dto.FilterSpecializationKategoriRequest{Name: c.QueryParam("name")}
	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListSpecializationKategori(c.Request().Context(),page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, err.Error(), nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetSpecializationKategoriByID ───────────────────────────────────────────────────
//
//	@Summary		Get SpecializationKategori
//	@Description	Get SpecializationKategori by :id
//	@Tags			kepegawaian/jabatan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"SpecializationKategori ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.SpecializationKategoriResponse}
//	@Router			/kepegawaian/jabatan/specialization_kategoris/{id} [get]
func (h *KepegawaianJabatanHandler) GetSpecializationKategoriByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetSpecializationKategoriByID(c.Request().Context(),id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreateSpecializationKategori ────────────────────────────────────────────────────
//
//	@Summary		Create SpecializationKategori
//	@Description	Create New SpecializationKategori
//	@Tags			kepegawaian/jabatan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateSpecializationKategoriRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.SpecializationKategoriResponse}
//	@Router			/kepegawaian/jabatan/specialization_kategoris [post]
func (h *KepegawaianJabatanHandler) CreateSpecializationKategori(c *echo.Context) error {
	var req dto.CreateSpecializationKategoriRequest
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
	item, err := h.service.CreateSpecializationKategori(c.Request().Context(),&req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdateSpecializationKategori ────────────────────────────────────────────────────
//
//	@Summary		Update SpecializationKategori
//	@Description	Update SpecializationKategori by :id
//	@Tags			kepegawaian/jabatan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"SpecializationKategori ID"
//	@Param			body	body		dto.UpdateSpecializationKategoriRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.SpecializationKategoriResponse}
//	@Router			/kepegawaian/jabatan/specialization_kategoris/{id} [put]
func (h *KepegawaianJabatanHandler) UpdateSpecializationKategori(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateSpecializationKategoriRequest
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
	item, err := h.service.UpdateSpecializationKategori(c.Request().Context(),id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "SpecializationKategori tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeleteSpecializationKategori ────────────────────────────────────────────────────
//
//	@Summary		Delete SpecializationKategori
//	@Description	Delete SpecializationKategori by :id
//	@Tags			kepegawaian/jabatan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"SpecializationKategori ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/kepegawaian/jabatan/specialization_kategoris/{id} [delete]
func (h *KepegawaianJabatanHandler) DeleteSpecializationKategori(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteSpecializationKategori(c.Request().Context(),id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "SpecializationKategori tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
