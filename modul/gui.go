package modul

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/xuri/excelize/v2"
)

var MyWindow fyne.Window
var RSLonlyflag uint8
var QRcode string
var Updateonlyflag uint8
var namaapp = "PRODUKSI v1.5.2"

func Loginapp() {
	a := app.New()
	err := ReadLoginFile()
	if err == nil {

		Guiapp(a)
		// loginWindow.Close()

	} else {
		// os.Remove("login.txt")
		loginWindow := a.NewWindow("Login")
		loginWindow.Resize(fyne.NewSize(300, 200))

		// Widget untuk username dan password
		usernameEntry := widget.NewEntry()
		usernameEntry.SetPlaceHolder("Username")

		passwordEntry := widget.NewPasswordEntry()
		passwordEntry.SetPlaceHolder("Password")

		// Tombol login
		loginButton := widget.NewButton("Login", func() {
			username := usernameEntry.Text
			password := passwordEntry.Text
			Istrueid := IpaLogin(username, password)
			// Validasi kredensial (contoh: username = "admin", password = "12345")
			if Istrueid == 0 {
				dialog.ShowInformation("Login Berhasil", "Selamat datang!", loginWindow)
				Guiapp(a)
				loginWindow.Close()
				// Buka aplikasi utama
			} else {
				dialog.ShowInformation("Login Gagal", "Username atau password salah.", loginWindow)
			}
		})

		// Layout login
		loginWindow.SetContent(container.NewVBox(
			widget.NewLabel("Silakan login untuk melanjutkan"),
			usernameEntry,
			passwordEntry,
			loginButton,
		))

		// Menampilkan jendela login
		loginWindow.ShowAndRun()
	}
	a.Run()
}

