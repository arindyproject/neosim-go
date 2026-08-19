package contracts

// Repository defines database operations.
// Method utama KepegawaianPegawai didefinisikan di pegawai_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Repository interface {
	KepegawaianPegawaiRepository
	// GEN:ITEM_REPOSITORY_INTERFACE
}

// Service defines business logic operations.
// Method utama KepegawaianPegawai didefinisikan di pegawai_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Service interface {
	KepegawaianPegawaiService
	// GEN:ITEM_SERVICE_INTERFACE
}
