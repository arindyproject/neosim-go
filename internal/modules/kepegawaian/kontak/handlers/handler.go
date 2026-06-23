package handlers

import (
	"neosim_go/internal/modules/kepegawaian/kontak/contracts"
)

// KepegawaianKontakHandler defines HTTP handlers
type KepegawaianKontakHandler struct {
	service contracts.Service
}

// NewKepegawaianKontakHandler membuat instance handler baru
func NewKepegawaianKontakHandler(service contracts.Service) *KepegawaianKontakHandler {
	return &KepegawaianKontakHandler{service: service}
}
