package main

import (
	"log"

	"github.com/erp-sppg/backend/internal/models"
	"gorm.io/gorm"
)

// SeedSOPCategories seeds the 7 default SOP categories and their checklist items
// based on the SOP Dapur MBG document. Uses FirstOrCreate to avoid duplicates on re-run.
func SeedSOPCategories(db *gorm.DB) {
	log.Println("Seeding SOP categories and checklist items...")

	type checklistSeed struct {
		Nama      string
		Deskripsi string
	}

	type categorySeed struct {
		Nama      string
		Deskripsi string
		Urutan    int
		Items     []checklistSeed
	}

	categories := []categorySeed{
		{
			Nama:      "Higienitas Dapur dan Sanitasi",
			Deskripsi: "Standar kebersihan dan sanitasi area dapur sesuai SOP Dapur MBG",
			Urutan:    1,
			Items: []checklistSeed{
				{Nama: "Lantai dapur bersih, kering, dan tidak licin", Deskripsi: "Pemeriksaan kondisi lantai dapur secara visual"},
				{Nama: "Tempat sampah tertutup dan dibuang secara berkala", Deskripsi: "Pengelolaan sampah sesuai jadwal"},
				{Nama: "Saluran air dan drainase berfungsi baik", Deskripsi: "Pemeriksaan kelancaran saluran pembuangan"},
				{Nama: "Area dapur bebas dari hama dan serangga", Deskripsi: "Pengendalian hama secara rutin"},
			},
		},
		{
			Nama:      "Standar Persiapan Makanan",
			Deskripsi: "Prosedur persiapan dan pengolahan makanan sesuai SOP Dapur MBG",
			Urutan:    2,
			Items: []checklistSeed{
				{Nama: "Bahan makanan dicuci sebelum diolah", Deskripsi: "Pencucian bahan baku sesuai prosedur"},
				{Nama: "Pisau dan talenan terpisah untuk bahan mentah dan matang", Deskripsi: "Pencegahan kontaminasi silang"},
				{Nama: "Makanan dimasak hingga suhu aman", Deskripsi: "Pemastian suhu pemasakan sesuai standar"},
				{Nama: "Bumbu dan bahan tambahan sesuai resep standar", Deskripsi: "Kepatuhan terhadap resep yang ditetapkan"},
				{Nama: "Porsi makanan sesuai standar yang ditetapkan", Deskripsi: "Konsistensi ukuran porsi"},
			},
		},
		{
			Nama:      "Penyimpanan dan Kontrol Suhu",
			Deskripsi: "Standar penyimpanan bahan dan kontrol suhu sesuai SOP Dapur MBG",
			Urutan:    3,
			Items: []checklistSeed{
				{Nama: "Bahan baku disimpan sesuai kategori (kering/basah/beku)", Deskripsi: "Pemisahan penyimpanan berdasarkan jenis bahan"},
				{Nama: "Suhu penyimpanan dingin terjaga (0-5°C)", Deskripsi: "Monitoring suhu lemari pendingin"},
				{Nama: "Suhu penyimpanan beku terjaga (<-18°C)", Deskripsi: "Monitoring suhu freezer"},
				{Nama: "Sistem FIFO (First In First Out) diterapkan", Deskripsi: "Rotasi stok bahan baku"},
			},
		},
		{
			Nama:      "Prosedur Pengiriman",
			Deskripsi: "Standar prosedur pengiriman makanan ke sekolah sesuai SOP Dapur MBG",
			Urutan:    4,
			Items: []checklistSeed{
				{Nama: "Wadah pengiriman bersih dan tertutup rapat", Deskripsi: "Kondisi kontainer pengiriman"},
				{Nama: "Suhu makanan terjaga selama pengiriman", Deskripsi: "Kontrol suhu selama distribusi"},
				{Nama: "Pengiriman tepat waktu sesuai jadwal", Deskripsi: "Ketepatan waktu pengiriman ke sekolah"},
				{Nama: "Dokumentasi serah terima lengkap", Deskripsi: "Kelengkapan bukti pengiriman (POD)"},
			},
		},
		{
			Nama:      "Kebersihan Staf dan APD",
			Deskripsi: "Standar kebersihan personal dan penggunaan APD sesuai SOP Dapur MBG",
			Urutan:    5,
			Items: []checklistSeed{
				{Nama: "Staf menggunakan seragam dan celemek bersih", Deskripsi: "Penggunaan pakaian kerja sesuai standar"},
				{Nama: "Staf menggunakan penutup kepala dan masker", Deskripsi: "Penggunaan APD wajib di area dapur"},
				{Nama: "Staf mencuci tangan sebelum dan sesudah bekerja", Deskripsi: "Kepatuhan prosedur cuci tangan"},
				{Nama: "Staf dalam kondisi sehat saat bekerja", Deskripsi: "Pemeriksaan kesehatan harian staf"},
			},
		},
		{
			Nama:      "Pemeliharaan Peralatan",
			Deskripsi: "Standar pemeliharaan dan perawatan peralatan dapur sesuai SOP Dapur MBG",
			Urutan:    6,
			Items: []checklistSeed{
				{Nama: "Peralatan masak dalam kondisi bersih dan layak pakai", Deskripsi: "Pemeriksaan kondisi peralatan masak"},
				{Nama: "Jadwal perawatan peralatan dilaksanakan rutin", Deskripsi: "Kepatuhan jadwal maintenance"},
				{Nama: "Peralatan rusak segera dilaporkan dan diperbaiki", Deskripsi: "Proses pelaporan dan perbaikan peralatan"},
			},
		},
		{
			Nama:      "Dokumentasi dan Pencatatan",
			Deskripsi: "Standar dokumentasi dan pencatatan operasional sesuai SOP Dapur MBG",
			Urutan:    7,
			Items: []checklistSeed{
				{Nama: "Log produksi harian diisi lengkap", Deskripsi: "Pencatatan aktivitas produksi harian"},
				{Nama: "Catatan penerimaan bahan baku terdokumentasi", Deskripsi: "Dokumentasi GRN (Goods Received Note)"},
				{Nama: "Laporan stok opname dilakukan berkala", Deskripsi: "Pelaksanaan dan dokumentasi stok opname"},
				{Nama: "Catatan suhu harian tersedia dan terisi", Deskripsi: "Log monitoring suhu penyimpanan"},
			},
		},
	}

	categoryCount := 0
	itemCount := 0

	for _, cat := range categories {
		sopCategory := models.SOPCategory{
			Nama:      cat.Nama,
			Deskripsi: cat.Deskripsi,
			Urutan:    cat.Urutan,
			IsActive:  true,
		}
		// Use FirstOrCreate with Nama as the unique key to avoid duplicates
		db.Where(models.SOPCategory{Nama: cat.Nama}).FirstOrCreate(&sopCategory)
		categoryCount++

		// Seed checklist items for this category
		for idx, item := range cat.Items {
			checklistItem := models.SOPChecklistItem{
				SOPCategoryID: sopCategory.ID,
				Nama:          item.Nama,
				Deskripsi:     item.Deskripsi,
				Urutan:        idx + 1,
				IsActive:      true,
			}
			// Use SOPCategoryID + Nama as unique key
			db.Where(models.SOPChecklistItem{
				SOPCategoryID: sopCategory.ID,
				Nama:          item.Nama,
			}).FirstOrCreate(&checklistItem)
			itemCount++
		}
	}

	log.Printf("Seeded %d SOP categories with %d checklist items\n", categoryCount, itemCount)
}