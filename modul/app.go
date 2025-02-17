package modul

import (
	"fmt"
	"strings"
)

var Data Bus

func Ceking_firmware(bord string) uint8 {
	// versi := CekVersion(bord)
	key, value := CekVersion(Data.Bord)

	result := CheckBinFolder(bord)
	if result == 0 { // Menambahkan tanda kurung buka setelah if

		Download_firmware(bord, key, value)
		// _, versi := Verifikasi_versi(bord)
		WriteVersionToFile(key, value, bord)
		// hapusFileFirmware()
		Setversifirmware(bord)
		return 0
		// Ekstrak_move("firmware.rar", "bin\\vcu\\")
		// Implementasi jika result == 0
	} else {
		versionData, err := ParseVersionFile(bord)
		if err != nil {
			fmt.Println("Gagal memparsing file versi.txt:", err)
			return 1
		}
		for keyold, valueold := range versionData {
			for i := 0; i < len(key); i++ {
				if keyold == key[i] {
					// fmt.Printf("old %s new %s", valueold, value[i])
					if CompareVersions(valueold, value[i]) {
						HapusSemuaFileDalamFolder(bord)
						Download_firmware(bord, key, value)
						// _, versi := Verifikasi_versi(bord)
						WriteVersionToFile(key, value, bord)
						Setversifirmware(bord)
						return 0
					}

				}

			}
		}
		Setversifirmware(bord)
		return 0

	}

	// return 1
}

func Programflash(bord string) {
	if strings.Contains(bord, "vcu") {
		// ExecuteOpenOCDvcu(Data)
		// SNVCU(Data)
		if RSLonlyflag == 1 {
			ExecuteOpenOCDble(Data)
		} else {
			if ExecuteOpenOCDble(Data) == nil {
				ExecuteOpenOCDvcu(Data)
			}
		}
		// Implementasi untuk board yang mengandung "vcu"
	} else if strings.Contains(bord, "hmi") {
		ExecuteOpenOCDhmi(Data)
		// SNhmi(Data)
		// Implementasi untuk board yang mengandung "hmi"
	} else if strings.Contains(bord, "bms") {
		ExecuteOpenOCDBMS(Data)
		// buatexsel(bord)
		// tambahbaris(bord, "ayambakar")
		// SnBmsString(Data, 2)
		// SnBms(Data)
		// updateCount("bmsh", 1)
		// getcount("bmsh")
		// timecek()
		// Implementasi untuk board yang mengandung "bms"
	} else if strings.Contains(bord, "keyless") {
		// ProgramKeyles()
		ExecuteOpenOCDkeyless(Data)
		// Implementasi untuk board yang mengandung "keyfob"
	}
}

func ProsesFlash() {

	if Ceking_firmware(Data.Bord) == 0 {
		Programflash(Data.Bord)
	}
}
