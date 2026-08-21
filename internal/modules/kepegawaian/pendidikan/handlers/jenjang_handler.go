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

// Method di bawah ini ditempelkan ke struct KepegawaianPendidikanHandler yang
// sama dengan handler entitas utama (lihat handlers/handler.go). Nama method
// diberi suffix Jenjang agar tidak bentrok dengan method entitas utama
// pada struct handler yang sama.

// ─── ListJenjang ──────────────────────────────────────────────────────
//
//	@Summary		Get list of Jenjang
//	@Description	Get paginated list of Jenjang
//	@Tags			kepegawaian/pendidikan/jenjang
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.JenjangResponse}
//	@Router			/kepegawaian/pendidikan/jenjangs [get]
func (h *KepegawaianPendidikanHandler) ListJenjang(c *echo.Context) error {
	filter := dto.FilterJenjangRequest{Code: c.QueryParam("code"), Label: c.QueryParam("label")}
	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListJenjang(c.Request().Context(), page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetJenjangByID ───────────────────────────────────────────────────
//
//	@Summary		Get Jenjang
//	@Description	Get Jenjang by :id
//	@Tags			kepegawaian/pendidikan/jenjang
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Jenjang ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.JenjangResponse}
//	@Router			/kepegawaian/pendidikan/jenjangs/{id} [get]
func (h *KepegawaianPendidikanHandler) GetJenjangByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetJenjangByID(c.Request().Context(), id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreateJenjang ────────────────────────────────────────────────────
//
//	@Summary		Create Jenjang
//	@Description	Create New Jenjang
//	@Tags			kepegawaian/pendidikan/jenjang
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateJenjangRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.JenjangResponse}
//	@Router			/kepegawaian/pendidikan/jenjangs [post]
func (h *KepegawaianPendidikanHandler) CreateJenjang(c *echo.Context) error {
	var req dto.CreateJenjangRequest
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
	item, err := h.service.CreateJenjang(c.Request().Context(), &req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdateJenjang ────────────────────────────────────────────────────
//
//	@Summary		Update Jenjang
//	@Description	Update Jenjang by :id
//	@Tags			kepegawaian/pendidikan/jenjang
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"Jenjang ID"
//	@Param			body	body		dto.UpdateJenjangRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.JenjangResponse}
//	@Router			/kepegawaian/pendidikan/jenjangs/{id} [put]
func (h *KepegawaianPendidikanHandler) UpdateJenjang(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateJenjangRequest
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
	item, err := h.service.UpdateJenjang(c.Request().Context(), id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "Jenjang tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeleteJenjang ────────────────────────────────────────────────────
//
//	@Summary		Delete Jenjang
//	@Description	Delete Jenjang by :id
//	@Tags			kepegawaian/pendidikan/jenjang
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Jenjang ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/kepegawaian/pendidikan/jenjangs/{id} [delete]
func (h *KepegawaianPendidikanHandler) DeleteJenjang(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteJenjang(c.Request().Context(), id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "Jenjang tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
