package modul

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func SnBms(data Bus) (uint32, uint32, string) {
	data1 := uint32(0)
	data2 := uint32(0)
	tegangan := [8]int{12, 24, 36, 48, 60, 72, 84, 96}
	// var newCount uint32
	num1, err := strconv.Atoi(data.Num.Tegangan)
	if err != nil {
		fmt.Println("Error:", err)
	}
	data1 = (uint32(num1) & 0xff) << 4
	num2, err := strconv.Atoi(data.Num.ParalelN)
	if err != nil {
		fmt.Println("Error:", err)
	}
	data1 |= (uint32(num2) & 0x0f) | (uint32(data.Num.Jenis[0])&0xff)<<8 | (uint32(data.Num.Type_[0])&0xff)<<16
	num3, err := strconv.Atoi(data.Num.Tahun_ex)
	if err != nil {
		fmt.Println("Error:", err)
	}
	th := num3 - 2000
	data1 |= (uint32(th) & 0xff) << 24

	num4, err := strconv.Atoi(data.Num.Bulan_ex)
	if err != nil {
		fmt.Println("Error:", err)
	}
	data2 |= (uint32(num4) & 0xff) << 12

	battrytipe := ""
	if data.Id == 2 {
		battrytipe = "hight"
	} else if data.Id == 1 {
		battrytipe = "low"
	}
	bms := Data_bms{
		Voltage:       tegangan[num1],
		ParallelNum:   num2,
		CellBrand:     data.Num.Jenis,
		CellBrandType: data.Num.Type_,
		CellProdYear:  num3,
		CellProdMonth: num4,
		ModelVersion:  data.Modelver,
		BatteryType:   battrytipe,
	}
	respData := UpdateBms(bms)

	if respData == nil {
		return data1, data2, ""
	}

	snStr := respData["sn"].(string)
	detail := respData["detail"].(map[string]interface{})
	tmp, _ := detail["pack_prod_month"].(float64)

	data2 |= (uint32(tmp) & 0xff) << 8
	tmp, _ = detail["pack_prod_year"].(float64)
	th = int(tmp) - 2000
	data2 |= (uint32(th) & 0xff)
	// respDetail, isOk := respData["detail"].(map[string]interface{})
	// if !isOk {
	// 	// TODO handle
	// }

	counter, ok := detail["counter"].(float64)
	if !ok {
		fmt.Println("Error: detail['counter'] tidak dapat dikonversi menjadi uint32")
		return data1, data2, ""
	}

	data2 |= (uint32(counter) & 0xffff) << 16

	buatexsel(data.Bord)
	// fmt.Println("Counter:", counter)
	// fmt.Printf("Converted number:%s", snStr)
	fmt.Printf("Converted number:%08x%08x", data1, data2)
	return data1, data2, snStr
}

func SnBmsString(sn1 uint32, sn2 uint32) string {
	// var tegangan string
	sn1Bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sn1Bytes, sn1)

	sn2Bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sn2Bytes, sn2)
	fmt.Println("SN:", sn1, sn2)

	// Byte 0
	voltageMap := map[uint8]string{
		0: "12", 1: "24", 2: "36", 3: "48", 4: "60", 5: "72", 6: "84", 7: "96",
	}
	voltageCode := (sn1Bytes[0] >> 4) & 0x0F
	parallel := sn1Bytes[0] & 0x0F

	voltage, ok := voltageMap[voltageCode]
	if !ok {
		voltage = "??"
	}

	// Byte 1 - Cell Brand
	cellBrand := string(sn1Bytes[1])

	// Byte 2 - Battery Type
	battType := string(sn1Bytes[2])

	// Byte 3 - Year Cell
	yearCell := int(sn1Bytes[3])

	// Byte 4 - Year Pack
	yearPack := int(sn2Bytes[0])

	// Byte 5 - Month Pack & Cell
	monthPack := int(sn2Bytes[1] & 0x0F)
	monthCell := int((sn2Bytes[1] >> 4) & 0x0F)

	// Byte 6-7 - Counter
	counter := binary.BigEndian.Uint16(sn2Bytes[2:4])

	// Final format: XX YY Z T AA CC BB DD NNNN
	result := fmt.Sprintf("%s%02d%s%s%02d%02d%02d%02d%04d",
		voltage, parallel, cellBrand, battType, yearCell, monthCell, yearPack, monthPack, counter)

	// switch data.Num.Tegangan {
	// case "0":
	// 	tegangan = "12"
	// case "1":
	// 	tegangan = "24"
	// case "2":
	// 	tegangan = "36"
	// case "3":
	// 	tegangan = "48"
	// case "4":
	// 	tegangan = "60"
	// case "5":
	// 	tegangan = "72"
	// case "6":
	// 	tegangan = "84"
	// case "7":
	// 	tegangan = "96"
	// }
	// message := fmt.Sprintf("%s%s%s%s%s%s%s%s%04d", tegangan, data.Num.ParalelN, data.Num.Jenis, data.Num.Type_, data.Num.Tahun_ex, data.Num.Bulan_ex, data.Num.Tahun_pb, data.Num.Bulan_pb, coun)

	fmt.Printf("Serial number:%s", result)
	return result
}

