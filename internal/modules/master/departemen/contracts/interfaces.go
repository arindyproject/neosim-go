package contracts

// Repository defines database operations.
// Method utama MasterDepartemen didefinisikan di departemen_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Repository interface {
	MasterDepartemenRepository
	// GEN:ITEM_REPOSITORY_INTERFACE
}

// Service defines business logic operations.
// Method utama MasterDepartemen didefinisikan di departemen_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Service interface {
	MasterDepartemenService
	// GEN:ITEM_SERVICE_INTERFACE
}
