package modul

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type SN struct {
	MCU      string
	Tegangan string
	ParalelN string
	Jenis    string
	Type_    string
	Bulan_pb string //bord produksi
	Tahun_pb string
	Bulan_ex string //external produksi
	Tahun_ex string
}

type Bus struct {
	Bord              string
	Id                uint32
	Model             uint32
	Modelver          string
	Versifirmapp      uint32
	Versifirmboot     uint32
	Versifirmapp2     uint32
	Versifirmboot2    uint32
	Versifirmappstr   string
	Versifirmbootstr  string
	Versifirmappstr2  string
	Versifirmbootstr2 string
	Num               SN
}

func DownloadFile(url string, filepath string) (uint32, uint32, error) {
	// Membuat request HTTP GET
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	// fmt.Println("Header Respons:", resp.Status)
	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("gagal mengambil metadata, status: %s", resp.Status)
	}
	var crc uint32
	var length uint32

	if crcStr, ok := resp.Header["Firmware-Crc"]; ok {
		if crcValue, err := strconv.ParseUint(crcStr[0], 10, 32); err == nil {
			crc = uint32(crcValue)
			// fmt.Printf("Firmware-Crc: %d\n", crc)
		} else {
			fmt.Println("Gagal mengonversi Firmware-Crc:", err)
		}
	}

	if lengthStr, ok := resp.Header["Firmware-Length"]; ok {
		if lengthValue, err := strconv.ParseUint(lengthStr[0], 10, 32); err == nil {
			length = uint32(lengthValue)
			// fmt.Printf("Firmware-Length: %d\n", length)
		} else {
			fmt.Println("Gagal mengonversi Firmware-Length:", err)
		}
	}

	// Membuka file untuk menyimpan data yang diunduh
	out, err := os.Create(filepath)
	if err != nil {
		return 0, 0, err
	}
	defer out.Close()

	// Menyalin konten dari response ke file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return 0, 0, err
	}

	return crc, length, nil
}
func Download_firmware(bord string, key []string, value []string) {
	var url string
	var filepath string
	var firmwareType string
	for i := 0; i < len(key); i++ {
		switch key[i] {
		case "ble_application":
			firmwareType = "ble-application"
		case "ble_bootloader":
			firmwareType = "ble-bootloader"
		default:
			firmwareType = key[i]
		}

		url = "https://fota.savart-ev.com/" + bord + "/" + firmwareType + "/download?m=v" + Data.Modelver + "&v=" + value[i]
		filepath = firmwareType + ".bin"
		crc, length, err := DownloadFile(url, filepath)
		if err != nil {
			println("Gagal mengunduh firmware:", err)
		} else {

			crcbin, lengthbin, eror := CalculateCrcAndLen(filepath)

			if eror == nil {
				if crc == crcbin && length == lengthbin {
					println("Firmware berhasil diunduh ", key[i])
					// Setelah mengunduh firmware
					savebin := fmt.Sprintf("./bin/%s/%d/", bord, Data.Model)
					err = MoveFileToFolder(filepath, savebin)
					if err != nil {
						println("Gagal memindahkan file:", err)
					}
				} else {
					println("CRC salah cek CRC", crc, crcbin, length, lengthbin)
				}
			}

		}
	}

	// else {
	// 	Ekstrak_move(filepath, "./bin/"+bord+"/")
	// }
	// Mencetak header balasan

}

func GetlistModelversi(bords string) ([]string, error) {
	metadataURL := "https://fota.savart-ev.com/" + bords + "/model-version"
	resp, err := http.Get(metadataURL)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil metadata: %v", err)
	}
	defer resp.Body.Close()

	// Membaca isi respons
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca respons: %v", err)
	}

	// Mengubah JSON menjadi struct untuk mengakses data version
	var data struct {
		Code string
		Data []struct {
			ID      int    `json:"id"`
			Type    string `json:"type"`
			Version string `json:"version"`
		} `json:"data"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("gagal mengurai JSON: %v", err)
	}

	// Mengumpulkan versi ke dalam array string
	var modelVersi []string
	for _, item := range data.Data {
		modelVersi = append(modelVersi, item.Version)
	}

	return modelVersi, nil
}

func CekVersion(bords string) ([]string, []string) {
	metadataURL := "https://fota.savart-ev.com/" + bords + "/model-version/v" + Data.Modelver
	resp, err := http.Get(metadataURL)
	if err != nil {
		fmt.Printf("gagal mengambil metadata: %v\n", err)
		return nil, nil
	}
	defer resp.Body.Close()

	// Membaca isi respons
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("gagal membaca respons: %v\n", err)
		return nil, nil
	}

	// Mengubah JSON menjadi struct untuk mengakses data version
	var data struct {
		Code string
		Data struct {
			ID            int                    `json:"id"`
			Type          string                 `json:"type"`
			Version       string                 `json:"version"`
			LatestVersion map[string]interface{} `json:"latest_version"`
		} `json:"data"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Printf("gagal mengurai JSON: %v\n", err)
		return nil, nil
	}

	// Membuat array untuk kunci dan nilai
	var keys []string
	var values []string
	for key, value := range data.Data.LatestVersion {
		keys = append(keys, key)
		values = append(values, fmt.Sprintf("%v", value)) // Mengonversi nilai ke string
		fmt.Printf("Kunci: %s, Nilai: %v\n", key, value)  // Mencetak kunci dan nilai
	}

	return keys, values
}

func IpaLogin(username string, password string) uint8 {
	// Menggabungkan username dan password dengan tanda ":"
	credentials := fmt.Sprintf("%s:%s", username, password)

	// Melakukan encoding ke Base64
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))

	// Membuat header Authorization
	authorizationHeader := fmt.Sprintf("Basic %s", encodedCredentials)

	// URL endpoint
	url := "https://team.savart-ev.com/api/login"

	// Membuat request GET
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
		return 1
	}

	// Menambahkan header Authorization ke request
	req.Header.Set("Authorization", authorizationHeader)

	// Membuat client HTTP
	client := &http.Client{}

	// Mengirim request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
		return 1
	}
	defer resp.Body.Close()

	// Membaca respons dari server
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
		return 1
	}

	// Menampilkan status dan body respons
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Response Body: %s\n", string(body))
	if resp.StatusCode == 200 {
		ParseResponseBody(string(body))
		return 0
	}
	return 1
}

func Getprofil(data string) error {
	// Membuat header Authorization
	authorizationHeader := fmt.Sprintf("Bearer %s", data)
	fmt.Printf("{%s}", authorizationHeader)
	// URL endpoint
	url := "https://team.savart-ev.com/api/auth/profile"

	// Membuat request GET
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("gagal parsing JSON: %v", err)

	}
	req.Header.Set("Content-Type", "Application/json")
	// Menambahkan header Authorization ke request
	req.Header.Set("Authorization", authorizationHeader)

	// Membuat client HTTP
	client := &http.Client{}

	// Mengirim request
	resp, err := client.Do(req)
	if err != nil {

		return fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Membaca respons dari server
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		return fmt.Errorf("Error sending request: %v", err)
	}

	// Menampilkan status dan body respons
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Response Body: %s\n", string(body))
	if resp.StatusCode == 200 {

		return nil
	}
	return fmt.Errorf("Error sending request: %v", err)
}
