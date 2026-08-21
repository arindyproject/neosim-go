package contracts

// Repository defines database operations.
// Method utama KepegawaianPendidikan didefinisikan di pendidikan_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Repository interface {
	KepegawaianPendidikanRepository
	JenjangRepository
	// GEN:ITEM_REPOSITORY_INTERFACE
}

// Service defines business logic operations.
// Method utama KepegawaianPendidikan didefinisikan di pendidikan_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Service interface {
	KepegawaianPendidikanService
	JenjangService
	// GEN:ITEM_SERVICE_INTERFACE
}
