package handlers

import (
	"neosim_go/config"
	"neosim_go/internal/modules/kepegawaian/pendidikan/contracts"
)

// KepegawaianPendidikanHandler adalah satu-satunya struct handler untuk sub-module ini.
// Item baru (mode add-item) TIDAK membuat handler baru — method HTTP-nya
// ditempelkan langsung ke struct ini di file terpisah (mis. handlers/tag_handler.go),
// dengan nama method bersuffix nama item (ListTag, CreateTag, dst) agar tidak
// bentrok dengan method CRUD entitas utama (ListPendidikan, CreatePendidikan, dst).
type KepegawaianPendidikanHandler struct {
	service contracts.Service
	cfg     *config.Config
}

// NewKepegawaianPendidikanHandler membuat instance handler baru
func NewKepegawaianPendidikanHandler(service contracts.Service, cfg *config.Config) *KepegawaianPendidikanHandler {
	return &KepegawaianPendidikanHandler{service: service, cfg: cfg}
}
