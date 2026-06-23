package handlers

import (
	"neosim_go/internal/modules/kepegawaian/identifikasi/contracts"
)

// KepegawaianIdentifikasiHandler defines HTTP handlers
type KepegawaianIdentifikasiHandler struct {
	service contracts.Service
}

// NewKepegawaianIdentifikasiHandler membuat instance handler baru
func NewKepegawaianIdentifikasiHandler(service contracts.Service) *KepegawaianIdentifikasiHandler {
	return &KepegawaianIdentifikasiHandler{service: service}
}
