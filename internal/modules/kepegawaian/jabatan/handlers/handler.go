package handlers

import (
	"neosim_go/internal/modules/kepegawaian/jabatan/contracts"
)

// KepegawaianJabatanHandler defines HTTP handlers
type KepegawaianJabatanHandler struct {
	service contracts.Service
}

// NewKepegawaianJabatanHandler membuat instance handler baru
func NewKepegawaianJabatanHandler(service contracts.Service) *KepegawaianJabatanHandler {
	return &KepegawaianJabatanHandler{service: service}
}
