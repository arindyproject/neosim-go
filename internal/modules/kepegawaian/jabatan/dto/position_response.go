package dto

import (
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// PositionKategoriSimpelResponse representasi ringkas PositionKategori,
// dipakai sebagai nested object di PositionResponse (pola sama dengan TipeSimpelResponse).

// PositionSimpelResponse representasi ringkas Position, dipakai untuk field
// Parent/Children supaya tidak rekursif penuh (hindari payload meledak di hierarki dalam).
type PositionSimpelResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// DepartemenSimpelResponse representasi ringkas Departemen (dari modul
// master/departemen), dipakai sebagai nested object di PositionResponse.
//
// ASUMSI: field "Nama" mengikuti konvensi model master lain di project ini
// (nama_lengkap, nama_institusi, dst). Sesuaikan ke "Name" kalau ternyata
// model Departemen kamu pakai field bahasa Inggris seperti Position.
type DepartemenSimpelResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// PositionResponse response untuk single Position
type PositionResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`

	//PositionKategoriID int64                           `json:"position_kategori_id"`
	PositionKategori *PositionKategoriSimpelResponse `json:"position_kategori,omitempty"`

	ParentID *int64                   `json:"parent_id"`
	Parent   *PositionSimpelResponse  `json:"parent,omitempty"`
	Children []PositionSimpelResponse `json:"children,omitempty"`

	//DepartmentID int64                      `json:"department_id"`
	Department *DepartemenSimpelResponse `json:"department,omitempty"`

	LevelHierarki int16  `json:"level_hierarki"`
	Kuota         *int16 `json:"kuota"`
	IsAktif       bool   `json:"is_aktif"`

	CreatedBy *he.UserData     `json:"created_by"`
	UpdatedBy *he.UserData     `json:"updated_by"`
	CreatedAt types.CustomTime `json:"created_at"`
	UpdatedAt types.CustomTime `json:"updated_at"`
}

type PositionResponseParams struct {
	Position   *models.Position
	Creator    *he.UserData
	Updater    *he.UserData
	Department *DepartemenSimpelResponse // hasil lookup lintas-module, tidak ikut GORM preload Position
}

// ToPositionResponse mengubah model menjadi response
func ToPositionResponse(params PositionResponseParams) *PositionResponse {
	if params.Position == nil {
		return nil
	}

	m := params.Position

	var kategoriResponse *PositionKategoriSimpelResponse
	if m.PositionKategori != nil {
		kategoriResponse = &PositionKategoriSimpelResponse{
			ID:    m.PositionKategori.ID,
			Code:  m.PositionKategori.Code,
			Label: m.PositionKategori.Label,
		}
	}

	var parentResponse *PositionSimpelResponse
	if m.Parent != nil {
		parentResponse = &PositionSimpelResponse{
			ID:   m.Parent.ID,
			Name: m.Parent.Name,
		}
	}

	var childrenResponse []PositionSimpelResponse
	if len(m.Children) > 0 {
		childrenResponse = make([]PositionSimpelResponse, 0, len(m.Children))
		for _, c := range m.Children {
			if c == nil {
				continue
			}
			childrenResponse = append(childrenResponse, PositionSimpelResponse{
				ID:   c.ID,
				Name: c.Name,
			})
		}
	}

	return &PositionResponse{
		ID:               m.ID,
		Name:             m.Name,
		Description:      m.Description,
		PositionKategori: kategoriResponse,
		ParentID:         m.ParentID,
		Parent:           parentResponse,
		Children:         childrenResponse,
		Department:       params.Department,
		LevelHierarki:    m.LevelHierarki,
		Kuota:            m.Kuota,
		IsAktif:          m.IsAktif,
		CreatedBy:        params.Creator,
		UpdatedBy:        params.Updater,
		CreatedAt:        types.CustomTime(m.CreatedAt),
		UpdatedAt:        types.CustomTime(m.UpdatedAt),
	}
}

// ToPositionListResponse mengubah slice model menjadi slice response.
// departmentsMap: key = department_id, dikumpulkan sekali di service layer
// lewat s.departemenRepo, bukan N+1 query per item.
func ToPositionListResponse(
	items []models.Position,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
	departmentsMap map[int64]*DepartemenSimpelResponse,
) []PositionResponse {
	responses := make([]PositionResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData
		var department *DepartemenSimpelResponse

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}
		if departmentsMap != nil && m.DepartmentID != nil {
			department = departmentsMap[*m.DepartmentID]
		}

		responses = append(responses, *ToPositionResponse(PositionResponseParams{
			Position:   &m,
			Creator:    creator,
			Updater:    updater,
			Department: department,
		}))
	}

	return responses
}

// PositionTreeNode representasi satu node di tree bagan organisasi Position.
// Berbeda dari PositionResponse: Children di sini rekursif penuh (bukan
// PositionSimpelResponse satu level), karena memang tujuannya menampilkan
// seluruh cabang di bawahnya.
type PositionTreeNode struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`

	PositionKategori *PositionKategoriSimpelResponse `json:"position_kategori,omitempty"`

	//DepartmentID int64                     `json:"department_id"`
	Department *DepartemenSimpelResponse `json:"department,omitempty"`

	LevelHierarki int16  `json:"level_hierarki"`
	Kuota         *int16 `json:"kuota"`
	IsAktif       bool   `json:"is_aktif"`

	Children []PositionTreeNode `json:"children,omitempty"`
}

