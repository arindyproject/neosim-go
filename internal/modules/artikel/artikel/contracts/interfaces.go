package contracts

// Repository defines database operations.
// Method utama Artikel didefinisikan di artikel_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Repository interface {
	ArtikelRepository
	// GEN:ITEM_REPOSITORY_INTERFACE
}

// Service defines business logic operations.
// Method utama Artikel didefinisikan di artikel_interfaces.go.
// Item tambahan (mode add-item) di-embed otomatis lewat marker di bawah.
type Service interface {
	ArtikelService
	// GEN:ITEM_SERVICE_INTERFACE
}
