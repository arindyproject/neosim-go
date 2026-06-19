
package repositories
import (
	"neosim_go/internal/modules/artikel/ketegori/contracts"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewArtikelKetegoriRepository membuat instance repository baru
func NewArtikelKetegoriRepository(db *gorm.DB) contracts.Repository {
	return &repository{db: db}
}
