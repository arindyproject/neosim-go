package handlers

import (
	"neosim_go/config"
	"neosim_go/internal/modules/kepegawaian/pegawai/contracts"
)

// KepegawaianPegawaiHandler adalah satu-satunya struct handler untuk sub-module ini.
// Item baru (mode add-item) TIDAK membuat handler baru — method HTTP-nya
// ditempelkan langsung ke struct ini di file terpisah (mis. handlers/tag_handler.go),
// dengan nama method bersuffix nama item (ListTag, CreateTag, dst) agar tidak
// bentrok dengan method CRUD entitas utama (ListPegawai, CreatePegawai, dst).
type KepegawaianPegawaiHandler struct {
	service contracts.Service
	cfg     *config.Config
}

// NewKepegawaianPegawaiHandler membuat instance handler baru
func NewKepegawaianPegawaiHandler(service contracts.Service, cfg *config.Config) *KepegawaianPegawaiHandler {
	return &KepegawaianPegawaiHandler{service: service, cfg: cfg}
}
