package handlers

import (
	"neosim_go/config"
	"neosim_go/internal/modules/kepegawaian/kontak/contracts"
)

// KepegawaianKontakHandler adalah satu-satunya struct handler untuk sub-module ini.
// Item baru (mode add-item) TIDAK membuat handler baru — method HTTP-nya
// ditempelkan langsung ke struct ini di file terpisah (mis. handlers/tag_handler.go),
// dengan nama method bersuffix nama item (ListTag, CreateTag, dst) agar tidak
// bentrok dengan method CRUD entitas utama (ListKontak, CreateKontak, dst).
type KepegawaianKontakHandler struct {
	service contracts.Service
	cfg     *config.Config
}

// NewKepegawaianKontakHandler membuat instance handler baru
func NewKepegawaianKontakHandler(service contracts.Service, cfg *config.Config) *KepegawaianKontakHandler {
	return &KepegawaianKontakHandler{service: service, cfg: cfg}
}
