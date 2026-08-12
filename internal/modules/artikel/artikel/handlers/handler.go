package handlers

import (
	"neosim_go/config"
	"neosim_go/internal/modules/artikel/artikel/contracts"
)

// ArtikelHandler adalah satu-satunya struct handler untuk sub-module ini.
// Item baru (mode add-item) TIDAK membuat handler baru — method HTTP-nya
// ditempelkan langsung ke struct ini di file terpisah (mis. handlers/tag_handler.go),
// dengan nama method bersuffix nama item (ListTag, CreateTag, dst) agar tidak
// bentrok dengan method CRUD entitas utama (ListArtikel, CreateArtikel, dst).
type ArtikelHandler struct {
	service contracts.Service
	cfg     *config.Config
}

// NewArtikelHandler membuat instance handler baru
func NewArtikelHandler(service contracts.Service, cfg *config.Config) *ArtikelHandler {
	return &ArtikelHandler{service: service, cfg: cfg}
}
