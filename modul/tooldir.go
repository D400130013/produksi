package modul

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nwaples/rardecode"
)

func CheckBinFolder(bord string) int {
	// Tentukan path ke folder bin
	folderPath := fmt.Sprintf("./bin/%s/%d/", bord, Data.Model)
	// Cek apakah folder bin ada
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		fmt.Println("Error membuat folder 'bin':", err)
		return 0
	}

	// Membaca isi folder bin
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		fmt.Println("Error membaca folder 'bin':", err)
		return 0
	}

	// Jika folder ada tetapi kosong
	if len(files) == 0 {
		fmt.Println("Folder 'bin' ada tapi kosong.")
		return 0
	}

	// Jika folder ada dan berisi file
	fmt.Println("Isi folder 'bin':")
	for _, file := range files {
		fmt.Println(file.Name())
	}

	// Mengembalikan 1 jika ada file
	return 1
}

func WriteVersionToFile(keys []string, values []string, bords string) error {
	// Membuka file untuk ditulis
	filename := fmt.Sprintf("./bin/%s/%d/versi.txt", bords, Data.Model)
	// filename := "./bin/" + bords + "/versi.txt"
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("gagal membuat file: %v", err)
	}
	defer file.Close()

	// Menulis kunci dan nilai ke dalam file
	for i := 0; i < len(keys); i++ {
		_, err := fmt.Fprintf(file, "%s: %s\n", keys[i], values[i])
		if err != nil {
			return fmt.Errorf("gagal menulis ke file: %v", err)
		}
	}

	fmt.Printf("Data berhasil ditulis ke %s\n", filename)
	return nil
}

func Setversifirmware(bords string) error {
	// filename := "./bin/" + bords + "/versi.txt"
	versionData, err := ParseVersionFile(bords)
	if err != nil {
		return fmt.Errorf("gagal memparsing file: %v", err)
	}

	if appStr, exists := versionData["application"]; exists {
		Data.Versifirmapp, _ = VersionToHex(appStr)

		// Sekarang Anda dapat menggunakan bootloaderUint32
		fmt.Printf("application sebagai uint32: %d\n", Data.Versifirmapp)
	} else {
		return fmt.Errorf("kunci bootloader tidak ditemukan")
	}

	if bootloaderStr, exists := versionData["bootloader"]; exists {
		Data.Versifirmboot, _ = VersionToHex(bootloaderStr)

		// Sekarang Anda dapat menggunakan bootloaderUint32
		fmt.Printf("Bootloader sebagai uint32: %d\n", Data.Versifirmboot)
	} else {
		return fmt.Errorf("kunci bootloader tidak ditemukan")
	}

	return nil
}

// WriteVersionToFile creates versi.txt and writes the version information to it.
// func WriteVersionToFile(bords string, content string) error {
// 	// Buka atau buat file versi.txt
// 	filename := "./bin/" + bords + "/versi.txt"
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return fmt.Errorf("failed to create file: %w", err)
// 	}
// 	defer file.Close()

// 	// Tulis konten ke dalam file
// 	_, err = file.WriteString(content)
// 	if err != nil {
// 		return fmt.Errorf("failed to write to file: %w", err)
// 	}

// 	return nil
// }

// FetchFileFromURL downloads and returns the contents of the file at the specified URL.
func FetchFileFromURL(url string) (string, error) {
	// Lakukan HTTP GET untuk mendapatkan isi file
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	// Periksa jika HTTP request berhasil
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch file, status code: %d", resp.StatusCode)
	}

	// Baca isi body dari response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Kembalikan isi file sebagai string
	return string(body), nil
}

func Get_versi(bords string) (int, string) {

	// URL dari file yang ingin diambil
	url := "https://firmware.dev.savart-ev.com/bin/" + bords + "/versi.txt"
	contents, err := FetchFileFromURL(url)
	if err != nil {
		fmt.Println("Error:", err)
		return 1, ""
	}

	lines := strings.Split(contents, "\n")
	if len(lines) > 0 {
		// content := strings.Join(lines, "\n")
		// content := fmt.Sprintf("%s", lines)
		// WriteVersionToFile("./bin/"+bords+"/versi.txt", lines[0])
		return 0, lines[0]
	} else {
		fmt.Println("Contents are empty")
	}
	return 1, ""
}

