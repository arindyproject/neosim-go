package contracts

// Repository defines database operations.
// Method utama KepegawaianKontak didefinisikan di kontak_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Repository interface {
	KepegawaianKontakRepository
	TipeRepository
	// GEN:ITEM_REPOSITORY_INTERFACE
}

// Service defines business logic operations.
// Method utama KepegawaianKontak didefinisikan di kontak_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Service interface {
	KepegawaianKontakService
	TipeService
	// GEN:ITEM_SERVICE_INTERFACE
}
