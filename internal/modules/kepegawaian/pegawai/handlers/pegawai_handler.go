package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"neosim_go/internal/modules/kepegawaian/pegawai/dto"
	"neosim_go/internal/shared/binding"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"

	"github.com/labstack/echo/v5"
)

// ─── ListPegawai ─────────────────────────────────────────────────────
//
//	@Summary		Get list of KepegawaianPegawai
//	@Description	Get paginated list of KepegawaianPegawai
//	@Tags			kepegawaian/pegawai
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name					query		string	false	"Filter by nama_lengkap (partial match)"
//	@Param			nik						query		string	false	"Filter by nik (partial match)"
//	@Param			nomor_pegawai			query		string	false	"Filter by nomor_pegawai (partial match)"
//	@Param			jenis_kelamin_id		query		int		false	"Filter by jenis_kelamin_id"
//	@Param			agama_id				query		int		false	"Filter by agama_id"
//	@Param			status_pernikahan_id	query		int		false	"Filter by status_pernikahan_id"
//	@Param			jenis_id				query		int		false	"Filter by jenis_id"
//	@Param			status_id				query		int		false	"Filter by status_id"
//	@Param			is_aktif				query		bool	false	"Filter by is_aktif"
//	@Param			sort_by					query		string	false	"Sort column"	Enums(nik, ihs_number, nomor_pegawai, nama_lengkap, jenis_kelamin_id, tanggal_lahir, tempat_lahir, golongan_darah_id, agama_id, status_pernikahan_id, kewarganegaraan, tanggal_masuk, tanggal_keluar, jenis_id, status_id, foto_url, is_aktif, created_at, updated_at)
//	@Param			sort_order				query		string	false	"Sort direction"	Enums(asc, desc)
//	@Param			page					query		int		false	"Page number"
//	@Param			page_size				query		int		false	"Page size"
//	@Success		200						{object}	response.MyGoResponse{data=[]dto.KepegawaianPegawaiResponse}
//	@Failure		400						{object}	response.MyGoResponse
//	@Router			/kepegawaian/pegawai [get]
func (h *KepegawaianPegawaiHandler) ListPegawai(c *echo.Context) error {
	// Helper kecil: parse query param opsional jadi pointer, simpan error pertama.
	var parseErr error
	int64Q := func(key string) *int64 {
		v := c.QueryParam(key)
		if v == "" {
			return nil
		}
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			if parseErr == nil {
				parseErr = fmt.Errorf("parameter %s harus berupa angka", key)
			}
			return nil
		}
		return &n
	}
	boolQ := func(key string) *bool {
		v := c.QueryParam(key)
		if v == "" {
			return nil
		}
		b, err := strconv.ParseBool(v)
		if err != nil {
			if parseErr == nil {
				parseErr = fmt.Errorf("parameter %s harus berupa true atau false", key)
			}
			return nil
		}
		return &b
	}

	filter := dto.FilterKepegawaianPegawaiRequest{
		// Filter teks
		Name:         c.QueryParam("name"),
		NIK:          c.QueryParam("nik"),
		NomorPegawai: c.QueryParam("nomor_pegawai"),

		// Filter ID (pointer: nil = tidak difilter)
		JenisKelaminID:     int64Q("jenis_kelamin_id"),
		AgamaID:            int64Q("agama_id"),
		StatusPernikahanID: int64Q("status_pernikahan_id"),
		JenisID:            int64Q("jenis_id"),
		StatusID:           int64Q("status_id"),

		// Filter boolean (pointer agar false tetap bisa dipakai)
		IsAktif: boolQ("is_aktif"),

		// Sorting
		SortBy:    c.QueryParam("sort_by"),
		SortOrder: c.QueryParam("sort_order"),
	}

	if parseErr != nil {
		return response.Response(c, http.StatusBadRequest, false, parseErr.Error(), nil, nil)
	}

	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListPegawai(c.Request().Context(), page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data : "+err.Error(), nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetPegawaiByID ───────────────────────────────────────────────────
//
//	@Summary		Get KepegawaianPegawai
//	@Description	Get KepegawaianPegawai by :id
//	@Tags			kepegawaian/pegawai
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"KepegawaianPegawai ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.KepegawaianPegawaiResponse}
//	@Router			/kepegawaian/pegawai/{id} [get]
func (h *KepegawaianPegawaiHandler) GetPegawaiByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetPegawaiByID(c.Request().Context(), id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreatePegawai ────────────────────────────────────────────────────
//
//	@Summary		Create KepegawaianPegawai
//	@Description	Create New KepegawaianPegawai
//	@Tags			kepegawaian/pegawai
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateKepegawaianPegawaiRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.KepegawaianPegawaiResponse}
//	@Router			/kepegawaian/pegawai [post]
func (h *KepegawaianPegawaiHandler) CreatePegawai(c *echo.Context) error {
	var req dto.CreateKepegawaianPegawaiRequest
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
	item, err := h.service.CreatePegawai(c.Request().Context(), &req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdatePegawai ────────────────────────────────────────────────────
//
//	@Summary		Update KepegawaianPegawai
//	@Description	Update KepegawaianPegawai by :id
//	@Tags			kepegawaian/pegawai
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"KepegawaianPegawai ID"
//	@Param			body	body		dto.UpdateKepegawaianPegawaiRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.KepegawaianPegawaiResponse}
//	@Router			/kepegawaian/pegawai/{id} [put]
func (h *KepegawaianPegawaiHandler) UpdatePegawai(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	var req dto.UpdateKepegawaianPegawaiRequest
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
	item, err := h.service.UpdatePegawai(c.Request().Context(), id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "KepegawaianPegawai tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeletePegawai ────────────────────────────────────────────────────
//
//	@Summary		Delete KepegawaianPegawai
//	@Description	Delete KepegawaianPegawai by :id
//	@Tags			kepegawaian/pegawai
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"KepegawaianPegawai ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/kepegawaian/pegawai/{id} [delete]
func (h *KepegawaianPegawaiHandler) DeletePegawai(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeletePegawai(c.Request().Context(), id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "KepegawaianPegawai tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