func extractRar(rarFile, destDir string) error {
	// Buka file RAR
	file, err := os.Open(rarFile)
	if err != nil {
		return fmt.Errorf("error membuka file RAR: %v", err)
	}
	defer file.Close()

	// Membuat pembaca RAR
	rarReader, err := rardecode.NewReader(file, "")
	if err != nil {
		return fmt.Errorf("error membuat pembaca RAR: %v", err)
	}

	// Iterasi melalui setiap file dalam arsip
	for {
		// Baca header file berikutnya
		header, err := rarReader.Next()
		if err == io.EOF {
			break // Selesai membaca arsip
		}
		if err != nil {
			return fmt.Errorf("error membaca header file: %v", err)
		}

		// Tentukan path tujuan
		destPath := filepath.Join(destDir, header.Name)

		// Buat direktori jika perlu
		if header.IsDir {
			err = os.MkdirAll(destPath, os.ModePerm) // Pastikan direktori dibuat
			if err != nil {
				return fmt.Errorf("error membuat direktori: %v", err)
			}
			continue
		}

		// Pastikan direktori tujuan ada sebelum membuat file
		err = os.MkdirAll(filepath.Dir(destPath), os.ModePerm) // Membuat direktori jika belum ada
		if err != nil {
			return fmt.Errorf("error membuat direktori untuk file tujuan: %v", err)
		}

		// Buka file tujuan
		outFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("error membuat file tujuan: %v", err)
		}

		// Salin konten dari arsip ke file tujuan
		_, err = io.Copy(outFile, rarReader)
		if err != nil {
			outFile.Close()
			return fmt.Errorf("error menyalin file: %v", err)
		}

		// Tutup file tujuan setelah selesai menyalin
		outFile.Close()
	}

	return nil
}

func Ekstrak_move(rarFile string, destDir string) {
	// rarFile = "archive.rar"
	// destDir = "./bin"

	// Pastikan folder tujuan ada
	err := os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error membuat folder tujuan:", err)
		return
	}

	// Ekstrak file RAR
	err = extractRar(rarFile, destDir)
	if err != nil {
		fmt.Println("Error mengekstrak file RAR:", err)
	} else {
		fmt.Println("File RAR berhasil diekstrak ke", destDir)
	}
}

// Fungsi untuk memparsing data dari versi.txt
func ParseVersionFile(bords string) (map[string]string, error) {
	// Membuka file untuk dibaca
	// filename := "./bin/" + bords + "/versi.txt"
	filename := fmt.Sprintf("./bin/%s/%d/versi.txt", bords, Data.Model)
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file: %v", err)
	}
	defer file.Close()

	// Membuat map untuk menyimpan kunci dan nilai
	versionData := make(map[string]string)

	// Membaca file baris per baris
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Memisahkan kunci dan nilai berdasarkan ": "
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			versionData[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error membaca file: %v", err)
	}

	return versionData, nil
}

func Verifikasi_versi(bords string) (int, string) {
	erro, versi_update := Get_versi(bords)
	if erro != 0 {
		fmt.Printf("Gagal membuka Get_versi: %v\n", erro)
		return 2, ""
	}
	// filename := "./bin/" + bords + "/versi.txt"
	filename := fmt.Sprintf("./bin/%s/%d/versi.txt", bords, Data.Model)
	_, err := os.Stat(filename)
	// Jika tidak ada error, file ada
	if err != nil {
		return 1, versi_update
	}

	// Buka file untuk dibaca
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Gagal membuka file: %v\n", err)
		return 2, ""
	}
	defer file.Close()

	// Baca baris pertama dari file
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		firstLine := scanner.Text() // Baris pertama
		// fmt.Printf("versi file: |%s|%s|\n", firstLine, versi_update)
		if strings.Contains(versi_update, firstLine) {
			// fmt.Println("string sesuai")
			return 0, ""
		} else {
			return 1, versi_update
		}
		// Kembalikan 0 (sukses) dan baris pertama
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error membaca file: %v\n", err)
		return 2, ""
	}

	// Jika tidak ada baris pertama, kembalikan error

	return 2, ""
}

func CompareVersions(oldVersion, newVersion string) bool {
	// Menghapus awalan "v" jika ada
	oldVersion = strings.TrimPrefix(oldVersion, "v")
	newVersion = strings.TrimPrefix(newVersion, "v")

	// Memisahkan string versi menjadi komponen angka
	oldParts := strings.Split(oldVersion, ".")
	newParts := strings.Split(newVersion, ".")

	// Looping untuk membandingkan setiap bagian versi
	for i := 0; i < len(oldParts); i++ {
		// Mengonversi setiap bagian versi ke integer untuk dibandingkan
		old, _ := strconv.Atoi(oldParts[i])
		new, _ := strconv.Atoi(newParts[i])

		if new > old {
			fmt.Printf("Versi baru %s lebih tinggi dari versi lama %s\n", newVersion, oldVersion)
			return true
		} else if new < old {
			fmt.Printf("Versi lama %s lebih tinggi dari versi baru %s\n", oldVersion, newVersion)
			return false
		}
	}

	// Jika semua bagian sama
	fmt.Println("Versi baru dan lama sama.")
	return false
}

