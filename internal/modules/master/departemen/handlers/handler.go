package handlers

import (
	"neosim_go/internal/modules/master/departemen/contracts"
)

// MasterDepartemenHandler defines HTTP handlers
type MasterDepartemenHandler struct {
	service contracts.Service
}

// NewMasterDepartemenHandler membuat instance handler baru
func NewMasterDepartemenHandler(service contracts.Service) *MasterDepartemenHandler {
	return &MasterDepartemenHandler{service: service}
}
