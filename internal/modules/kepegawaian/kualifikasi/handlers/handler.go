package handlers

import (
	"neosim_go/internal/modules/kepegawaian/kualifikasi/contracts"
)

// KepegawaianKualifikasiHandler defines HTTP handlers
type KepegawaianKualifikasiHandler struct {
	service contracts.Service
}

// NewKepegawaianKualifikasiHandler membuat instance handler baru
func NewKepegawaianKualifikasiHandler(service contracts.Service) *KepegawaianKualifikasiHandler {
	return &KepegawaianKualifikasiHandler{service: service}
}