func SnHmiString(sn1 uint32, sn2 uint32) (string, uint32) {
	sn1Bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sn1Bytes, sn1)

	// Konversi sn2 ke array 4 byte
	sn2Bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sn2Bytes, sn2)
	sn2f := uint32(sn2Bytes[2])<<24 | uint32(sn2Bytes[3])<<16 | uint32(sn2Bytes[1])<<8 | uint32(sn2Bytes[0]) // kerena di saat flash bayte 2 dan 3 dibalik lagi .

	caun := uint16(sn2Bytes[2])<<8 | uint16(sn2Bytes[3])
	fmt.Printf("Nilai caun: %d %02x %02x \n", caun, sn2Bytes[2], sn2Bytes[3])

	message := fmt.Sprintf("%c%c%c%02d%02d%02d%04d", sn1Bytes[0], sn1Bytes[1], sn1Bytes[2], sn1Bytes[3], sn2Bytes[0], sn2Bytes[1], caun)
	// switch data.Num.Tegangan {
	// case "0":
	// 	tegangan = "12"
	// case "1":
	// 	tegangan = "24"
	// case "2":
	// 	tegangan = "36"
	// case "3":
	// 	tegangan = "48"
	// case "4":
	// 	tegangan = "60"
	// case "5":
	// 	tegangan = "72"
	// case "6":
	// 	tegangan = "84"
	// case "7":
	// 	tegangan = "96"
	// }
	// message := fmt.Sprintf("%s%s%s%s%s%s%s%s%04d", tegangan, data.Num.ParalelN, data.Num.Jenis, data.Num.Type_, data.Num.Tahun_ex, data.Num.Bulan_ex, data.Num.Tahun_pb, data.Num.Bulan_pb, coun)

	fmt.Printf("Serial number:%s", message)
	return message, sn2f
}

func SnVCUString(sn1 uint32, sn2 uint32) string {
	// Konversi sn1 ke array 4 byte
	sn1Bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sn1Bytes, sn1)

	// Konversi sn2 ke array 4 byte
	sn2Bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sn2Bytes, sn2)

	caun := uint16(sn2Bytes[2])<<8 | uint16(sn2Bytes[3])
	fmt.Printf("Nilai caun: %d %02x %02x \n", caun, sn2Bytes[2], sn2Bytes[3])

	message := fmt.Sprintf("%c%c%02d%02d%02d%04d", sn1Bytes[0], sn1Bytes[1], sn1Bytes[2], sn2Bytes[0], sn2Bytes[1], caun)

	// // Contoh penggunaan bytes untuk membuat string
	// // Misalnya kita ingin mengambil byte pertama dari sn1 sebagai tegangan
	// tegangan := fmt.Sprintf("%c", sn1Bytes[0]>>4) // Mengambil 4 bit pertama

	// // Mengambil 4 bit terakhir dari byte pertama sn1 untuk paralel
	// paralel := fmt.Sprintf("%d", sn1Bytes[0]&0x0F)

	// // Mengambil byte kedua sn1 untuk jenis
	// jenis := string([]byte{sn1Bytes[1]})

	// // Mengambil byte ketiga sn1 untuk tipe
	// tipe := string([]byte{sn1Bytes[2]})

	// // Mengambil byte keempat sn1 untuk tahun
	// tahun := fmt.Sprintf("%d", 2000+int(sn1Bytes[3]))

	// // Mengambil byte pertama sn2 untuk bulan
	// bulan := fmt.Sprintf("%d", sn2Bytes[0])

	// // Mengambil byte kedua sn2 untuk tahun produksi
	// tahunPb := fmt.Sprintf("%d", 2000+int(sn2Bytes[1]))

	// // Mengambil byte ketiga sn2 untuk bulan produksi
	// bulanPb := fmt.Sprintf("%d", sn2Bytes[2])

	// // Mengambil 2 byte terakhir sn2 untuk counter
	// counter := binary.LittleEndian.Uint16(sn2Bytes[2:4])

	// message := fmt.Sprintf("%s%s%s%s%s%s%s%s%04d",
	//     tegangan, paralel, jenis, tipe, tahun, bulan, tahunPb, bulanPb, counter)

	fmt.Printf("Serial number: %s\n", message)
	return message
	// return ""
}

