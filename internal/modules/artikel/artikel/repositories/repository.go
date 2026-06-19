
package repositories
import (
	"neosim_go/internal/modules/artikel/artikel/contracts"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewArtikelRepository membuat instance repository baru
func NewArtikelRepository(db *gorm.DB) contracts.Repository {
	return &repository{db: db}
}
