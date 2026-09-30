package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/kepegawaian/pegawai/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

var (
	namaDepan    = []string{"Ahmad", "Budi", "Citra", "Dewi", "Eko", "Fitri", "Gunawan", "Hana", "Indra", "Joko", "Kartika", "Lestari", "Mega", "Nur", "Putri", "Rudi", "Siti", "Tono", "Wahyu", "Yuni"}
	namaBelakang = []string{"Saputra", "Wijaya", "Pratama", "Santoso", "Lestari", "Hidayat", "Nugroho", "Rahayu", "Kurniawan", "Setiawan", "Utami", "Permana"}
	kotaLahir    = []string{"Madiun", "Ngawi", "Surabaya", "Malang", "Kediri", "Ponorogo", "Magetan", "Solo", "Yogyakarta", "Jakarta"}
)

func pick(list []string) string { return list[rng.Intn(len(list))] }

// randMasterID mengembalikan 1 atau 2 (ID master yang diasumsikan sudah ada)
func randMasterID() int64 { return int64(rng.Intn(2) + 1) }

func randomDate(fromYear, toYear int) time.Time {
	start := time.Date(fromYear, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(toYear, 12, 31, 0, 0, 0, 0, time.UTC)
	days := int(end.Sub(start).Hours() / 24)
	return start.AddDate(0, 0, rng.Intn(days+1))
}

// KepegawaianPegawaiFactory membuat data KepegawaianPegawai untuk testing/seeding
type KepegawaianPegawaiFactory struct {
	overrides map[string]interface{}
}

func NewKepegawaianPegawaiFactory() *KepegawaianPegawaiFactory {
	return &KepegawaianPegawaiFactory{overrides: make(map[string]interface{})}
}

// With meng-override field. Key yang didukung:
// nik, nomor_pegawai, nama_lengkap, jenis_kelamin_id, golongan_darah_id,
// agama_id, status_pernikahan_id, kewarganegaraan, jenis_id, status_id,
// user_id, is_aktif, tanggal_keluar
func (f *KepegawaianPegawaiFactory) With(field string, value interface{}) *KepegawaianPegawaiFactory {
	f.overrides[field] = value
	return f
}

func (f *KepegawaianPegawaiFactory) Make() *models.KepegawaianPegawai {
	// NIK: kode wilayah 3577 + 12 digit acak (total 16 digit)
	nik := fmt.Sprintf("3577%012d", rng.Int63n(1_000_000_000_000))
	nomor := fmt.Sprintf("PEG-%06d", rng.Intn(1_000_000))

	golDarahID := randMasterID()

	// 90% WNI, 10% WNA
	kewarganegaraan := "WNI"
	if rng.Intn(10) == 0 {
		kewarganegaraan = "WNA"
	}

	m := &models.KepegawaianPegawai{
		NIK:                nik,
		NomorPegawai:       nomor,
		NamaLengkap:        pick(namaDepan) + " " + pick(namaBelakang),
		JenisKelaminID:     randMasterID(),
		TanggalLahir:       randomDate(1970, 2000),
		TempatLahir:        pick(kotaLahir),
		GolonganDarahID:    &golDarahID,
		AgamaID:            randMasterID(),
		StatusPernikahanID: randMasterID(),
		Kewarganegaraan:    &kewarganegaraan,
		TanggalMasuk:       randomDate(2010, 2024),
		IsAktif:            true,
	}

	if v, ok := f.overrides["nik"]; ok {
		m.NIK = v.(string)
	}
	if v, ok := f.overrides["nomor_pegawai"]; ok {
		m.NomorPegawai = v.(string)
	}
	if v, ok := f.overrides["nama_lengkap"]; ok {
		m.NamaLengkap = v.(string)
	}
	if v, ok := f.overrides["jenis_kelamin_id"]; ok {
		m.JenisKelaminID = v.(int64)
	}
	if v, ok := f.overrides["golongan_darah_id"]; ok {
		if v == nil {
			m.GolonganDarahID = nil
		} else {
			id := v.(int64)
			m.GolonganDarahID = &id
		}
	}
	if v, ok := f.overrides["agama_id"]; ok {
		m.AgamaID = v.(int64)
	}
	if v, ok := f.overrides["status_pernikahan_id"]; ok {
		m.StatusPernikahanID = v.(int64)
	}
	if v, ok := f.overrides["kewarganegaraan"]; ok {
		if v == nil {
			m.Kewarganegaraan = nil
		} else {
			s := v.(string)
			m.Kewarganegaraan = &s
		}
	}
	if v, ok := f.overrides["jenis_id"]; ok {
		m.JenisID = v.(int64)
	}
	if v, ok := f.overrides["status_id"]; ok {
		m.StatusID = v.(int64)
	}
	if v, ok := f.overrides["user_id"]; ok {
		id := v.(int64)
		m.UserID = &id
	}
	if v, ok := f.overrides["is_aktif"]; ok {
		m.IsAktif = v.(bool)
	}
	if v, ok := f.overrides["tanggal_keluar"]; ok {
		t := v.(time.Time)
		m.TanggalKeluar = &t
		m.IsAktif = false
	}

	return m
}

func (f *KepegawaianPegawaiFactory) MakeMany(count int) []*models.KepegawaianPegawai {
	items := make([]*models.KepegawaianPegawai, count)
	for i := 0; i < count; i++ {
		items[i] = NewKepegawaianPegawaiFactory().Make()
	}
	return items
}