func SnKeylessString(sn1 uint32, sn2 uint32) (string, uint32) {
	// Konversi sn1 ke array 4 byte
	sn1Bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sn1Bytes, sn1)

	// Konversi sn2 ke array 4 byte
	sn2Bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sn2Bytes, sn2)
	sn2f := uint32(sn2Bytes[2])<<24 | uint32(sn2Bytes[3])<<16 | uint32(sn2Bytes[1])<<8 | uint32(sn2Bytes[0]) // kerena di saat flash bayte 2 dan 3 dibalik lagi .

	caun := uint16(sn2Bytes[2])<<8 | uint16(sn2Bytes[3])
	fmt.Printf("Nilai caun: %d %02x %02x \n", caun, sn2Bytes[3], sn2Bytes[2])

	message := fmt.Sprintf("%c%c%02d%02d%02d%04d", sn1Bytes[0], sn1Bytes[1], sn1Bytes[2], sn2Bytes[0], sn2Bytes[1], caun)

	// // Contoh penggunaan bytes untuk membuat string
	// // Misalnya kita ingin mengambil byte pertama dari sn1 sebagai tegangan
	// tegangan := fmt.Sprintf("%c", sn1Bytes[0]>>4) // Mengambil 4 bit pertama

	// // Mengambil 4 bit terakhir dari byte pertama sn1 untuk paralel
	// paralel := fmt.Sprintf("%d", sn1Bytes[0]&0x0F)

	// // Mengambil byte kedua sn1 untuk jenis
	// jenis := string([]byte{sn1Bytes[1]})

	// // Mengambil byte ketiga sn1 untuk tipe
	// tipe := string([]byte{sn1Bytes[2]})

	// // Mengambil byte keempat sn1 untuk tahun
	// tahun := fmt.Sprintf("%d", 2000+int(sn1Bytes[3]))

	// // Mengambil byte pertama sn2 untuk bulan
	// bulan := fmt.Sprintf("%d", sn2Bytes[0])

	// // Mengambil byte kedua sn2 untuk tahun produksi
	// tahunPb := fmt.Sprintf("%d", 2000+int(sn2Bytes[1]))

	// // Mengambil byte ketiga sn2 untuk bulan produksi
	// bulanPb := fmt.Sprintf("%d", sn2Bytes[2])

	// // Mengambil 2 byte terakhir sn2 untuk counter
	// counter := binary.LittleEndian.Uint16(sn2Bytes[2:4])

	// message := fmt.Sprintf("%s%s%s%s%s%s%s%s%04d",
	//     tegangan, paralel, jenis, tipe, tahun, bulan, tahunPb, bulanPb, counter)

	fmt.Printf("Serial number: %s\n", message)
	return message, sn2f
	// return ""
}

