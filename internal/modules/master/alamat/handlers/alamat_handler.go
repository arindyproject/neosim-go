package handlers

import (
	"net/http"
	"strconv"

	"neosim_go/internal/modules/master/alamat/contracts"
	"neosim_go/internal/modules/master/alamat/dto"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"

	rbacMiddlewares "neosim_go/internal/modules/rbac/middlewares"

	"github.com/labstack/echo/v5"
)

// MasterAlamatHandler defines HTTP handlers
type MasterAlamatHandler struct {
	service contracts.Service
}

// NewMasterAlamatHandler membuat instance handler baru
func NewMasterAlamatHandler(service contracts.Service) *MasterAlamatHandler {
	return &MasterAlamatHandler{service: service}
}

// buildAuthContext membuat AuthContext dari JWT claims di context
func buildAuthContext(c *echo.Context) contracts.AuthContext {
	userID, _ := rbacMiddlewares.GetUserIDFromContext(c)
	isSuperadmin := rbacMiddlewares.IsSuperadmin(c)
	return contracts.AuthContext{
		UserID:       userID,
		IsSuperadmin: isSuperadmin,
	}
}

// ─── Private Helpers ───────────────────────────────────────────────────────────

func parseID(c *echo.Context) (int64, error) {
	return strconv.ParseInt(c.Param("id"), 10, 64)
}

func getActorID(c *echo.Context) *int64 {
	if userID, ok := c.Get("userID").(int64); ok {
		return &userID
	}
	return nil
}

func parsePagination(c *echo.Context) (page, pageSize int) {
	page, pageSize = 1, 10
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
	return
}

// parseOptionalInt64Query mengambil query param sebagai *int64, nil jika kosong/invalid
func parseOptionalInt64Query(c *echo.Context, key string) *int64 {
	raw := c.QueryParam(key)
	if raw == "" {
		return nil
	}
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return nil
	}
	return &v
}

func notFoundStatus(err error, notFoundMsg string) int {
	if err.Error() == notFoundMsg {
		return http.StatusNotFound
	}
	return http.StatusBadRequest
}

// ─── Handlers ──────────────────────────────────────────────────────────────────

