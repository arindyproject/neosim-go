package contracts

// Repository defines database operations.
// Method utama ArtikelKategori didefinisikan di kategori_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Repository interface {
	ArtikelKategoriRepository
	TagRepository
	// GEN:ITEM_REPOSITORY_INTERFACE
}

// Service defines business logic operations.
// Method utama ArtikelKategori didefinisikan di kategori_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Service interface {
	ArtikelKategoriService
	TagService
	// GEN:ITEM_SERVICE_INTERFACE
}
