package handlers

import (
	"neosim_go/config"
	"neosim_go/internal/modules/kepegawaian/alamat/contracts"
)

// KepegawaianAlamatHandler adalah satu-satunya struct handler untuk sub-module ini.
// Item baru (mode add-item) TIDAK membuat handler baru — method HTTP-nya
// ditempelkan langsung ke struct ini di file terpisah (mis. handlers/tag_handler.go),
// dengan nama method bersuffix nama item (ListTag, CreateTag, dst) agar tidak
// bentrok dengan method CRUD entitas utama (ListAlamat, CreateAlamat, dst).
type KepegawaianAlamatHandler struct {
	service contracts.Service
	cfg     *config.Config
}

// NewKepegawaianAlamatHandler membuat instance handler baru
func NewKepegawaianAlamatHandler(service contracts.Service, cfg *config.Config) *KepegawaianAlamatHandler {
	return &KepegawaianAlamatHandler{service: service, cfg: cfg}
}
