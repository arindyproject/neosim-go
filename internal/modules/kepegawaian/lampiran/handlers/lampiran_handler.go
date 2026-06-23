package handlers

import (
	"net/http"
	"strconv"

	"neosim_go/internal/modules/kepegawaian/lampiran/dto"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"
	he "neosim_go/internal/shared/httputil"

	"github.com/labstack/echo/v5"
)

// ─── List ──────────────────────────────────────────────────────────────────────
//
//	@Summary		Get list of KepegawaianLampiran
//	@Description	Get paginated list of KepegawaianLampiran
//	@Tags			kepegawaian/lampiran
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.KepegawaianLampiranResponse}
//	@Router			/api/v1/kepegawaian/lampiran [get]
func (h *KepegawaianLampiranHandler) List(c *echo.Context) error {
	page, pageSize := 1, 10
	filter := dto.FilterKepegawaianLampiranRequest{
		Name: c.QueryParam("name"),
	}
	if p := c.QueryParam("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.QueryParam("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.List(page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetByID ───────────────────────────────────────────────────────────────────
//
//	@Summary		Get KepegawaianLampiran
//	@Description	Get KepegawaianLampiran by :id
//	@Tags			kepegawaian/lampiran
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"KepegawaianLampiran ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.KepegawaianLampiranResponse}
//	@Router			/api/v1/kepegawaian/lampiran/{id} [get]
func (h *KepegawaianLampiranHandler) GetByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetByID(id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── Create ────────────────────────────────────────────────────────────────────
//
//	@Summary		Create KepegawaianLampiran
//	@Description	Create New KepegawaianLampiran
//	@Tags			kepegawaian/lampiran
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateKepegawaianLampiranRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.KepegawaianLampiranResponse}
//	@Router			/api/v1/kepegawaian/lampiran [post]
func (h *KepegawaianLampiranHandler) Create(c *echo.Context) error {
	var req dto.CreateKepegawaianLampiranRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.Create(&req, he.GetActorID(c), actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── Update ────────────────────────────────────────────────────────────────────
//
//	@Summary		Update KepegawaianLampiran
//	@Description	Update KepegawaianLampiran by :id
//	@Tags			kepegawaian/lampiran
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"KepegawaianLampiran ID"
//	@Param			body	body		dto.UpdateKepegawaianLampiranRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.KepegawaianLampiranResponse}
//	@Router			/api/v1/kepegawaian/lampiran/{id} [put]
func (h *KepegawaianLampiranHandler) Update(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	var req dto.UpdateKepegawaianLampiranRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.Update(id, &req, he.GetActorID(c), actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "KepegawaianLampiran tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── Delete ────────────────────────────────────────────────────────────────────
//
//	@Summary		Delete KepegawaianLampiran
//	@Description	Delete KepegawaianLampiran by :id
//	@Tags			kepegawaian/lampiran
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"KepegawaianLampiran ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/api/v1/kepegawaian/lampiran/{id} [delete]
func (h *KepegawaianLampiranHandler) Delete(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.Delete(id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "KepegawaianLampiran tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
