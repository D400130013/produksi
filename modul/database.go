package modul

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Data_vcu struct {
	Mcu           string `json:"mcu"`             // ST : STM32, MT: "Mobiletek"
	VehicleTypeID int    `json:"vehicle_type_id"` // 1 : s1p, 2 : s1, 3: a
	ModelVersion  string `json:"model_version"`
}

type Data_bms struct {
	Voltage       int    `json:"voltage"`         // Tegangan baterai dalam volt
	ParallelNum   int    `json:"parallel_num"`    // Jumlah paralel baterai
	CellBrand     string `json:"cell_brand"`      // Merek baterai, contoh: E untuk EVE, P untuk Panasonic
	CellBrandType string `json:"cell_brand_type"` // Tipe baterai, contoh: A, B, C, dst.
	CellProdYear  int    `json:"cell_prod_year"`  // Tahun produksi baterai
	CellProdMonth int    `json:"cell_prod_month"` // Bulan produksi baterai, 1 untuk Januari, 12 untuk Desember
	ModelVersion  string `json:"model_version"`
	BatteryType   string `json:"battery_type"` // Tipe baterai, contoh: "hight"
}

type Data_hmi struct {
	Mcu           string `json:"mcu"`             // Tipe MCU, contoh: STM32 atau Mobiletek
	LcdTypeId     string `json:"lcd_type_id"`     // Tipe LCD, contoh: S
	VehicleTypeId int    `json:"vehicle_type_id"` // Tipe Kendaraan, contoh: 1 untuk s1p, 2 untuk s1, 3 untuk a
	ModelVersion  string `json:"model_version"`   // Versi Model, contoh: 1.1
}

type Data_keyless struct {
	Mcu           string `json:"mcu"`             // Tipe MCU, contoh: STM32 atau Mobiletek
	VehicleTypeId int    `json:"vehicle_type_id"` // Tipe Kendaraan, contoh: 1 untuk s1p, 2 untuk s1, 3 untuk a
	ModelVersion  string `json:"model_version"`   // Versi Model, contoh: 1.1
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "server_produksi"
)

var db *sql.DB

var IP string = "https://part.savart-ev.com/api"

// var IP string = "https://part.dev.savart-ev.com/api"

// var PORT string = ""

func App_DB_connection() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Getcount(bord string) uint32 {
	db, err := App_DB_connection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Menyiapkan pernyataan SQL untuk pembaruan
	var count int
	query := "SELECT coun FROM bords WHERE nama = $1"
	err = db.QueryRow(query, bord).Scan(&count)
	if err == sql.ErrNoRows {
		log.Printf("Tidak ada data ditemukan untuk nama: %s\n", bord)
		return 0 // Atau nilai default lainnya
	} else if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Count untuk nama '%s' adalah: %d\n", bord, count)
	return uint32(count)
}

func UpdateCount(bord string, newCount int) error {
	db, err := App_DB_connection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Menyiapkan pernyataan SQL untuk pembaruan
	query := "UPDATE bords SET coun = $1 WHERE nama = $2"
	_, err = db.Exec(query, newCount, bord)
	if err != nil {
		return err
	}

	return nil
}

func Gettgl(bord string) int {
	db, err := App_DB_connection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Menyiapkan pernyataan SQL untuk pembaruan
	var count int
	query := "SELECT tgl FROM bords WHERE nama = $1"
	err = db.QueryRow(query, bord).Scan(&count)
	if err == sql.ErrNoRows {
		log.Printf("Tidak ada data ditemukan untuk nama: %s\n", bord)
		return 0 // Atau nilai default lainnya
	} else if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Tanggal untuk nama '%s' adalah: %d\n", bord, count)
	return int(count)
}

func Updatetgl(bord string, newCount int) error {
	db, err := App_DB_connection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Menyiapkan pernyataan SQL untuk pembaruan
	query := "UPDATE bords SET tgl = $1 WHERE nama = $2"
	_, err = db.Exec(query, newCount, bord)
	if err != nil {
		return err
	}

	return nil
}

//rest api

