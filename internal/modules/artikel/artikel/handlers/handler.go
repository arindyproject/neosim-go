
package handlers

import (
	"neosim_go/internal/modules/artikel/artikel/contracts"
)


// ArtikelHandler defines HTTP handlers
type ArtikelHandler struct {
	service contracts.Service
}

// NewArtikelHandler membuat instance handler baru
func NewArtikelHandler(service contracts.Service) *ArtikelHandler {
	return &ArtikelHandler{service: service}
}