func SNhmi(data Bus) (uint32, uint32, string, error) {
	data1 := uint32(0)
	data2 := uint32(0)
	num3, err := strconv.Atoi(data.Num.Type_)
	if err != nil {
		fmt.Println("Error:", err)
		return data1, data2, "", errors.New("Error: " + err.Error())
	}
	data1 |= (uint32(num3)&0xff)<<24 | (uint32(data.Num.Jenis[0])&0xff)<<16 | (uint32(data.Num.MCU[1])&0xff)<<8 | (uint32(data.Num.MCU[0]) & 0xff)

	hmi := Data_hmi{
		Mcu:           data.Num.MCU,
		LcdTypeId:     data.Num.Jenis,
		VehicleTypeId: num3,
		ModelVersion:  data.Modelver,
	}
	respData := UpdateHmi(hmi)

	if respData == nil {
		return data1, data2, "", errors.New("gagal mendapatkan data dari server")
	}

	snStr := respData["sn"].(string)
	detail := respData["detail"].(map[string]interface{})
	tmp, _ := detail["month"].(float64)

	data2 |= (uint32(tmp) & 0xff) << 8
	tmp, _ = detail["year"].(float64)
	th := int(tmp) - 2000
	data2 |= (uint32(th) & 0xff)

	counter, ok := detail["counter"].(float64)
	if !ok {
		fmt.Println("Error: detail['counter'] tidak dapat dikonversi menjadi uint32")
		return data1, data2, "", errors.New("Error: " + err.Error())
	}

	data2 |= (uint32(counter) & 0xffff) << 16

	buatexsel(data.Bord)
	fmt.Println("Counter:", counter)
	// fmt.Printf("Converted number:%c %c ", data.Num.MCU[0], data.Num.MCU[1])
	fmt.Printf("Converted number:%08x%08x", data1, data2)
	return data1, data2, snStr, nil
}

func SNVCU(data Bus) (uint32, uint32, string, uint32, error) {
	data1 := uint32(0)
	data2 := uint32(0)
	num3, err := strconv.Atoi(data.Num.Type_)
	if err != nil {
		fmt.Println("Error:", err)
	}
	data1 |= (uint32(num3)&0xff)<<16 | (uint32(data.Num.MCU[1])&0xff)<<8 | (uint32(data.Num.MCU[0]) & 0xff)

	vcu := Data_vcu{
		Mcu:           data.Num.MCU,
		VehicleTypeID: num3,
		ModelVersion:  data.Modelver,
	}
	respData := UpdateVcu(vcu)

	if respData == nil {
		// TODO handle
		fmt.Println("Error: UpdateVcu tidak dapat dikonversi menjadi uint32")
		return data1, data2, "", 0, errors.New("UpdateVcu tidak dapat dikonversi menjadi uint32")
	}

	snStr := respData["sn"].(string)
	detail := respData["detail"].(map[string]interface{})
	tmp, _ := detail["month"].(float64)

	data2 |= (uint32(tmp) & 0xff) << 8
	tmp, _ = detail["year"].(float64)
	th := int(tmp) - 2000
	data2 |= (uint32(th) & 0xff)

	tmp, _ = detail["vin"].(float64)

	counter, ok := detail["counter"].(float64)
	if !ok {
		fmt.Println("Error: detail['counter'] tidak dapat dikonversi menjadi uint32")
		return data1, data2, "", 0, errors.New("detail['counter'] tidak dapat dikonversi menjadi uint32")
	}

	data2 |= (uint32(counter) & 0xffff) << 16

	buatexsel(data.Bord)
	fmt.Println("Counter:", counter)
	// fmt.Printf("Converted number:%c %c ", data.Num.MCU[0], data.Num.MCU[1])
	fmt.Printf("Converted number %d :%08x%08x", uint32(tmp), data1, data2)
	return data1, data2, snStr, uint32(tmp), nil
}

func SNKeyless(data Bus) (uint32, uint32, string, error) {
	data1 := uint32(0)
	data2 := uint32(0)
	num3, err := strconv.Atoi(data.Num.Type_)
	if err != nil {
		fmt.Println("Error:", err)
	}
	data1 |= (uint32(num3)&0xff)<<16 | (uint32(data.Num.MCU[1])&0xff)<<8 | (uint32(data.Num.MCU[0]) & 0xff)

	keyless := Data_keyless{
		Mcu:           data.Num.MCU,
		VehicleTypeId: num3,
		ModelVersion:  data.Modelver,
	}
	respData := UpdateKeyless(keyless)

	if respData == nil {
		// TODO handle
		fmt.Println("Error: UpdateKeyless tidak dapat dikonversi menjadi uint32")
		return data1, data2, "", errors.New("UpdateKeyless tidak dapat dikonversi menjadi uint32")
	}

	snStr := respData["sn"].(string)
	detail := respData["detail"].(map[string]interface{})
	tmp, _ := detail["month"].(float64)

	data2 |= (uint32(tmp) & 0xff) << 8
	tmp, _ = detail["year"].(float64)
	th := int(tmp) - 2000
	data2 |= (uint32(th) & 0xff)

	counter, ok := detail["counter"].(float64)
	if !ok {
		fmt.Println("Error: detail['counter'] tidak dapat dikonversi menjadi uint32")
		return data1, data2, "", errors.New("detail['counter'] tidak dapat dikonversi menjadi uint32")
	}

	data2 |= (uint32(counter) & 0xffff) << 16

	buatexsel(data.Bord)
	fmt.Println("Counter:", counter)
	// fmt.Printf("Converted number:%c %c ", data.Num.MCU[0], data.Num.MCU[1])
	fmt.Printf("Converted number:%08x%08x", data1, data2)
	return data1, data2, snStr, nil
}

