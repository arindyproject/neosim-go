package handlers

import (
	"neosim_go/internal/modules/master/alamat/contracts"
)

// MasterAlamatHandler defines HTTP handlers
type MasterAlamatHandler struct {
	service contracts.Service
}

// NewMasterAlamatHandler membuat instance handler baru
func NewMasterAlamatHandler(service contracts.Service) *MasterAlamatHandler {
	return &MasterAlamatHandler{service: service}
}
