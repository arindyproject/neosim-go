package handlers

import (
	"neosim_go/internal/modules/kepegawaian/identifier/contracts"
)

// KepegawaianIdentifierHandler defines HTTP handlers
type KepegawaianIdentifierHandler struct {
	service contracts.Service
}

// NewKepegawaianIdentifierHandler membuat instance handler baru
func NewKepegawaianIdentifierHandler(service contracts.Service) *KepegawaianIdentifierHandler {
	return &KepegawaianIdentifierHandler{service: service}
}
