package modul

import (
	"fmt"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
)

// Fungsi untuk membuat file Excel baru dengan header "Tanggal" dan "Serial Number"
func createExcelFile(filename string) (*excelize.File, error) {
	f := excelize.NewFile()

	// Buat header di kolom A dan B
	f.SetCellValue("Sheet1", "A1", "Tanggal")
	f.SetCellValue("Sheet1", "B1", "Serial Number")
	f.SetCellValue("Sheet1", "C1", "ID")

	// Atur lebar kolom
	// Set lebar kolom A dari Tanggal dan kolom B dari Serial Number
	err := f.SetColWidth("Sheet1", "A", "A", 20) // Lebar kolom A diatur ke 20
	if err != nil {
		return nil, err
	}
	err = f.SetColWidth("Sheet1", "B", "B", 30) // Lebar kolom B diatur ke 30
	if err != nil {
		return nil, err
	}
	err = f.SetColWidth("Sheet1", "C", "C", 30) // Lebar kolom B diatur ke 30
	if err != nil {
		return nil, err
	}

	// Simpan file sementara sebelum ditambah baris
	if err := f.SaveAs(filename); err != nil {
		return nil, err
	}

	return f, nil
}

// Fungsi untuk menambahkan baris baru dengan hanya input SN, tanggal diisi otomatis
func addRowToExcel(filename string, sn string, id string) error {
	// Buka file Excel
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Dapatkan jumlah baris saat ini di sheet
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return err
	}

	// Tentukan baris baru
	rowNum := len(rows) + 1

	// Set tanggal di kolom A dan SN di kolom B
	f.SetCellValue("Sheet1", fmt.Sprintf("A%d", rowNum), time.Now().Format("2006-01-02"))
	f.SetCellValue("Sheet1", fmt.Sprintf("B%d", rowNum), sn)
	f.SetCellValue("Sheet1", fmt.Sprintf("C%d", rowNum), id)

	// Simpan file setelah ditambah baris
	if err := f.SaveAs(filename); err != nil {
		return err
	}

	return nil
}

func buatexsel(bord string) {
	folderPath := "./data_produksi/" + bord
	filename := folderPath + "/" + bord + ".xlsx"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("File tidak ditemukan:", filename)
	} else {
		fmt.Println("File ditemukan:", filename)
		return
	}
	err := ensureFolderExists(folderPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Buat file Excel baru dengan header
	_, err = createExcelFile(filename)
	if err != nil {
		fmt.Println("Gagal membuat file:", err)
		return
	}
}

func tambahbaris(bord string, SN string, id string) error {
	folderPath := "./data_produksi/" + bord
	filename := folderPath + "/" + bord + ".xlsx"
	if err := addRowToExcel(filename, SN, id); err != nil {
		fmt.Println("Gagal menambah baris:", err)
		return err
	}
	return nil
}

func readExcelData(filePath string, sheet string) ([]string, error) {
	// Membuka file Excel
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Membaca semua baris dari sheet pertama
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	// Mengambil data dari kolom pertama setiap baris
	var data []string
	for _, row := range rows {
		if len(row) > 0 {
			data = append(data, row[0]) // Mengambil nilai dari kolom pertama
		}
	}

	return data, nil
}
