package handlers

import (
	"neosim_go/config"
	"neosim_go/internal/modules/master/alamat/contracts"
)

// MasterAlamatHandler defines HTTP handlers
type MasterAlamatHandler struct {
	service contracts.Service
	cfg     *config.Config
}

// NewMasterAlamatHandler membuat instance handler baru
func NewMasterAlamatHandler(service contracts.Service, cfg *config.Config) *MasterAlamatHandler {
	return &MasterAlamatHandler{service: service, cfg: cfg}
}