func SNconvert(sn string, bord string) (uint32, uint32, error) {
	parts := strings.Fields(sn)
	if len(parts) == 0 {
		fmt.Println("SN kosong atau format salah")
		return 0, 0, errors.New("SN kosong atau format salah")
	}
	fmt.Println("SN:", len(parts))
	switch strings.ToLower(bord) {
	case "hmi":
		// Format: ST S 01 2403 0001
		// → XX Y ZZ BBBB NNNN
		sn1, sn2, err := ParseHmiString(sn)
		if err != nil {
			fmt.Println("Error:", err)
			return 0, 0, err
		}
		return sn1, sn2, nil

	case "vcu":
		// Format: ST 01 00 2403 0001
		// → XX ZZ CC BBBB NNNN
		sn1, sn2, err := ParseSnVCU(sn)
		if err != nil {
			fmt.Println("Error:", err)
			return 0, 0, err
		}
		return sn1, sn2, nil

	case "keyfob":
		// Format: RS 01 2403 0001
		// → XX ZZ BBBB NNNN

		sn1, sn2, err := ParseKeylessString(sn)
		if err != nil {
			fmt.Println("Error:", err)
			return 0, 0, err
		}
		return sn1, sn2, nil

	case "bms":
		// Format: 72 07 E A 2301 2403 0001
		// → XX YY Z T AAAA BBBB NNNN
		sn1, sn2, err := ParseBmsString(sn)
		if err != nil {
			fmt.Println("Error:", err)
			return 0, 0, err
		}
		return sn1, sn2, nil

	default:
		return 0, 0, errors.New("board tidak dikenal")
	}

	return 0, 0, nil
}

func ParseKeylessString(message string) (uint32, uint32, error) {
	if len(message) < 12 {
		return 0, 0, fmt.Errorf("serial number terlalu pendek")
	}

	// Ambil bagian string
	ch1 := message[0]
	ch2 := message[1]
	sn1b2, err1 := strconv.Atoi(message[2:4])
	sn2b0, err2 := strconv.Atoi(message[4:6])
	sn2b1, err3 := strconv.Atoi(message[6:8])
	caun, err4 := strconv.Atoi(message[8:12])

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return 0, 0, fmt.Errorf("gagal parsing angka dari serial number")
	}

	// Susun ulang sn1Bytes
	sn1Bytes := []byte{
		ch1,
		ch2,
		byte(sn1b2),
		0, // byte ke-4 tidak dipakai di serial number (default 0)
	}

	// Susun ulang sn2Bytes
	sn2Bytes := []byte{
		byte(sn2b0),
		byte(sn2b1),
		byte(caun & 0xFF),
		byte((caun >> 8) & 0xFF),
	}

	// Konversi ke uint32 (Little Endian)
	sn1 := binary.LittleEndian.Uint32(sn1Bytes)
	sn2 := binary.LittleEndian.Uint32(sn2Bytes)

	return sn1, sn2, nil
}

