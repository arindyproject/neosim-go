package handlers

import (
	"neosim_go/config"
	"neosim_go/internal/modules/artikel/kategori/contracts"
)

// ArtikelKategoriHandler adalah satu-satunya struct handler untuk sub-module ini.
// Item baru (mode add-item) TIDAK membuat handler baru — method HTTP-nya
// ditempelkan langsung ke struct ini di file terpisah (mis. handlers/tag_handler.go),
// dengan nama method bersuffix nama item (ListTag, CreateTag, dst) agar tidak
// bentrok dengan method CRUD entitas utama (ListKategori, CreateKategori, dst).
type ArtikelKategoriHandler struct {
	service contracts.Service
	cfg     *config.Config
}

// NewArtikelKategoriHandler membuat instance handler baru
func NewArtikelKategoriHandler(service contracts.Service, cfg *config.Config) *ArtikelKategoriHandler {
	return &ArtikelKategoriHandler{service: service, cfg: cfg}
}
