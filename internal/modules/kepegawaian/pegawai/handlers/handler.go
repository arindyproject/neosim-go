package handlers

import (
	"neosim_go/internal/modules/kepegawaian/pegawai/contracts"
)

// KepegawaianPegawaiHandler defines HTTP handlers
type KepegawaianPegawaiHandler struct {
	service contracts.Service
}

// NewKepegawaianPegawaiHandler membuat instance handler baru
func NewKepegawaianPegawaiHandler(service contracts.Service) *KepegawaianPegawaiHandler {
	return &KepegawaianPegawaiHandler{service: service}
}
