package handlers

import (
	"neosim_go/config"
	"neosim_go/internal/modules/master/master/contracts"
)

// MasterHandler defines HTTP handlers
type MasterHandler struct {
	service contracts.Service
	cfg     *config.Config
}

// NewMasterHandler membuat instance handler baru
func NewMasterHandler(service contracts.Service, cfg *config.Config) *MasterHandler {
	return &MasterHandler{service: service, cfg: cfg}
}
