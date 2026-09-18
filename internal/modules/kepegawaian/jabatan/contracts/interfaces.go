package contracts

// Repository defines database operations.
// Method utama KepegawaianJabatan didefinisikan di jabatan_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Repository interface {
	KepegawaianJabatanRepository
	PositionRepository
	JobTitleRepository
	SpecializationRepository
	// GEN:ITEM_REPOSITORY_INTERFACE
}

// Service defines business logic operations.
// Method utama KepegawaianJabatan didefinisikan di jabatan_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Service interface {
	KepegawaianJabatanService
	PositionService
	JobTitleService
	SpecializationService
	// GEN:ITEM_SERVICE_INTERFACE
}
