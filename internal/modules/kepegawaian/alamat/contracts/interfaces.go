package contracts

// Repository defines database operations.
// Method utama KepegawaianAlamat didefinisikan di alamat_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Repository interface {
	KepegawaianAlamatRepository
	TipeRepository
	// GEN:ITEM_REPOSITORY_INTERFACE
}

// Service defines business logic operations.
// Method utama KepegawaianAlamat didefinisikan di alamat_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Service interface {
	KepegawaianAlamatService
	TipeService
	// GEN:ITEM_SERVICE_INTERFACE
}
