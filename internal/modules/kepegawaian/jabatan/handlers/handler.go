package handlers

import (
	"neosim_go/config"
	"neosim_go/internal/modules/kepegawaian/jabatan/contracts"
)

// KepegawaianJabatanHandler adalah satu-satunya struct handler untuk sub-module ini.
// Item baru (mode add-item) TIDAK membuat handler baru — method HTTP-nya
// ditempelkan langsung ke struct ini di file terpisah (mis. handlers/tag_handler.go),
// dengan nama method bersuffix nama item (ListTag, CreateTag, dst) agar tidak
// bentrok dengan method CRUD entitas utama (ListJabatan, CreateJabatan, dst).
type KepegawaianJabatanHandler struct {
	service contracts.Service
	cfg     *config.Config
}

// NewKepegawaianJabatanHandler membuat instance handler baru
func NewKepegawaianJabatanHandler(service contracts.Service, cfg *config.Config) *KepegawaianJabatanHandler {
	return &KepegawaianJabatanHandler{service: service, cfg: cfg}
}
