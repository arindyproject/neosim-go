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
// diberi suffix Position agar tidak bentrok dengan method entitas utama
// pada struct handler yang sama.

// ─── ListPosition ──────────────────────────────────────────────────────
//
//	@Summary		Get list of Position
//	@Description	Get paginated list of Position
//	@Tags			kepegawaian/jabatan/positions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name					query		string	false	"Filter by name (partial match)"
//	@Param			position_kategori_id	query		int		false	"Filter by kategori jabatan"
//	@Param			parent_id				query		int		false	"Filter by parent position ID"
//	@Param			department_id			query		int		false	"Filter by department ID"
//	@Param			is_aktif				query		bool	false	"Filter by status aktif"
//	@Param			is_root					query		bool	false	"true = hanya posisi puncak hierarki (parent_id NULL)"
//	@Param			page					query		int		false	"Page number"
//	@Param			page_size				query		int		false	"Page size"
//	@Success		200						{object}	response.MyGoResponse{data=[]dto.PositionResponse}
//	@Router			/kepegawaian/jabatan/positions [get]
func (h *KepegawaianJabatanHandler) ListPosition(c *echo.Context) error {
	filter := dto.FilterPositionRequest{
		Name: c.QueryParam("name"),
	}

	if v := c.QueryParam("position_kategori_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			filter.PositionKategoriID = &id
		}
	}
	if v := c.QueryParam("parent_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			filter.ParentID = &id
		}
	}
	if v := c.QueryParam("department_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			filter.DepartmentID = &id
		}
	}
	if v := c.QueryParam("is_aktif"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			filter.IsAktif = &b
		}
	}
	if v := c.QueryParam("is_root"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			filter.IsRoot = &b
		}
	}

	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListPosition(c.Request().Context(), page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, err.Error(), nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetPositionByID ───────────────────────────────────────────────────
//
//	@Summary		Get Position
//	@Description	Get Position by :id
//	@Tags			kepegawaian/jabatan/positions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Position ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.PositionResponse}
//	@Router			/kepegawaian/jabatan/positions/{id} [get]
func (h *KepegawaianJabatanHandler) GetPositionByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetPositionByID(c.Request().Context(), id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreatePosition ────────────────────────────────────────────────────
//
//	@Summary		Create Position
//	@Description	Create New Position
//	@Tags			kepegawaian/jabatan/positions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreatePositionRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.PositionResponse}
//	@Router			/kepegawaian/jabatan/positions [post]
func (h *KepegawaianJabatanHandler) CreatePosition(c *echo.Context) error {
	var req dto.CreatePositionRequest
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
	item, err := h.service.CreatePosition(c.Request().Context(), &req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdatePosition ────────────────────────────────────────────────────
//
//	@Summary		Update Position
//	@Description	Update Position by :id
//	@Tags			kepegawaian/jabatan/positions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"Position ID"
//	@Param			body	body		dto.UpdatePositionRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.PositionResponse}
//	@Router			/kepegawaian/jabatan/positions/{id} [put]
func (h *KepegawaianJabatanHandler) UpdatePosition(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdatePositionRequest
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
	item, err := h.service.UpdatePosition(c.Request().Context(), id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "Position tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeletePosition ────────────────────────────────────────────────────
//
//	@Summary		Delete Position
//	@Description	Delete Position by :id
//	@Tags			kepegawaian/jabatan/positions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Position ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/kepegawaian/jabatan/positions/{id} [delete]
func (h *KepegawaianJabatanHandler) DeletePosition(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeletePosition(c.Request().Context(), id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "Position tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}

// ─── GetPositionTree ───────────────────────────────────────────────────
//
//	@Summary		Get Position hierarchy tree
//	@Description	Menampilkan seluruh bagan organisasi Position sebagai tree bersarang, dibangun dari data yang ada di DB (bukan hardcode)
//	@Tags			kepegawaian/jabatan/positions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			only_aktif	query		bool	false	"true = hanya posisi aktif (default: true)"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.PositionTreeNode}
//	@Router			/kepegawaian/jabatan/positions/tree [get]
func (h *KepegawaianJabatanHandler) GetPositionTree(c *echo.Context) error {
	onlyAktif := true
	if v := c.QueryParam("only_aktif"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			onlyAktif = b
		}
	}

	actor := he.BuildAuthContext(c)
	tree, err := h.service.GetPositionTree(c.Request().Context(), onlyAktif, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", tree, nil)
}
