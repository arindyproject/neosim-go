package handlers

import (
	"neosim_go/config"
	"neosim_go/internal/modules/artikel/artikel/contracts"
)

// ArtikelHandler defines HTTP handlers
type ArtikelHandler struct {
	service contracts.Service
	cfg     *config.Config
}

// NewArtikelHandler membuat instance handler baru
func NewArtikelHandler(service contracts.Service, cfg *config.Config) *ArtikelHandler {
	return &ArtikelHandler{service: service, cfg: cfg}
}
