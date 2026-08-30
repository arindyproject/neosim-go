package contracts

// Repository defines database operations.
// Method utama KepegawaianKualifikasi didefinisikan di kualifikasi_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Repository interface {
	KepegawaianKualifikasiRepository
	TipeRepository
	// GEN:ITEM_REPOSITORY_INTERFACE
}

// Service defines business logic operations.
// Method utama KepegawaianKualifikasi didefinisikan di kualifikasi_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Service interface {
	KepegawaianKualifikasiService
	TipeService
	// GEN:ITEM_SERVICE_INTERFACE
}
