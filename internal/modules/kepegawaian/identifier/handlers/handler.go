package handlers

import (
	"neosim_go/config"
	"neosim_go/internal/modules/kepegawaian/identifier/contracts"
)

// KepegawaianIdentifierHandler adalah satu-satunya struct handler untuk sub-module ini.
// Item baru (mode add-item) TIDAK membuat handler baru — method HTTP-nya
// ditempelkan langsung ke struct ini di file terpisah (mis. handlers/tag_handler.go),
// dengan nama method bersuffix nama item (ListTag, CreateTag, dst) agar tidak
// bentrok dengan method CRUD entitas utama (ListIdentifier, CreateIdentifier, dst).
type KepegawaianIdentifierHandler struct {
	service contracts.Service
	cfg     *config.Config
}

// NewKepegawaianIdentifierHandler membuat instance handler baru
func NewKepegawaianIdentifierHandler(service contracts.Service, cfg *config.Config) *KepegawaianIdentifierHandler {
	return &KepegawaianIdentifierHandler{service: service, cfg: cfg}
}
