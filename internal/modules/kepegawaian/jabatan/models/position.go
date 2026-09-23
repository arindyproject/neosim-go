package models

import (
	"time"

	"gorm.io/gorm"
)

// Position represents the kepegawaian_jabatan_positions table in database
type Position struct {
	ID          int64   `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name        string  `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description *string `gorm:"column:description;type:text" json:"description"`

	// Kategori jabatan: direksi | manajerial | koordinator | pelaksana.
	// Dinormalisasi sebagai FK ke PositionKategori, bukan varchar enum,
	// supaya konsisten dengan pola tabel kategori lain di modul ini.
	PositionKategoriID int64             `gorm:"column:position_kategori_id;not null;index" json:"position_kategori_id"`
	PositionKategori   *PositionKategori `gorm:"foreignKey:PositionKategoriID" json:"position_kategori,omitempty"`

	// Self-referencing untuk bagan organisasi berjenjang.
	// ParentID NULL = puncak hierarki (mis. Direktur).
	ParentID *int64      `gorm:"column:parent_id;index" json:"parent_id"`
	Parent   *Position   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []*Position `gorm:"foreignKey:ParentID" json:"children,omitempty"`

	// Unit tempat jabatan ini melekat. Nullable karena sebagian jabatan
	// struktural lintas unit (Direktur, Wakil Direktur tidak terikat 1 unit).
	// Disimpan sebagai int64 polos (bukan GORM association) kalau modul
	// department ada di module lain, untuk menghindari import silang
	// antar module dalam arsitektur modular monolith.
	DepartmentID *int64 `gorm:"column:department_id;index" json:"department_id"`

	// 1 = level tertinggi (Direktur), makin besar makin rendah.
	// Dipakai untuk sorting bagan organisasi dan alur approval berjenjang.
	LevelHierarki int16 `gorm:"column:level_hierarki;not null;default:1" json:"level_hierarki"`

	// Jumlah slot yang tersedia untuk jabatan ini (mis. Kepala Ruangan = 1 per unit).
	// NULL berarti tidak dibatasi.
	Kuota *int16   `gorm:"column:kuota" json:"kuota"`
	Point *float64 `gorm:"column:point;type:decimal(10,2);default:0.00" json:"point"`

	IsAktif bool `gorm:"column:is_aktif;not null;default:true" json:"is_aktif"`

	CreatedBy *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`
}

func (Position) TableName() string {
	return "kepegawaian_jabatan_positions"
}
