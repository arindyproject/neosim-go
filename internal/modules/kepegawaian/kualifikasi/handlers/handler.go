package handlers

import (
	"neosim_go/config"
	"neosim_go/internal/modules/kepegawaian/kualifikasi/contracts"
)

// KepegawaianKualifikasiHandler adalah satu-satunya struct handler untuk sub-module ini.
// Item baru (mode add-item) TIDAK membuat handler baru — method HTTP-nya
// ditempelkan langsung ke struct ini di file terpisah (mis. handlers/tag_handler.go),
// dengan nama method bersuffix nama item (ListTag, CreateTag, dst) agar tidak
// bentrok dengan method CRUD entitas utama (ListKualifikasi, CreateKualifikasi, dst).
type KepegawaianKualifikasiHandler struct {
	service contracts.Service
	cfg     *config.Config
}

// NewKepegawaianKualifikasiHandler membuat instance handler baru
func NewKepegawaianKualifikasiHandler(service contracts.Service, cfg *config.Config) *KepegawaianKualifikasiHandler {
	return &KepegawaianKualifikasiHandler{service: service, cfg: cfg}
}
