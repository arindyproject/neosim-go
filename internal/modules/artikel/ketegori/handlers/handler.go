
package handlers

import (
	"neosim_go/internal/modules/artikel/ketegori/contracts"
)


// ArtikelKetegoriHandler defines HTTP handlers
type ArtikelKetegoriHandler struct {
	service contracts.Service
}

// NewArtikelKetegoriHandler membuat instance handler baru
func NewArtikelKetegoriHandler(service contracts.Service) *ArtikelKetegoriHandler {
	return &ArtikelKetegoriHandler{service: service}
}
