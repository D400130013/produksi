package modul

import (
	"encoding/binary"
	"fmt"
	"strconv"
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

	if respData != nil {
		// TODO handle
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

func SnBmsString(data Bus, coun uint32) string {
	var tegangan string

	switch data.Num.Tegangan {
	case "0":
		tegangan = "12"
	case "1":
		tegangan = "24"
	case "2":
		tegangan = "36"
	case "3":
		tegangan = "48"
	case "4":
		tegangan = "60"
	case "5":
		tegangan = "72"
	case "6":
		tegangan = "84"
	case "7":
		tegangan = "96"
	}
	message := fmt.Sprintf("%s%s%s%s%s%s%s%s%04d", tegangan, data.Num.ParalelN, data.Num.Jenis, data.Num.Type_, data.Num.Tahun_ex, data.Num.Bulan_ex, data.Num.Tahun_pb, data.Num.Bulan_pb, coun)

	fmt.Printf("Serial number:%s", message)
	return message
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

func SNhmi(data Bus) (uint32, uint32, string) {
	data1 := uint32(0)
	data2 := uint32(0)
	num3, err := strconv.Atoi(data.Num.Type_)
	if err != nil {
		fmt.Println("Error:", err)
	}
	data1 |= (uint32(num3)&0xff)<<24 | (uint32(data.Num.Jenis[0])&0xff)<<16 | (uint32(data.Num.MCU[1])&0xff)<<8 | (uint32(data.Num.MCU[0]) & 0xff)

	hmi := Data_hmi{
		Mcu:           data.Num.MCU,
		LcdTypeId:     data.Num.Jenis,
		VehicleTypeId: num3,
		ModelVersion:  data.Modelver,
	}
	respData := UpdateHmi(hmi)

	if respData != nil {
		// TODO handle
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
		return data1, data2, ""
	}

	data2 |= (uint32(counter) & 0xffff) << 16

	buatexsel(data.Bord)
	fmt.Println("Counter:", counter)
	// fmt.Printf("Converted number:%c %c ", data.Num.MCU[0], data.Num.MCU[1])
	fmt.Printf("Converted number:%08x%08x", data1, data2)
	return data1, data2, snStr
}

func SNVCU(data Bus) (uint32, uint32, string, uint32) {
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

	if respData != nil {
		// TODO handle
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
		return data1, data2, "", 0
	}

	data2 |= (uint32(counter) & 0xffff) << 16

	buatexsel(data.Bord)
	fmt.Println("Counter:", counter)
	// fmt.Printf("Converted number:%c %c ", data.Num.MCU[0], data.Num.MCU[1])
	fmt.Printf("Converted number %d :%08x%08x", uint32(tmp), data1, data2)
	return data1, data2, snStr, uint32(tmp)
}

func SNKeyless(data Bus) (uint32, uint32, string) {
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

	if respData != nil {
		// TODO handle
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
		return data1, data2, ""
	}

	data2 |= (uint32(counter) & 0xffff) << 16

	buatexsel(data.Bord)
	fmt.Println("Counter:", counter)
	// fmt.Printf("Converted number:%c %c ", data.Num.MCU[0], data.Num.MCU[1])
	fmt.Printf("Converted number:%08x%08x", data1, data2)
	return data1, data2, snStr
}
