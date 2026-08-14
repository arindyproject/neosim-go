package handlers

import (
	"neosim_go/config"
	"neosim_go/internal/modules/kepegawaian/identifier/contracts"
)

// KepegawaianIdentifierHandler defines HTTP handlers
type KepegawaianIdentifierHandler struct {
	service contracts.Service
	cfg     *config.Config
}

// NewKepegawaianIdentifierHandler membuat instance handler baru
func NewKepegawaianIdentifierHandler(service contracts.Service, cfg *config.Config) *KepegawaianIdentifierHandler {
	return &KepegawaianIdentifierHandler{service: service, cfg: cfg}
}