// var Data Bus
func Guiapp(myApp fyne.App) {

	// // Menampilkan window
	// MyWindow.ShowAndRun()
	// myApp := app.New()
	MyWindow = myApp.NewWindow(namaapp)

	// Mengatur ukuran jendela
	MyWindow.Resize(fyne.NewSize(400, 600))

	// Data untuk tipe cell dan tipe battery dari Excel
	typeCellMap := make(map[string][]string)
	// Membaca file Excel dan mengisi typeCellMap
	filePath := "battrymenu.xlsx" // Ganti dengan path file Excel yang benar
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		dialog.ShowInformation("Gagal", "battrymenu.xlsx", MyWindow)
		return
	}

	rows, err := f.GetRows("Sheet1") // Ganti dengan nama sheet yang benar
	if err != nil {
		dialog.ShowInformation("Gagal", "battrymenu.xlsx", MyWindow)
		return
	}
	// Mengisi typeCellMap berdasarkan data dari Excel
	for _, row := range rows {
		if len(row) >= 2 {

			tipeCell := row[0]
			tipeBattery := row[1]
			typeCellMap[tipeCell] = append(typeCellMap[tipeCell], tipeBattery)
		}
	}

	// Membuat dropdown dengan pilihan
	options := []string{"vcu", "bms", "hmi", "keyless"}
	dropdown := widget.NewSelect(options, func(selected string) {
		MyWindow.SetTitle(namaapp + " " + selected)
		// MyWindow.SetTitle("Pilihan: " + selected)
	})

	// Membuat label dan entry untuk jumlah paralel
	modelidLabel := widget.NewLabel("Masukkan modelid:")
	modelidEntry := widget.NewEntry()
	modelidLabel.Hide() // Disembunyikan awalnya
	modelidEntry.Hide() // Disembunyikan awalnya
	modelversion := widget.NewSelect([]string{}, nil)
	modelversion.Hide() // Disembunyikan awalnya
	// Logika untuk menentukan huruf berdasarkan pilihan tipe cell
	modelversion.OnChanged = func(selected string) {
		if selected != "" {
			Data.Modelver = GetStringAfterV(selected)
			Data.Model, _ = VersionToHex(selected)
		}

	}

	// Membuat label dan entry untuk SN
	parallelsLabel := widget.NewLabel("Masukkan Jumlah parallels:")
	parallelsEntry := widget.NewEntry()
	parallelsLabel.Hide() // Disembunyikan awalnya
	parallelsEntry.Hide() // Disembunyikan awalnya

	// Opsi untuk dropdown bulan dan tahun
	months := []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"}
	years := []string{"2020", "2021", "2022", "2023", "2024", "2025", "2026", "2027", "2028", "2029", "2030"}

	// Membuat dropdown untuk Bulan/Tahun Battery
	batteryLabel := widget.NewLabel("Bulan/Tahun Battery:")
	batteryMonthDropdown := widget.NewSelect(months, nil)
	batteryYearDropdown := widget.NewSelect(years, nil)
	batteryLabel.Hide()
	batteryMonthDropdown.Hide()
	batteryYearDropdown.Hide()

	// Membuat dropdown untuk Bulan/Tahun Pembuatan Cell
	cellLabel := widget.NewLabel("Bulan/Tahun Pembuatan Cell:")
	cellMonthDropdown := widget.NewSelect(months, nil)
	cellYearDropdown := widget.NewSelect(years, nil)
	cellLabel.Hide()
	cellMonthDropdown.Hide()
	cellYearDropdown.Hide()

	type_motor, err := readExcelData("menu.xlsx", "type")
	if err != nil {
		log.Fatal(err)
	}
	// Membuat dropdown untuk Tipe Cell
	typemotorLabel := widget.NewLabel("Tipe motor:")
	typemotordropdown := widget.NewSelect(type_motor, func(selected string) {
		// fmt.Println("Selected:", selected)
		Data.Num.Type_ = getStringBeforeDot(selected)
	})
	typemotorLabel.Hide()
	typemotordropdown.Hide()

	type_lcd, err := readExcelData("menu.xlsx", "lcd")
	if err != nil {
		log.Fatal(err)
	}
	// Membuat dropdown untuk Tipe Cell
	typelcdLabel := widget.NewLabel("LCD:")
	typelcddropdown := widget.NewSelect(type_lcd, func(selected string) {
		// fmt.Println("Selected:", selected)
		Data.Num.Jenis = getStringBeforeDot(selected)
	})
	typelcdLabel.Hide()
	typelcddropdown.Hide()

	TypeMCULabel := widget.NewLabel("MCU:")

	mcudropdown := widget.NewSelect([]string{}, nil)
	TypeMCULabel.Hide() // Disembunyikan awalnya
	mcudropdown.Hide()  // Disembunyikan awalnya
	// Logika untuk menentukan huruf berdasarkan pilihan tipe cell
	mcudropdown.OnChanged = func(selected string) {
		if selected != "" {
			Data.Num.MCU = getStringBeforeDot(selected)
		}

	}

	// Membuat dropdown untuk Tipe Cell
	cellbrandLabel := widget.NewLabel("Brand Cell:")
	cellbrandOptions := make([]string, 0, len(typeCellMap))
	for key := range typeCellMap {
		cellbrandOptions = append(cellbrandOptions, key)
	}
	cellbrandDropdown := widget.NewSelect(cellbrandOptions, nil)
	cellbrandLabel.Hide()
	cellbrandDropdown.Hide()
	cellTypeLabel := widget.NewLabel("Tipe Cell:")
	cellTypeDropdown := widget.NewSelect([]string{}, nil)
	cellTypeLabel.Hide()
	cellTypeDropdown.Hide()
	// Variabel untuk menyimpan Brand cell terpilih
	var cellbrandShort string

	// Logika untuk menentukan huruf berdasarkan pilihan tipe cell
	cellbrandDropdown.OnChanged = func(selected string) {

		cellbrandShort = GetFirstLetter(selected)
		if batteries, ok := typeCellMap[selected]; ok {
			cellTypeDropdown.Options = batteries
			cellTypeDropdown.Refresh() // Refresh untuk memperbarui pilihan
		}
		cellTypeLabel.Show()    // Menampilkan label Tipe Battery
		cellTypeDropdown.Show() // Menampilkan dropdown Tipe Battery
	}

	// Variabel untuk menyimpan tipe cell terpilih
	var cellTypeShort string

	// Logika untuk menentukan huruf berdasarkan pilihan tipe cell
	cellTypeDropdown.OnChanged = func(selected string) {
		if selected != "" {
			cellTypeShort = GetFirstLetter(selected)
		}

	}

	// Membuat dropdown untuk tegangan
	teganganLabel := widget.NewLabel("Tegangan Battry:")
	teganganOptions := []string{"12", "24", "36", "48", "60", "72", "84", "96"}
	teganganDropdown := widget.NewSelect(teganganOptions, nil)
	teganganLabel.Hide()
	teganganDropdown.Hide()

	// Variabel untuk menyimpan tipe cell terpilih
	var teganganShort string

	// Logika untuk menentukan huruf berdasarkan pilihan tipe cell
	teganganDropdown.OnChanged = func(selected string) {
		switch selected {
		case "12":
			teganganShort = "0"
		case "24":
			teganganShort = "1"
		case "36":
			teganganShort = "2"
		case "48":
			teganganShort = "3"
		case "60":
			teganganShort = "4"
		case "72":
			teganganShort = "5"
		case "84":
			teganganShort = "6"
		case "96":
			teganganShort = "7"
		default:
			teganganShort = ""
		}
	}

	// Menambahkan RadioGroup untuk Hardware Model (High Set & Low Set)
	hardwareLabel := widget.NewLabel("Hardware Model:")
	hardwareOptions := []string{"High Set", "Low Set"}
	hardwareRadioGroup := widget.NewRadioGroup(hardwareOptions, nil)
	hardwareLabel.Hide()
	hardwareRadioGroup.Hide()

	// Variabel untuk menyimpan nilai angka dari Hardware Model
	var hardwareModelValue uint32
	hardwareModelValue = 0
	// Logika untuk menentukan nilai berdasarkan pilihan Hardware Model
	hardwareRadioGroup.OnChanged = func(selected string) {
		if selected == "High Set" {
			hardwareModelValue = 2
		} else if selected == "Low Set" {
			hardwareModelValue = 1
		}
	}
	RSLonlyflag = 0
	RSLonlyLabel := widget.NewLabel("Hardware Model:")
	check1 := widget.NewCheck("RSL ONLY", func(checked bool) {
		if checked {
			// println("Poin Centang 1 diaktifkan")
			RSLonlyflag = 1

		} else {
			// println("Poin Centang 1 dinonaktifkan")
			RSLonlyflag = 0
		}
	})
	check3 := widget.NewCheck("STM ONLY", func(checked bool) {
		if checked {
			// println("Poin Centang 1 diaktifkan")
			RSLonlyflag = 2

		} else {
			// println("Poin Centang 1 dinonaktifkan")
			RSLonlyflag = 0
		}
	})
	check4 := widget.NewCheck("test ONLY", func(checked bool) {
		if checked {
			// println("Poin Centang 1 diaktifkan")
			RSLonlyflag = 3

		} else {
			// println("Poin Centang 1 dinonaktifkan")
			RSLonlyflag = 0
		}
	})
	RSLonlyLabel.Hide()
	check1.Hide()
	check3.Hide()
	check4.Hide()

	Updateonlyflag = 0
	UpdateonlyLabel := widget.NewLabel("Update Only:")
	check2 := widget.NewCheck("Update Only", func(checked bool) {
		if checked {
			// println("Poin Centang 1 diaktifkan")
			Updateonlyflag = 1

		} else {
			// println("Poin Centang 1 dinonaktifkan")
			Updateonlyflag = 0
		}
	})
	UpdateonlyLabel.Hide()
	check2.Hide()
	// Membuat container horizontal untuk Hardware Model (High Set & Low Set)
	hardwareRow := container.NewHBox(
		hardwareLabel,
		hardwareRadioGroup,
	)
	typemotorRow := container.NewHBox(
		typemotorLabel,
		typemotordropdown,
	)

	// Mengubah fungsi dropdown agar menampilkan input ketika "bms" dipilih
	dropdown.OnChanged = func(selected string) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic in dropdown:", r)
			}
		}()

		if selected == "" {
			dialog.ShowInformation("Error", "Pilihan tidak boleh kosong", MyWindow)
			return
		}

		MyWindow.SetTitle(namaapp + " " + selected)
		modelidLabel.Hide()
		modelidEntry.Hide()
		modelversion.Hide()
		parallelsLabel.Hide()
		parallelsEntry.Hide()
		batteryLabel.Hide()
		batteryMonthDropdown.Hide()
		batteryYearDropdown.Hide()
		cellLabel.Hide()
		cellMonthDropdown.Hide()
		cellYearDropdown.Hide()
		cellbrandLabel.Hide()
		cellbrandDropdown.Hide()
		cellTypeLabel.Hide()
		cellTypeDropdown.Hide()
		teganganLabel.Hide()
		teganganDropdown.Hide()
		hardwareLabel.Hide()
		hardwareRadioGroup.Hide()
		typemotorLabel.Hide()
		typemotordropdown.Hide()
		TypeMCULabel.Hide()
		mcudropdown.Hide()
		typelcdLabel.Hide()
		typelcddropdown.Hide()
		RSLonlyLabel.Hide()
		check1.Hide()
		check3.Hide()
		check4.Hide()
		UpdateonlyLabel.Hide()
		check2.Hide()
		if selected == "bms" {
			modelver, erri := GetlistModelversi(selected)
			if erri != nil {
				log.Fatal(err)
			}
			modelversion.Options = modelver
			modelversion.Refresh() // Refresh untuk memperbarui pilihan
			modelversion.Show()
			modelidLabel.Show()
			// modelidEntry.Show()
			parallelsLabel.Show()
			parallelsEntry.Show()
			cellLabel.Show()
			cellMonthDropdown.Show()
			cellYearDropdown.Show()
			cellbrandLabel.Show()
			cellbrandDropdown.Show()
			teganganLabel.Show()
			teganganDropdown.Show()
			hardwareLabel.Show()
			hardwareRadioGroup.Show()
			UpdateonlyLabel.Show()
			check2.Show()
		} else if selected == "vcu" {
			modelver, erri := GetlistModelversi(selected)
			if erri != nil {
				log.Fatal(err)
			}
			mcu, err := readExcelData("menu.xlsx", "vcumcu")
			if err != nil {
				log.Fatal(err)
			}
			modelversion.Options = modelver
			modelversion.Refresh() // Refresh untuk memperbarui pilihan
			modelversion.Show()
			mcudropdown.Options = mcu
			mcudropdown.Refresh() // Refresh untuk memperbarui pilihan
			mcudropdown.Show()
			TypeMCULabel.Show()
			modelidLabel.Show()
			// modelidEntry.Show()
			typemotorLabel.Show()
			typemotordropdown.Show()
			RSLonlyLabel.Show()
			check1.Show()
			check3.Show()
			check4.Show()

		} else if selected == "hmi" {
			modelver, erri := GetlistModelversi(selected)
			if erri != nil {
				log.Fatal(err)
			}
			modelversion.Options = modelver
			modelversion.Refresh() // Refresh untuk memperbarui pilihan
			modelversion.Show()
			mcu, err := readExcelData("menu.xlsx", "hmimcu")
			if err != nil {
				log.Fatal(err)
			}
			mcudropdown.Options = mcu
			mcudropdown.Refresh() // Refresh untuk memperbarui pilihan
			mcudropdown.Show()
			TypeMCULabel.Show()
			modelidLabel.Show()
			// modelidEntry.Show()
			typemotorLabel.Show()
			typemotordropdown.Show()

			typelcdLabel.Show()
			typelcddropdown.Show()
		} else if selected == "keyless" {
			mcu, err := readExcelData("menu.xlsx", "keylessmcu")
			if err != nil {
				log.Fatal(err)
			}
			mcudropdown.Options = mcu
			mcudropdown.Refresh() // Refresh untuk memperbarui pilihan
			mcudropdown.Show()
			modelver, erri := GetlistModelversi(selected)
			if erri != nil {
				log.Fatal(err)
			}
			modelversion.Options = modelver
			modelversion.Refresh() // Refresh untuk memperbarui pilihan
			modelversion.Show()
			TypeMCULabel.Show()
			modelidLabel.Show()
			// modelidEntry.Show()
			typemotorLabel.Show()
			typemotordropdown.Show()

		}
	}

	// Menambahkan tombol Submit
	submitButton := widget.NewButton("Program", func() {
		flashon := 0
		if dropdown.Selected == "bms" {
			// Validasi jika SN atau input lain kosong
			if parallelsEntry.Text == "" || cellMonthDropdown.Selected == "" || cellYearDropdown.Selected == "" || cellTypeShort == "" {
				dialog.ShowInformation("Error", "Semua field harus diisi!", MyWindow)

				// Data.Bord = dropdown.Selected
				// ProsesFlash()

			} else {
				Data.Num.Tegangan = teganganShort
				Data.Num.ParalelN = parallelsEntry.Text
				Data.Num.Jenis = cellbrandShort
				Data.Num.Type_ = cellTypeShort
				Data.Num.Bulan_ex = cellMonthDropdown.Selected
				Data.Num.Tahun_ex = cellYearDropdown.Selected
				Data.Id = hardwareModelValue

				Data.Bord = dropdown.Selected
				// dialog.ShowInformation("Info", "WAIT TO FLASH ", MyWindow)

				// ProsesFlash()
				flashon = 1
			}
		} else if dropdown.Selected == "vcu" {
			// Validasi jika SN atau input lain kosong
			if mcudropdown.Selected == "" || typemotordropdown.Selected == "" {
				dialog.ShowInformation("Error", "Semua field harus diisi!", MyWindow)

			} else {

				Data.Bord = dropdown.Selected
				// dialog.ShowInformation("Info", "WAIT TO FLASH ", MyWindow)
				// ProsesFlash()
				flashon = 1

				// fmt.Println("Data 1:", joy["application"].(string))
				// fmt.Println("Data 2:", joy["bootloader"].(string))
				// fmt.Println("Data 3:", joy["ble_application"].(string))
				// fmt.Println("Data 4:", joy["ble_bootloader"].(string))

			}
		} else if dropdown.Selected == "hmi" {
			// Validasi jika SN atau input lain kosong
			if mcudropdown.Selected == "" || typelcddropdown.Selected == "" || typemotordropdown.Selected == "" {
				dialog.ShowInformation("Error", "Semua field harus diisi!", MyWindow)

			} else {

				Data.Bord = dropdown.Selected
				// dialog.ShowInformation("Info", "WAIT TO FLASH ", MyWindow)
				// ProsesFlash()
				flashon = 1
			}
		} else if dropdown.Selected == "keyless" {
			// Validasi jika SN atau input lain kosong
			if mcudropdown.Selected == "" || typemotordropdown.Selected == "" {
				dialog.ShowInformation("Error", "Semua field harus diisi!", MyWindow)

			} else {

				Data.Bord = dropdown.Selected
				flashon = 1
				// dialog.ShowInformation("Info", "WAIT TO FLASH ", MyWindow)
				// ProsesFlash()
			}
		} else {
			dialog.ShowInformation("Info", "Pilihan: "+dropdown.Selected, MyWindow)
			flashon = 0
		}

		if flashon == 1 {
			entry := widget.NewEntry()
			IPc := widget.NewEntry()

			resultLabel := widget.NewLabel("Klik tombol untuk mulai scan QR")

			butQRcode := widget.NewButton("Scan QR dari Kamera", func() {
				resultLabel.SetText("Scanning...")

				go func() {
					IPcame := IPc.Text
					cameraURL := "http://" + IPcame + ":8080/shot.jpg"
					text, err := ScanQRCodeFromURL(cameraURL)
					if err != nil {
						resultLabel.SetText("Gagal: " + err.Error())
					} else {

						entry.SetText(text)
						resultLabel.SetText("Hasil QR: " + text)
						dialog.ShowInformation("QR Code Ditemukan", text, MyWindow)
					}
				}()
			})
			content := container.NewVBox(
				widget.NewLabel("Masukkan QRCODE:"),
				entry,
				widget.NewLabel("Masukkan IPcamera:"),
				IPc,
				resultLabel,
				butQRcode,
			)

			dialog.ShowCustomConfirm("Input QRCODE", "OK", "Cancel", content, func(b bool) {
				if b {
					QRcode = entry.Text
					// label.SetText(fmt.Sprintf("Ini adalah namanya: %s", name))
					fmt.Println("Selected: %s", QRcode)
					dialog.ShowInformation("Info", "WAIT TO FLASH "+QRcode, MyWindow)
					ProsesFlash()
				}
			}, MyWindow)
		}

	})

	// Susun dropdown bulan dan tahun di satu baris untuk Battery
	batteryRow := container.NewHBox(
		batteryMonthDropdown,
		batteryYearDropdown,
	)

	// Susun dropdown bulan dan tahun di satu baris untuk Battery
	mcuRow := container.NewHBox(
		TypeMCULabel,
		mcudropdown,
	)

	// Susun dropdown bulan dan tahun di satu baris untuk Battery
	lcdRow := container.NewHBox(
		typelcdLabel,
		typelcddropdown,
	)

	// Susun dropdown bulan dan tahun di satu baris untuk Battery
	rslvcuRow := container.NewHBox(
		RSLonlyLabel,
		check1,
		check3,
		check4,
	)

	// Susun dropdown bulan dan tahun di satu baris untuk Cell
	cellRow := container.NewHBox(
		cellMonthDropdown,
		cellYearDropdown,
	)
	jeniscellRow := container.NewHBox(
		cellbrandLabel,
		cellbrandDropdown,
		cellTypeLabel,
		cellTypeDropdown,
	)
	settingBattryRow := container.NewHBox(
		teganganLabel,
		teganganDropdown,
		parallelsLabel,
		parallelsEntry,
	)

	// Susun konten dalam tampilan vertikal
	content := container.NewVBox(
		widget.NewLabel("Pilih opsi:"),
		dropdown,
		modelidLabel,
		modelidEntry,
		modelversion,
		settingBattryRow,
		jeniscellRow,
		cellLabel,
		cellRow,
		mcuRow,
		lcdRow,
		typemotorRow,
		batteryLabel,
		batteryRow,
		hardwareRow,
		rslvcuRow,
		UpdateonlyLabel,
		check2,
		submitButton,
		widget.NewButton("Quit", func() {
			myApp.Quit()
		}),
	)

	// Mengatur konten ke dalam jendela
	MyWindow.SetContent(content)

	// Menampilkan jendela
	MyWindow.Show()
}

func Dialoginfo(msg string) {
	dialog.ShowInformation("Info", msg, MyWindow)
}

func Dialogeror(msg string) {
	dialog.ShowInformation("Error", msg, MyWindow)
}