// Negara ===========================================================================
// ─────────────── List ─────────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Get list of Negara
//	@Description	Get paginated list of Negara
//	@Tags			master/alamat/negara
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			code		query		string	false	"Filter by code"
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.NegaraResponse}
//	@Router			/master/alamat/negara [get]
//
// ListNegara handles GET /api/v1/master/alamat/negara
func (h *MasterAlamatHandler) ListNegara(c *echo.Context) error {
	page, pageSize := parsePagination(c)
	filter := dto.FilterNegaraRequest{
		Code: c.QueryParam("code"),
		Name: c.QueryParam("name"),
	}

	items, total, err := h.service.ListNegara(page, pageSize, &filter)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── GetByID ────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Get Negara
//	@Description	Get Negara by :id
//	@Tags			master/alamat/negara
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Negara ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.NegaraResponse}
//	@Router			/master/alamat/negara/{id} [get]
//
// GetByIDNegara handles GET /api/v1/master/alamat/negara/:id
func (h *MasterAlamatHandler) GetByIDNegara(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	item, err := h.service.GetByIDNegara(id)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Create Negara
//	@Description	Create New Negara
//	@Tags			master/alamat/negara
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateNegaraRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.NegaraResponse}
//	@Router			/master/alamat/negara [post]
//
// CreateNegara handles POST /api/v1/master/alamat/negara
func (h *MasterAlamatHandler) CreateNegara(c *echo.Context) error {
	var req dto.CreateNegaraRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := buildAuthContext(c)
	item, err := h.service.CreateNegara(&req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Update Negara
//	@Description	Update Negara by :id
//	@Tags			master/alamat/negara
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"Negara ID"
//	@Param			body	body		dto.UpdateNegaraRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.NegaraResponse}
//	@Router			/master/alamat/negara/{id} [put]
//
// UpdateNegara handles PUT /api/v1/master/alamat/negara/:id
func (h *MasterAlamatHandler) UpdateNegara(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateNegaraRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := buildAuthContext(c)
	item, err := h.service.UpdateNegara(id, &req, actor)
	if err != nil {
		return response.Response(c, notFoundStatus(err, "Negara tidak ditemukan"), false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Delete Negara
//	@Description	Delete Negara by :id
//	@Tags			master/alamat/negara
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Negara ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/master/alamat/negara/{id} [delete]
//
// DeleteNegara handles DELETE /api/v1/master/alamat/negara/:id
func (h *MasterAlamatHandler) DeleteNegara(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := buildAuthContext(c)
	if err := h.service.DeleteNegara(id, actor); err != nil {
		return response.Response(c, notFoundStatus(err, "Negara tidak ditemukan"), false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
} // ───────────── Delete ─────────────────────────────────────────────────────────
// Negara ===========================================================================

// Provinsi ==========================================================================
// ─────────────── List ─────────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Get list of Provinsi
//	@Description	Get paginated list of Provinsi
//	@Tags			master/alamat/provinsi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			negara_id	query		int		false	"Filter by negara_id"
//	@Param			code		query		string	false	"Filter by code"
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.ProvinsiResponse}
//	@Router			/master/alamat/provinsi [get]
//
// ListProvinsi handles GET /api/v1/master/alamat/provinsi
func (h *MasterAlamatHandler) ListProvinsi(c *echo.Context) error {
	page, pageSize := parsePagination(c)
	negaraID := parseOptionalInt64Query(c, "negara_id")
	filter := dto.FilterProvinsiRequest{
		Code: c.QueryParam("code"),
		Name: c.QueryParam("name"),
	}

	items, total, err := h.service.ListProvinsi(page, pageSize, negaraID, &filter)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── GetByID ────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Get Provinsi
//	@Description	Get Provinsi detail by :id, termasuk statistik turunan
//	@Tags			master/alamat/provinsi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Provinsi ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.ProvinsiDetailResponse}
//	@Router			/master/alamat/provinsi/{id} [get]
//
// GetByIDProvinsi handles GET /api/v1/master/alamat/provinsi/:id
func (h *MasterAlamatHandler) GetByIDProvinsi(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	item, err := h.service.GetByIDProvinsi(id)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Create Provinsi
//	@Description	Create New Provinsi
//	@Tags			master/alamat/provinsi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateProvinsiRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.ProvinsiResponse}
//	@Router			/master/alamat/provinsi [post]
//
// CreateProvinsi handles POST /api/v1/master/alamat/provinsi
func (h *MasterAlamatHandler) CreateProvinsi(c *echo.Context) error {
	var req dto.CreateProvinsiRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := buildAuthContext(c)
	item, err := h.service.CreateProvinsi(&req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Update Provinsi
//	@Description	Update Provinsi by :id
//	@Tags			master/alamat/provinsi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"Provinsi ID"
//	@Param			body	body		dto.UpdateProvinsiRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.ProvinsiResponse}
//	@Router			/master/alamat/provinsi/{id} [put]
//
// UpdateProvinsi handles PUT /api/v1/master/alamat/provinsi/:id
func (h *MasterAlamatHandler) UpdateProvinsi(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateProvinsiRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := buildAuthContext(c)
	item, err := h.service.UpdateProvinsi(id, &req, actor)
	if err != nil {
		return response.Response(c, notFoundStatus(err, "Provinsi tidak ditemukan"), false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Delete Provinsi
//	@Description	Delete Provinsi by :id
//	@Tags			master/alamat/provinsi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Provinsi ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/master/alamat/provinsi/{id} [delete]
//
// DeleteProvinsi handles DELETE /api/v1/master/alamat/provinsi/:id
func (h *MasterAlamatHandler) DeleteProvinsi(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := buildAuthContext(c)
	if err := h.service.DeleteProvinsi(id, actor); err != nil {
		return response.Response(c, notFoundStatus(err, "Provinsi tidak ditemukan"), false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
} // ───────────── Delete ─────────────────────────────────────────────────────────
// Provinsi ==========================================================================

// Kota/Kabupaten =====================================================================
// ─────────────── List ─────────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Get list of Kota/Kabupaten
//	@Description	Get paginated list of Kota/Kabupaten
//	@Tags			master/alamat/kota
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			provinsi_id	query		int		false	"Filter by provinsi_id"
//	@Param			code		query		string	false	"Filter by code"
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.KotaKabupatenResponse}
//	@Router			/master/alamat/kota [get]
//
// ListKotaKabupaten handles GET /api/v1/master/kota
func (h *MasterAlamatHandler) ListKotaKabupaten(c *echo.Context) error {
	page, pageSize := parsePagination(c)
	provinsiID := parseOptionalInt64Query(c, "provinsi_id")
	filter := dto.FilterKotaKabupatenRequest{
		Code: c.QueryParam("code"),
		Name: c.QueryParam("name"),
	}

	items, total, err := h.service.ListKotaKabupaten(page, pageSize, provinsiID, &filter)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── GetByID ────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Get Kota/Kabupaten
//	@Description	Get Kota/Kabupaten detail by :id, termasuk statistik turunan
//	@Tags			master/alamat/kota
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Kota/Kabupaten ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.KotaKabupatenDetailResponse}
//	@Router			/master/alamat/kota/{id} [get]
//
// GetByIDKotaKabupaten handles GET /api/v1/master/kota/:id
func (h *MasterAlamatHandler) GetByIDKotaKabupaten(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	item, err := h.service.GetByIDKotaKabupaten(id)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Create Kota/Kabupaten
//	@Description	Create New Kota/Kabupaten
//	@Tags			master/alamat/kota
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateKotaKabupatenRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.KotaKabupatenResponse}
//	@Router			/master/alamat/kota [post]
//
// CreateKotaKabupaten handles POST /api/v1/master/kota
func (h *MasterAlamatHandler) CreateKotaKabupaten(c *echo.Context) error {
	var req dto.CreateKotaKabupatenRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := buildAuthContext(c)
	item, err := h.service.CreateKotaKabupaten(&req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Update Kota/Kabupaten
//	@Description	Update Kota/Kabupaten by :id
//	@Tags			master/alamat/kota
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int								true	"Kota/Kabupaten ID"
//	@Param			body	body		dto.UpdateKotaKabupatenRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.KotaKabupatenResponse}
//	@Router			/master/alamat/kota/{id} [put]
//
// UpdateKotaKabupaten handles PUT /api/v1/master/kota/:id
func (h *MasterAlamatHandler) UpdateKotaKabupaten(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateKotaKabupatenRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := buildAuthContext(c)
	item, err := h.service.UpdateKotaKabupaten(id, &req, actor)
	if err != nil {
		return response.Response(c, notFoundStatus(err, "Kota/Kabupaten tidak ditemukan"), false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Delete Kota/Kabupaten
//	@Description	Delete Kota/Kabupaten by :id
//	@Tags			master/alamat/kota
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Kota/Kabupaten ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/master/alamat/kota/{id} [delete]
//
// DeleteKotaKabupaten handles DELETE /api/v1/master/kota/:id
func (h *MasterAlamatHandler) DeleteKotaKabupaten(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := buildAuthContext(c)
	if err := h.service.DeleteKotaKabupaten(id, actor); err != nil {
		return response.Response(c, notFoundStatus(err, "Kota/Kabupaten tidak ditemukan"), false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
} // ───────────── Delete ─────────────────────────────────────────────────────────
// Kota/Kabupaten =====================================================================

// Kecamatan ===========================================================================
// ─────────────── List ─────────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Get list of Kecamatan
//	@Description	Get paginated list of Kecamatan
//	@Tags			master/alamat/kecamatan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			kota_kabupaten_id	query		int		false	"Filter by kota_kabupaten_id"
//	@Param			code				query		string	false	"Filter by code"
//	@Param			name				query		string	false	"Filter by name (partial match)"
//	@Param			page				query		int		false	"Page number"
//	@Param			page_size			query		int		false	"Page size"
//	@Success		200					{object}	response.MyGoResponse{data=[]dto.KecamatanResponse}
//	@Router			/master/alamat/kecamatan [get]
//
// ListKecamatan handles GET /api/v1/master/kecamatan
func (h *MasterAlamatHandler) ListKecamatan(c *echo.Context) error {
	page, pageSize := parsePagination(c)
	kotaKabupatenID := parseOptionalInt64Query(c, "kota_kabupaten_id")
	filter := dto.FilterKecamatanRequest{
		Code: c.QueryParam("code"),
		Name: c.QueryParam("name"),
	}

	items, total, err := h.service.ListKecamatan(page, pageSize, kotaKabupatenID, &filter)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── GetByID ────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Get Kecamatan
//	@Description	Get Kecamatan detail by :id, termasuk statistik turunan
//	@Tags			master/alamat/kecamatan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Kecamatan ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.KecamatanDetailResponse}
//	@Router			/master/alamat/kecamatan/{id} [get]
//
// GetByIDKecamatan handles GET /api/v1/master/kecamatan/:id
func (h *MasterAlamatHandler) GetByIDKecamatan(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	item, err := h.service.GetByIDKecamatan(id)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Create Kecamatan
//	@Description	Create New Kecamatan
//	@Tags			master/alamat/kecamatan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateKecamatanRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.KecamatanResponse}
//	@Router			/master/alamat/kecamatan [post]
//
// CreateKecamatan handles POST /api/v1/master/kecamatan
func (h *MasterAlamatHandler) CreateKecamatan(c *echo.Context) error {
	var req dto.CreateKecamatanRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := buildAuthContext(c)
	item, err := h.service.CreateKecamatan(&req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Update Kecamatan
//	@Description	Update Kecamatan by :id
//	@Tags			master/alamat/kecamatan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"Kecamatan ID"
//	@Param			body	body		dto.UpdateKecamatanRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.KecamatanResponse}
//	@Router			/master/alamat/kecamatan/{id} [put]
//
// UpdateKecamatan handles PUT /api/v1/master/kecamatan/:id
func (h *MasterAlamatHandler) UpdateKecamatan(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateKecamatanRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := buildAuthContext(c)
	item, err := h.service.UpdateKecamatan(id, &req, actor)
	if err != nil {
		return response.Response(c, notFoundStatus(err, "Kecamatan tidak ditemukan"), false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Delete Kecamatan
//	@Description	Delete Kecamatan by :id
//	@Tags			master/alamat/kecamatan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Kecamatan ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/master/alamat/kecamatan/{id} [delete]
//
// DeleteKecamatan handles DELETE /api/v1/master/kecamatan/:id
func (h *MasterAlamatHandler) DeleteKecamatan(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := buildAuthContext(c)
	if err := h.service.DeleteKecamatan(id, actor); err != nil {
		return response.Response(c, notFoundStatus(err, "Kecamatan tidak ditemukan"), false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
} // ───────────── Delete ─────────────────────────────────────────────────────────
// Kecamatan ===========================================================================

// Kelurahan/Desa ========================================================================
// ─────────────── List ─────────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Get list of Kelurahan/Desa
//	@Description	Get paginated list of Kelurahan/Desa
//	@Tags			master/alamat/desa
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			kecamatan_id	query		int		false	"Filter by kecamatan_id"
//	@Param			code			query		string	false	"Filter by code"
//	@Param			name			query		string	false	"Filter by name (partial match)"
//	@Param			postal_code		query		string	false	"Filter by postal code"
//	@Param			page			query		int		false	"Page number"
//	@Param			page_size		query		int		false	"Page size"
//	@Success		200				{object}	response.MyGoResponse{data=[]dto.KelurahanDesaResponse}
//	@Router			/master/alamat/desa [get]
//
// ListKelurahanDesa handles GET /api/v1/master/desa
func (h *MasterAlamatHandler) ListKelurahanDesa(c *echo.Context) error {
	page, pageSize := parsePagination(c)
	kecamatanID := parseOptionalInt64Query(c, "kecamatan_id")
	filter := dto.FilterKelurahanDesaRequest{
		Code:       c.QueryParam("code"),
		Name:       c.QueryParam("name"),
		PostalCode: c.QueryParam("postal_code"),
	}

	items, total, err := h.service.ListKelurahanDesa(page, pageSize, kecamatanID, &filter)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── GetByID ────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Get Kelurahan/Desa
//	@Description	Get Kelurahan/Desa detail by :id, termasuk jalur hierarki lengkap
//	@Tags			master/alamat/desa
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Kelurahan/Desa ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.KelurahanDesaDetailResponse}
//	@Router			/master/alamat/desa/{id} [get]
//
// GetByIDKelurahanDesa handles GET /api/v1/master/desa/:id
func (h *MasterAlamatHandler) GetByIDKelurahanDesa(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	item, err := h.service.GetByIDKelurahanDesa(id)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Create Kelurahan/Desa
//	@Description	Create New Kelurahan/Desa
//	@Tags			master/alamat/desa
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateKelurahanDesaRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.KelurahanDesaResponse}
//	@Router			/master/alamat/desa [post]
//
// CreateKelurahanDesa handles POST /api/v1/master/desa
func (h *MasterAlamatHandler) CreateKelurahanDesa(c *echo.Context) error {
	var req dto.CreateKelurahanDesaRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := buildAuthContext(c)
	item, err := h.service.CreateKelurahanDesa(&req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Update Kelurahan/Desa
//	@Description	Update Kelurahan/Desa by :id
//	@Tags			master/alamat/desa
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int								true	"Kelurahan/Desa ID"
//	@Param			body	body		dto.UpdateKelurahanDesaRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.KelurahanDesaResponse}
//	@Router			/master/alamat/desa/{id} [put]
//
// UpdateKelurahanDesa handles PUT /api/v1/master/desa/:id
func (h *MasterAlamatHandler) UpdateKelurahanDesa(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateKelurahanDesaRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := buildAuthContext(c)
	item, err := h.service.UpdateKelurahanDesa(id, &req, actor)
	if err != nil {
		return response.Response(c, notFoundStatus(err, "Kelurahan/Desa tidak ditemukan"), false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Delete Kelurahan/Desa
//	@Description	Delete Kelurahan/Desa by :id
//	@Tags			master/alamat/desa
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Kelurahan/Desa ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/master/alamat/desa/{id} [delete]
//
// DeleteKelurahanDesa handles DELETE /api/v1/master/desa/:id
func (h *MasterAlamatHandler) DeleteKelurahanDesa(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := buildAuthContext(c)
	if err := h.service.DeleteKelurahanDesa(id, actor); err != nil {
		return response.Response(c, notFoundStatus(err, "Kelurahan/Desa tidak ditemukan"), false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
} // ───────────── Delete ─────────────────────────────────────────────────────────
// Kelurahan/Desa ========================================================================
