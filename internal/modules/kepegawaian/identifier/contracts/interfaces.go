package contracts

// Repository defines database operations.
// Method utama KepegawaianIdentifier didefinisikan di identifier_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Repository interface {
	KepegawaianIdentifierRepository
	TipeRepository
	// GEN:ITEM_REPOSITORY_INTERFACE
}

// Service defines business logic operations.
// Method utama KepegawaianIdentifier didefinisikan di identifier_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Service interface {
	KepegawaianIdentifierService
	TipeService
	// GEN:ITEM_SERVICE_INTERFACE
}
