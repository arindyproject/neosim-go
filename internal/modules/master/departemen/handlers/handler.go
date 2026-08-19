package handlers

import (
	"neosim_go/config"
	"neosim_go/internal/modules/master/departemen/contracts"
)

// MasterDepartemenHandler adalah satu-satunya struct handler untuk sub-module ini.
// Item baru (mode add-item) TIDAK membuat handler baru — method HTTP-nya
// ditempelkan langsung ke struct ini di file terpisah (mis. handlers/tag_handler.go),
// dengan nama method bersuffix nama item (ListTag, CreateTag, dst) agar tidak
// bentrok dengan method CRUD entitas utama (ListDepartemen, CreateDepartemen, dst).
type MasterDepartemenHandler struct {
	service contracts.Service
	cfg     *config.Config
}

// NewMasterDepartemenHandler membuat instance handler baru
func NewMasterDepartemenHandler(service contracts.Service, cfg *config.Config) *MasterDepartemenHandler {
	return &MasterDepartemenHandler{service: service, cfg: cfg}
}
