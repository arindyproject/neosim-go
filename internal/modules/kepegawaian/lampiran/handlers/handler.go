package handlers

import (
	"neosim_go/internal/modules/kepegawaian/lampiran/contracts"
)

// KepegawaianLampiranHandler defines HTTP handlers
type KepegawaianLampiranHandler struct {
	service contracts.Service
}

// NewKepegawaianLampiranHandler membuat instance handler baru
func NewKepegawaianLampiranHandler(service contracts.Service) *KepegawaianLampiranHandler {
	return &KepegawaianLampiranHandler{service: service}
}