func ParseBmsString(serial string) (uint32, uint32, error) {
	if len(serial) < 14 {
		return 0, 0, fmt.Errorf("serial number terlalu pendek")
	}

	// 1. Ambil bagian string
	voltageStr := serial[0:2]
	parallelStr := serial[2:4]
	cellBrand := serial[4:5]
	battType := serial[5:6]
	yearCellStr := serial[6:8]
	monthCellStr := serial[8:10]
	yearPackStr := serial[10:12]
	monthPackStr := serial[12:14]
	counterStr := serial[14:]

	// 2. Konversi angka
	parallel, _ := strconv.Atoi(parallelStr)
	yearCell, _ := strconv.Atoi(yearCellStr)
	monthCell, _ := strconv.Atoi(monthCellStr)
	yearPack, _ := strconv.Atoi(yearPackStr)
	monthPack, _ := strconv.Atoi(monthPackStr)
	counter, _ := strconv.Atoi(counterStr)

	// 3. Cari voltage code
	voltageMap := map[string]uint8{
		"12": 0, "24": 1, "36": 2, "48": 3, "60": 4,
		"72": 5, "84": 6, "96": 7,
	}
	voltageCode, ok := voltageMap[voltageStr]
	if !ok {
		return 0, 0, fmt.Errorf("voltage %s tidak dikenal", voltageStr)
	}

	// 4. Susun sn1Bytes
	sn1Bytes := []byte{
		(byte(voltageCode)<<4 | byte(parallel&0x0F)),
		cellBrand[0],
		battType[0],
		byte(yearCell),
	}

	// 5. Susun sn2Bytes
	sn2Bytes := []byte{
		byte(yearPack),
		byte((monthCell << 4) | (monthPack & 0x0F)),
		byte((counter >> 8) & 0xFF),
		byte(counter & 0xFF),
	}

	// 6. Konversi ke uint32
	sn1 := binary.LittleEndian.Uint32(sn1Bytes)
	sn2 := binary.LittleEndian.Uint32(sn2Bytes)

	return sn1, sn2, nil
}

func ParseHmiString(serial string) (uint32, uint32, error) {
	if len(serial) < 13 {
		return 0, 0, fmt.Errorf("serial number terlalu pendek")
	}

	// 1. Ambil bagian string
	ch1 := serial[0]
	ch2 := serial[1]
	ch3 := serial[2]
	sn1b3Str := serial[3:5]
	sn2b0Str := serial[5:7]
	sn2b1Str := serial[7:9]
	caunStr := serial[9:13]

	// 2. Konversi angka
	sn1b3, err1 := strconv.Atoi(sn1b3Str)
	sn2b0, err2 := strconv.Atoi(sn2b0Str)
	sn2b1, err3 := strconv.Atoi(sn2b1Str)
	caun, err4 := strconv.Atoi(caunStr)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return 0, 0, fmt.Errorf("gagal parsing angka dari serial number")
	}

	// 3. Susun sn1Bytes
	sn1Bytes := []byte{
		ch1,
		ch2,
		ch3,
		byte(sn1b3),
	}

	// 4. Susun sn2Bytes
	sn2Bytes := []byte{
		byte(sn2b0),
		byte(sn2b1),
		byte(caun & 0xFF),
		byte((caun >> 8) & 0xFF),
	}

	// 5. Konversi ke uint32
	sn1 := binary.LittleEndian.Uint32(sn1Bytes)
	sn2 := binary.LittleEndian.Uint32(sn2Bytes)

	return sn1, sn2, nil
}

func ParseSnVCU(serial string) (uint32, uint32, error) {
	if len(serial) < 11 {
		return 0, 0, fmt.Errorf("serial terlalu pendek: %s", serial)
	}

	// Pecah field sesuai format:
	// AB0307090258
	char1 := serial[0]                   // sn1Bytes[0]
	char2 := serial[1]                   // sn1Bytes[1]
	val1, _ := strconv.Atoi(serial[2:4]) // sn1Bytes[2]
	val2, _ := strconv.Atoi(serial[4:6]) // sn2Bytes[0]
	val3, _ := strconv.Atoi(serial[6:8]) // sn2Bytes[1]
	caun, _ := strconv.Atoi(serial[8:])  // sn2Bytes[2:3] (big endian)

	// Bangun sn1Bytes
	sn1Bytes := []byte{char1, char2, byte(val1), 0}
	sn1 := binary.LittleEndian.Uint32(sn1Bytes)

	// Bangun sn2Bytes
	sn2Bytes := make([]byte, 4)
	sn2Bytes[0] = byte(val2)
	sn2Bytes[1] = byte(val3)
	sn2Bytes[2] = byte(caun >> 8)
	sn2Bytes[3] = byte(caun & 0xFF)
	sn2 := binary.LittleEndian.Uint32(sn2Bytes)

	return sn1, sn2, nil
}