func UpdateVcu(data Data_vcu) map[string]interface{} {
	// Konversi struct menjadi JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil
	}

	// Membuat request POST dengan JSON
	url := IP + "/vcu/generate-sn"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	// Menambahkan header untuk tipe konten JSON
	req.Header.Set("Content-Type", "application/json")

	// Mengirim request menggunakan http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil
	}
	defer resp.Body.Close()

	// Menampilkan response dari server
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte
	fmt.Println(string(body))              // convert to string before print
	var result map[string]interface{}

	// Decode JSON menjadi map
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	// Menampilkan isi map
	fmt.Println("Code:", result["code"])

	// Mengakses nested data
	datajson, ok := result["data"].(map[string]interface{})
	if !ok {
		fmt.Println("Error: data tidak dapat dikonversi menjadi map[string]interface{}")
		return nil
	}
	return datajson // convert to string before print
}

// type BmsGenerateSNResp struct {
//     SN string `json:"sn"`
//     Detail struct {
// 		Counter
//     }
// }

func UpdateBms(data Data_bms) map[string]interface{} {
	// Konversi struct menjadi JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil
	}

	// Membuat request POST dengan JSON
	url := IP + "/bms/generate-sn"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	// Menambahkan header untuk tipe konten JSON
	req.Header.Set("Content-Type", "application/json")

	// Mengirim request menggunakan http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil
	}
	defer resp.Body.Close()

	// Menampilkan response dari server
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte

	var result map[string]interface{}

	// Decode JSON menjadi map
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	// Menampilkan isi map
	fmt.Println("Code:", result["code"])

	// Mengakses nested data
	datajson, ok := result["data"].(map[string]interface{})
	if !ok {
		fmt.Println("Error: data tidak dapat dikonversi menjadi map[string]interface{}")
		return nil
	}
	return datajson
	// fmt.Println(string(body))              // convert to string before print
}

func UpdateHmi(data Data_hmi) map[string]interface{} {
	// Konversi struct menjadi JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil
	}

	// Membuat request POST dengan JSON
	url := IP + "/hmi/generate-sn"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	// Menambahkan header untuk tipe konten JSON
	req.Header.Set("Content-Type", "application/json")

	// Mengirim request menggunakan http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil
	}
	defer resp.Body.Close()

	// Menampilkan response dari server
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte
	fmt.Println(string(body))              // convert to string before print
	var result map[string]interface{}

	// Decode JSON menjadi map
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	// Menampilkan isi map
	fmt.Println("Code:", result["code"])

	// Mengakses nested data
	datajson, ok := result["data"].(map[string]interface{})
	if !ok {
		fmt.Println("Error: data tidak dapat dikonversi menjadi map[string]interface{}")
		return nil
	}
	return datajson
}

func UpdateKeyless(data Data_keyless) map[string]interface{} {
	// Konversi struct menjadi JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil
	}

	// Membuat request POST dengan JSON
	url := IP + "/keyless/generate-sn"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	// Menambahkan header untuk tipe konten JSON
	req.Header.Set("Content-Type", "application/json")

	// Mengirim request menggunakan http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil
	}
	defer resp.Body.Close()

	// Menampilkan response dari server
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte
	fmt.Println(string(body))              // convert to string before print
	var result map[string]interface{}

	// Decode JSON menjadi map
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	// Menampilkan isi map
	fmt.Println("Code:", result["code"])

	// Mengakses nested data
	datajson, ok := result["data"].(map[string]interface{})
	if !ok {
		fmt.Println("Error: data tidak dapat dikonversi menjadi map[string]interface{}")
		return nil
	}
	return datajson
}

func Updatestatus(sn string, bord string, status string, qr_code string, data Bus) error {

	// Konversi struct menjadi JSON
	jsonData, err := json.Marshal(map[string]string{"status": status, "qr_code": qr_code, "bootloader_version": data.Versifirmbootstr, "firmware_version": data.Versifirmappstr, "ble_bootloader_version": data.Versifirmbootstr2, "ble_firmware_version": data.Versifirmappstr2})
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err
	}

	// Membuat request POST dengan JSON
	url := IP + "/" + bord + "/" + sn
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}
	fmt.Println(url)
	// Menambahkan header untuk tipe konten JSON
	req.Header.Set("Content-Type", "application/json")

	// Mengirim request menggunakan http.Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return err
	}
	defer resp.Body.Close()
	// Cek status respons

	fmt.Println("sudah terkirim")
	// Menampilkan response dari server
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte
	fmt.Println(string(body))
	if resp.StatusCode != http.StatusOK {
		fmt.Println("request failed with status: %s", resp.Status)
		return fmt.Errorf("request failed with status: %s", resp.Status)
	}
	return err
	// convert to string before print
}