func HapusSemuaFileDalamFolder(bord string) {
	// folderPath := "./bin/" + bord + "/" // Path ke folder yang ingin dihapus file-filenya
	folderPath := fmt.Sprintf("./bin/%s/%d/", bord, Data.Model)
	// Baca isi folder
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		fmt.Printf("Gagal membaca folder: %v\n", err)
		return
	}

	// Loop melalui semua file di dalam folder
	for _, file := range files {
		filePath := filepath.Join(folderPath, file.Name())

		// Hapus file atau folder
		if file.IsDir() {
			// Jika file adalah folder, hapus beserta isinya
			err = os.RemoveAll(filePath)
		} else {
			// Jika file biasa, hapus file
			err = os.Remove(filePath)
		}

		if err != nil {
			fmt.Printf("Gagal menghapus %s: %v\n", filePath, err)
		} else {
			fmt.Printf("%s berhasil dihapus.\n", filePath)
		}
	}
}

func hapusFileFirmware() error {
	filePath := "firmware.rar"
	err := os.Remove(filePath)
	if err != nil {
		// Jika terjadi error saat menghapus file, kembalikan error
		return err
	}
	// File berhasil dihapus
	return nil
}

func GetFirstLetter(word string) string {
	if len(word) == 0 {
		return ""
	}
	return string(word[0])
}

// Function to convert version string to hexadecimal representation
func VersionToHex(version string) (uint32, error) {
	versi := GetStringAfterV(version)
	// Split the version string into its components
	parts := strings.Split(versi, ".")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid version format")
	}

	// Parse the components
	y, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	x, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	z, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, err
	}

	// Combine the hexadecimal components
	hexVersion := (uint32(z)<<24 | (uint32(x)&0x00FF)<<16 | (uint32(x) & 0xFF00) | uint32(y))

	return hexVersion, nil
}

// Fungsi untuk memastikan folder ada
func ensureFolderExists(folderPath string) error {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, os.ModePerm) // Membuat folder
		if err != nil {
			return fmt.Errorf("gagal membuat folder: %w", err)
		}
	}
	return nil
}

func rearrangeBytes(input uint32) uint32 {
	// Pisahkan masing-masing byte dari uint32 (ABCD)
	byteA := (input >> 24) & 0xFF // Byte A (8 bit paling atas)
	byteB := (input >> 16) & 0xFF // Byte B (8 bit kedua)
	byteC := (input >> 8) & 0xFF  // Byte C (8 bit ketiga)
	byteD := input & 0xFF         // Byte D (8 bit paling bawah)

	// Gabungkan menjadi BACD (B menjadi byte pertama, A menjadi byte kedua, C dan D tetap sama)
	result := (byteB << 24) | (byteA << 16) | (byteC << 8) | byteD

	return result
}

// Fungsi untuk mengambil semua string sebelum titik
func getStringBeforeDot(s string) string {
	// Memisahkan string menggunakan .
	parts := strings.Split(s, ".")
	// Mengembalikan bagian pertama sebelum titik
	return parts[0]
}

// Fungsi untuk mengambil semua string setelah 'v'
func GetStringAfterV(s string) string {
	// Mencari indeks dari 'v' dalam string
	vIndex := strings.Index(s, "v")
	// Mengembalikan bagian string setelah 'v'
	return s[vIndex+1:]
}

func MoveFileToFolder(sourcePath string, destinationFolder string) error {

	// Membuat path tujuan
	destinationPath := filepath.Join(destinationFolder, filepath.Base(sourcePath))

	// Memindahkan file
	err := os.Rename(sourcePath, destinationPath)
	if err != nil {
		return fmt.Errorf("gagal memindahkan file: %v", err)
	}

	fmt.Printf("File berhasil dipindahkan ke: %s\n", destinationPath)
	return nil
}

// createBinFile membuat file biner (.bin) dan menulis data ke dalamnya.
func CreateBinFile(filename string, data []byte) error {
	// Membuka file dalam mode tulis dan biner, membuat file jika belum ada.
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("gagal membuat file: %w", err)
	}
	defer file.Close()

	// Menulis data ke dalam file
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("gagal menulis ke file: %w", err)
	}

	return nil
}
