package handlers

import (
	"neosim_go/config"
	"neosim_go/internal/modules/artikel/kategori/contracts"
)

// ArtikelKategoriHandler defines HTTP handlers
type ArtikelKategoriHandler struct {
	service contracts.Service
	cfg     *config.Config
}

// NewArtikelKategoriHandler membuat instance handler baru
func NewArtikelKategoriHandler(service contracts.Service, cfg *config.Config) *ArtikelKategoriHandler {
	return &ArtikelKategoriHandler{service: service, cfg: cfg}
}
