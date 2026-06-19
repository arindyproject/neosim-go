package handlers

import (
	"neosim_go/internal/modules/master/master/contracts"
)

// MasterHandler defines HTTP handlers
type MasterHandler struct {
	service contracts.Service
}

// NewMasterHandler membuat instance handler baru
func NewMasterHandler(service contracts.Service) *MasterHandler {
	return &MasterHandler{service: service}
}