// ToPositionTree membangun tree bersarang dari flat list Position yang sudah
// ada (hasil FindAllPositions, sudah preload PositionKategori).
//
// departmentsMap: key = department_id, dikumpulkan sekali di service layer
// lewat s.departemenRepo (departemen beda module, jadi bukan GORM preload
// Position) — sama seperti yang dipakai ToPositionListResponse. Boleh nil
// kalau tree ini tidak butuh info department sama sekali.
//
// Kompleksitas O(n): satu pass bikin map parent_id -> children, lalu DFS
// dari tiap root (parent_id NULL). Position yang parent_id-nya menunjuk ke
// ID yang tidak ada di 'items' (mis. sudah dihapus, atau ternyata bukan
// root maupun anak siapa pun) otomatis tidak muncul di tree — bukan bug,
// itu memang bukan bagian dari struktur organisasi yang aktif.
func ToPositionTree(items []models.Position, departmentsMap map[int64]*DepartemenSimpelResponse) []PositionTreeNode {
	childrenOf := make(map[int64][]models.Position)
	var roots []models.Position

	for _, m := range items {
		if m.ParentID == nil {
			roots = append(roots, m)
		} else {
			childrenOf[*m.ParentID] = append(childrenOf[*m.ParentID], m)
		}
	}

	visited := make(map[int64]bool, len(items))

	var build func(m models.Position) PositionTreeNode
	build = func(m models.Position) PositionTreeNode {
		visited[m.ID] = true

		var kategoriResponse *PositionKategoriSimpelResponse
		if m.PositionKategori != nil {
			kategoriResponse = &PositionKategoriSimpelResponse{
				ID:    m.PositionKategori.ID,
				Label: m.PositionKategori.Label,
			}
		}

		var departmentResponse *DepartemenSimpelResponse
		if departmentsMap != nil && m.DepartmentID != nil {
			departmentResponse = departmentsMap[*m.DepartmentID]
		}

		node := PositionTreeNode{
			ID:               m.ID,
			Name:             m.Name,
			Description:      m.Description,
			PositionKategori: kategoriResponse,
			Department:       departmentResponse,
			LevelHierarki:    m.LevelHierarki,
			Kuota:            m.Kuota,
			IsAktif:          m.IsAktif,
		}

		for _, child := range childrenOf[m.ID] {
			// guard: cegah infinite loop kalau data korup (parent_id bersiklus,
			// mis. A->B->A) yang lolos dari validasi self-parent di service.
			if visited[child.ID] {
				continue
			}
			node.Children = append(node.Children, build(child))
		}

		return node
	}

	nodes := make([]PositionTreeNode, 0, len(roots))
	for _, root := range roots {
		nodes = append(nodes, build(root))
	}
	return nodes
}
