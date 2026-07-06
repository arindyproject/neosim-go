package handlers

import (
	"net/http"
	"strconv"

	"neosim_go/internal/modules/artikel/kategori/contracts"
	"neosim_go/internal/modules/artikel/kategori/dto"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"
	he "neosim_go/internal/shared/httputil"

	"github.com/labstack/echo/v5"
)

// TagHandler defines HTTP handlers for Tag
type TagHandler struct {
	service contracts.TagService
}

func NewTagHandler(service contracts.TagService) *TagHandler {
	return &TagHandler{service: service}
}

// @Summary		Get list of Tag
// @Tags			artikel/kategori
// @Security		BearerAuth
// @Param			name		query		string	false	"Filter by name"
// @Param			page		query		int		false	"Page number"
// @Param			page_size	query		int		false	"Page size"
// @Success		200			{object}	response.MyGoResponse{data=[]dto.TagResponse}
// @Router			/api/v1/artikel/kategori/tags [get]
func (h *TagHandler) List(c *echo.Context) error {
	page, pageSize := 1, 10
	filter := dto.FilterTagRequest{Name: c.QueryParam("name")}
	if p := c.QueryParam("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 { page = v }
	}
	if ps := c.QueryParam("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 && v <= 100 { pageSize = v }
	}
	actor := he.BuildAuthContext(c)
	items, total, err := h.service.List(page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// @Summary		Get Tag by ID
// @Tags			artikel/kategori
// @Security		BearerAuth
// @Param			id	path		int	true	"Tag ID"
// @Success		200	{object}	response.MyGoResponse{data=dto.TagResponse}
// @Router			/api/v1/artikel/kategori/tags/{id} [get]
func (h *TagHandler) GetByID(c *echo.Context) error {
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

// @Summary		Create Tag
// @Tags			artikel/kategori
// @Security		BearerAuth
// @Param			body	body		dto.CreateTagRequest	true	"Create Request"
// @Success		201		{object}	response.MyGoResponse{data=dto.TagResponse}
// @Router			/api/v1/artikel/kategori/tags [post]
func (h *TagHandler) Create(c *echo.Context) error {
	var req dto.CreateTagRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.Create(&req,  actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// @Summary		Update Tag
// @Tags			artikel/kategori
// @Security		BearerAuth
// @Param			id		path		int							true	"Tag ID"
// @Param			body	body		dto.UpdateTagRequest	true	"Update Request"
// @Success		200		{object}	response.MyGoResponse{data=dto.TagResponse}
// @Router			/api/v1/artikel/kategori/tags/{id} [put]
func (h *TagHandler) Update(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	var req dto.UpdateTagRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.Update(id, &req,  actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "Tag tidak ditemukan" { status = http.StatusNotFound }
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// @Summary		Delete Tag
// @Tags			artikel/kategori
// @Security		BearerAuth
// @Param			id	path		int	true	"Tag ID"
// @Success		200	{object}	response.MyGoResponse{}
// @Router			/api/v1/artikel/kategori/tags/{id} [delete]
func (h *TagHandler) Delete(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.Delete(id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "Tag tidak ditemukan" { status = http.StatusNotFound }
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
