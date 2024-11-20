package modul

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func flashFirmware(firmwarePath string) error {
	// Path ke STM32_Programmer_CLI.exe
	programmerCLI := "C:\\Program Files\\STMicroelectronics\\STM32Cube\\STM32CubeProgrammer\\bin\\STM32_Programmer_CLI.exe"

	// Perintah untuk mem-flash firmware menggunakan STM32CubeProgrammer
	cmd := exec.Command(programmerCLI, "-c", "port=SWD", "-w", firmwarePath, "0x08000000", "--verify")

	// Jalankan perintah dan ambil output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("flashing firmware failed: %v. Output: %s", err, string(output))
	}

	fmt.Println("Firmware flashing successful!")
	return nil
}

func ExecuteOpenOCDvcu(data Bus) error {
	filePath := fmt.Sprintf("./bin/vcu/%d/bootloader.bin", data.Model)
	filePath1 := fmt.Sprintf("./bin/vcu/%d/application.bin", data.Model)
	bintmp := "output.bin"
	datasn1, datasn2, sn, id_ := SNVCU(data)
	data.Id = id_
	dummy := uint64(0xffffffff)
	// Periksa apakah file ada
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("File gak ada")
		Updatestatus(sn, data.Bord, "ERROR")
		Dialogeror("Gagal Flash VCU" + filePath)
		return fmt.Errorf("file tidak ditemukan: %s", filePath)
	}
	datasn3 := rearrangeBytes(datasn2)
	// Menyiapkan buffer data dalam format byte
	value := make([]byte, 32) // Sesuaikan ukuran buffer

	// Menulis nilai ke dalam buffer dengan urutan
	binary.LittleEndian.PutUint32(value[0:], datasn1) // 8 byte untuk SN
	binary.LittleEndian.PutUint32(value[4:], datasn3)
	binary.LittleEndian.PutUint32(value[8:], uint32(dummy))       // 4 byte untuk Dummy
	binary.LittleEndian.PutUint32(value[12:], uint32(dummy))      // 4 byte untuk Dummy
	binary.LittleEndian.PutUint32(value[16:], uint32(dummy))      // 4 byte untuk Dummy
	binary.LittleEndian.PutUint32(value[20:], data.Versifirmboot) // 4 byte untuk VIN
	binary.LittleEndian.PutUint32(value[24:], data.Model)
	binary.LittleEndian.PutUint32(value[28:], data.Id)
	CreateBinFile(bintmp, value)

	// Perintah untuk menjalankan OpenOCD
	cmd := exec.Command(".\\openocd\\bin\\openocd.exe",
		"-f", ".\\openocd\\share\\openocd\\scripts\\interface\\stlink-v2.cfg",
		"-f", ".\\openocd\\share\\openocd\\scripts\\target\\stm32h7x_dual_bank.cfg",
		"-c", fmt.Sprintf("init; halt; flash erase_sector 0 0 7; flash erase_sector 1 0 5; program %s 0x08000000; halt; flash write_image ./output.bin 0x0803ffe0; program %s 0x08040000; halt; flash filld 0x080ffff8 0x00000000%08x 1;  reset; exit", filePath, filePath1, data.Versifirmapp))
	// fmt.Sprintf("init; halt; flash erase_sector 0 0 7; flash erase_sector 1 0 5; program %s 0x08000000; flash filld 0x0803fff8 0x%08x%08x 1;halt; program %s 0x08040000; halt;flash filld 0x0803FFD0 0x%08x%08x 1 ; halt; flash fillw 0x0803FFF4 0x%08x 1 ; halt; flash filld 0x080ffff8 0x00000000%08x 1;  reset; exit", filePath, data.Id, data.Model, filePath1, datasn3, datasn1, data.Versifirmboot, data.Versifirmapp))
	// Jalankan perintah dan ambil output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(cmd.String())
		fmt.Println("Output:", string(output))

		Updatestatus(sn, data.Bord, "ERROR")
		Dialogeror("Gagal Flash VCU" + filePath)
		err = os.Remove(bintmp)
		if err != nil {
			return fmt.Errorf("gagal menghapus file: %w", err)
		}
		return fmt.Errorf("eksekusi OpenOCD gagal: %v. Output: %s", err, string(output))
	}

	fmt.Println("OpenOCD eksekusi berhasil!")
	Dialoginfo("flash sukses" + sn)
	id := fmt.Sprintf("%d", data.Id)
	tambahbaris(data.Bord, sn, id)
	Updatestatus(sn, data.Bord, "SUCCESS")
	// fmt.Println("Output:", string(output))
	err = os.Remove(bintmp)
	if err != nil {
		return fmt.Errorf("gagal menghapus file: %w", err)
	}
	return nil
}

func ProgramRSL() error {
	// Perintah st-flash untuk membaca dari alamat 0x08000000 dengan panjang 0x1000

	// Ganti tempat-tempat ini dengan nilai yang sesuai
	jlinkScriptFile := "ble_vcu.jlink"
	CreateJLinkScript(jlinkScriptFile)
	// Mengasumsikan J-Link ada dalam PATH sistem
	jlinkPath, err := exec.LookPath(".\\JLink\\JLink.exe")
	if err != nil {
		// fmt.Println("JLinkExe tidak ditemukan dalam PATH. Pastikan perangkat lunak J-Link terpasang dan ada dalam PATH.")
		return err
	}
	fmt.Println("PROCESS BLE VCU...")
	device := "RSL10"
	iface := "SWD"
	speed := "4000"

	// Buka skrip J-Link
	cmd := exec.Command(jlinkPath,
		"-device", device,
		"-if", iface,
		"-speed", speed,
		"-CommanderScript", jlinkScriptFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		// fmt.Println("Error menjalankan J-Link:", err)
		Dialogeror("Gagal Flash RSL")
		return err
	}

	Dialoginfo("flash sukses RSL")
	// fmt.Println("Proses flashing berhasil.")
	return nil

}

func CreateJLinkScript(fileName string) error {
	model := fmt.Sprintf("%d", Data.Model)
	jlinkScript := "Erase\nloadbin \"./bin/vcu/" + model + "/ble-bootloader.bin\", 0x00100000\n" +
		"r\n" + "loadbin \"./bin/vcu/" + model + "/ble-application.bin\", 0x00107000\n" +
		"r\n" + "q\n"

	err := ioutil.WriteFile(fileName, []byte(jlinkScript), 0644)
	if err != nil {
		return err
	}
	fmt.Printf("file .jlink '%s' berhasil dibuat.\n", fileName)
	return nil
}

func ExecuteOpenOCDhmi(data Bus) error {
	// filePath := "./bin/hmi/bootloader.bin"
	// filePath1 := "./bin/hmi/application.bin"
	filePath := fmt.Sprintf("./bin/hmi/%d/bootloader.bin", data.Model)
	filePath1 := fmt.Sprintf("./bin/hmi/%d/application.bin", data.Model)
	datasn1, datasn2, sn := SNhmi(data)
	// Periksa apakah file ada
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("File gak ada")
		Updatestatus(sn, data.Bord, "ERROR")
		Dialogeror("Gagal Flash" + filePath)
		return fmt.Errorf("file tidak ditemukan: %s", filePath)
	}

	datasn3 := rearrangeBytes(datasn2)

	// Perintah untuk menjalankan OpenOCD

	cmd := exec.Command(".\\openocd\\bin\\openocd.exe",
		"-f", ".\\openocd\\share\\openocd\\scripts\\interface\\stlink-v2.cfg",
		"-f", ".\\openocd\\share\\openocd\\scripts\\target\\stm32h7x_dual_bank.cfg",
		"-c", fmt.Sprintf("init; halt; flash erase_sector 0 0 7;flash erase_sector 1 0 5; program %s 0x08000000; flash filld 0x0803FFD0 0x%08x%08x 1; halt;program %s 0x08040000; flash filld 0x0803FFF8 0x%08x%08x 1; flash filld 0x0811FFF8 0x00000000%08x 1; halt; reset; exit", filePath, datasn3, datasn1, filePath1, data.Model, data.Versifirmboot, data.Versifirmapp))

	// Jalankan perintah dan ambil output
	output, err := cmd.CombinedOutput()
	if err != nil {

		Updatestatus(sn, data.Bord, "ERROR")
		Dialogeror("Gagal Flash" + filePath)
		fmt.Println("Output:", string(output))
		fmt.Println(cmd.String())
		return fmt.Errorf("eksekusi OpenOCD gagal: %v. Output: %s", err, string(output))
	}

	// time.Sleep(1 * time.Second)
	// cmd = exec.Command(".\\openocd\\bin\\openocd.exe",
	// 	"-f", ".\\openocd\\share\\openocd\\scripts\\interface\\stlink-v2.cfg",
	// 	"-f", ".\\openocd\\share\\openocd\\scripts\\target\\stm32h7x_dual_bank.cfg",
	// 	"-c", fmt.Sprintf("init; halt; flash filld 0x0803FFE0 0x%08x%08x 1; halt;reset;exit", datasn3, datasn1))

	// // Jalankan perintah dan ambil output
	// output, err = cmd.CombinedOutput()
	// if err != nil {

	// 	Updatestatus(sn, data.Bord, "ERROR")
	// 	Dialogeror("Gagal Flash" + filePath)
	// 	fmt.Println("Output:", string(output))
	// 	fmt.Println(cmd.String())
	// 	return fmt.Errorf("eksekusi OpenOCD gagal: %v. Output: %s", err, string(output))
	// }

	fmt.Println("OpenOCD eksekusi berhasil!")
	fmt.Println("Output:", string(output))
	Dialoginfo("flash sukses" + sn)
	fmt.Println("OpenOCD eksekusi berhasil!", filePath)
	tambahbaris(data.Bord, sn, data.Modelver)
	Updatestatus(sn, data.Bord, "SUCCESS")
	return nil
}

func ExecuteOpenOCDBMS(data Bus) error {
	filePath := fmt.Sprintf("./bin/bms/%d/bootloader.bin", data.Model)
	filePath1 := fmt.Sprintf("./bin/bms/%d/application.bin", data.Model)
	loop := 0
	datasn1, datasn2, sn := SnBms(data)

	// numcoun := int((datasn2 & 0xffff0000) >> 16)

	// Periksa apakah file ada
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// fmt.Println("File gak ada %s", filePath)
		Dialogeror("File gak ada" + filePath)
		Updatestatus(sn, data.Bord, "ERROR")
		return fmt.Errorf("file tidak ditemukan: %s", filePath)
	}
	// SnBms(data)

	datasn3 := rearrangeBytes(datasn2)
ulang:
	// Perintah untuk menjalankan OpenOCD
	cmd := exec.Command(".\\openocd\\bin\\openocd.exe",
		"-f", ".\\openocd\\share\\openocd\\scripts\\interface\\stlink-v2.cfg",
		"-f", ".\\openocd\\share\\openocd\\scripts\\target\\stm32f1x.cfg",
		"-c", fmt.Sprintf("flash init; init; halt; flash erase_sector 0 0 127; flash write_image erase %s 0x08000000;flash filld 0x08004fe0 0x%08x%08x 1; flash filld 0x08004ff0 0x%08xffffffff 1; flash filld 0x08004ff8 0x%08x%08x 1; flash write_image erase %s 0x08005000; halt;flash filld 0x0801fff8 0x00000000%08x 1; reset; exit", filePath, datasn3, datasn1, data.Versifirmboot, data.Id, data.Model, filePath1, data.Versifirmapp))

	// Jalankan perintah dan ambil output
	output, err := cmd.CombinedOutput()
	if err != nil {
		if loop <= 1 {
			loop++
			goto ulang

		}
		// fmt.Println("Output:", string(output))
		Updatestatus(sn, data.Bord, "ERROR")
		Dialogeror("Gagal Flash" + filePath)
		return fmt.Errorf("eksekusi OpenOCD gagal: %v. Output: %s", err, string(output))
	}
	typepack := ""
	if data.Id == 2 {
		typepack = "HIGH"
	} else {
		typepack = "LOW"
	}

	Dialoginfo("flash sukses" + sn)
	fmt.Println("OpenOCD eksekusi berhasil!", filePath)
	tambahbaris(data.Bord, sn, typepack)
	Updatestatus(sn, data.Bord, "SUCCESS")
	// fmt.Println("Output:", string(output))
	return nil
}

func ProgramKeyles() error {
	// Perintah st-flash untuk membaca dari alamat 0x08000000 dengan panjang 0x1000

	// Ganti tempat-tempat ini dengan nilai yang sesuai
	jlinkScriptFile := "ble_keyless.jlink"
	CreateJLinkScriptkeyles(jlinkScriptFile)
	// Mengasumsikan J-Link ada dalam PATH sistem
	jlinkPath, err := exec.LookPath(".\\JLink\\JLink.exe")
	if err != nil {
		// fmt.Println("JLinkExe tidak ditemukan dalam PATH. Pastikan perangkat lunak J-Link terpasang dan ada dalam PATH.")
		return err
	}
	fmt.Println("PROCESS BLE VCU...")
	device := "RSL10"
	iface := "SWD"
	speed := "4000"

	// Buka skrip J-Link
	cmd := exec.Command(jlinkPath,
		"-device", device,
		"-if", iface,
		"-speed", speed,
		"-CommanderScript", jlinkScriptFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error menjalankan J-Link: %v\n", err)
		Dialogeror("Gagal Flash RSL")
		return err
	}

	Dialoginfo("flash sukses RSL")
	// fmt.Println("Proses flashing berhasil.")
	return nil

}

func CreateJLinkScriptkeyles(fileName string) error {
	model := fmt.Sprintf("%d", Data.Model)
	jlinkScript := "Erase\nloadbin \"./bin/keyless/" + model + "/application.bin\", 0x00100000\n" +
		"r\n" + "q\n"

	err := ioutil.WriteFile(fileName, []byte(jlinkScript), 0644)
	if err != nil {
		return err
	}
	fmt.Printf("file .jlink '%s' berhasil dibuat.\n", fileName)
	return nil
}

func ExecuteOpenOCDkeyless(data Bus) error {
	// filePath := "./bin/hmi/bootloader.bin"
	// filePath1 := "./bin/hmi/application.bin"

	filePath := fmt.Sprintf("./bin/keyless/%d/application.bin", data.Model)
	datasn1, datasn2, sn := SNKeyless(data)
	// Periksa apakah file ada
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("File gak ada")
		Updatestatus(sn, data.Bord, "ERROR")
		Dialogeror("Gagal Flash" + filePath)
		return fmt.Errorf("file tidak ditemukan: %s", filePath)
	}

	datasn3 := rearrangeBytes(datasn2)

	// Perintah untuk menjalankan OpenOCD

	cmd := exec.Command(".\\openocd\\bin\\openocd.exe",
		"-f", ".\\openocd\\share\\openocd\\scripts\\interface\\jlink.cfg",
		"-f", ".\\openocd\\share\\openocd\\scripts\\target\\rsl10.cfg",
		"-c", fmt.Sprintf("adapter speed 1000; flash init; init; halt; flash protect 0 0 2 off; flash erase_sector 0 0 191; program %s 0x00100000; flash filld 0x0015fff0 0x%08x%08x 1;halt;flash filld 0x0015fff8 0x%08x%08x 1;halt; reset; exit;", filePath, datasn3, datasn1, data.Model, data.Versifirmapp))

	// Jalankan perintah dan ambil output
	output, err := cmd.CombinedOutput()
	if err != nil {

		Updatestatus(sn, data.Bord, "ERROR")
		Dialogeror("Gagal Flash" + filePath)
		fmt.Println("Output:", string(output))
		fmt.Println(cmd.String())
		return fmt.Errorf("eksekusi OpenOCD gagal: %v. Output: %s", err, string(output))
	}

	fmt.Println("OpenOCD eksekusi berhasil!")
	fmt.Println("Output:", string(output))
	Dialoginfo("flash sukses" + sn)
	fmt.Println("OpenOCD eksekusi berhasil!", filePath)
	tambahbaris(data.Bord, sn, data.Modelver)
	Updatestatus(sn, data.Bord, "SUCCESS")
	return nil
}

func ExecuteOpenOCDble(data Bus) error {
	// filePath := "./bin/hmi/bootloader.bin"
	// filePath1 := "./bin/hmi/application.bin"

	filePath := fmt.Sprintf("./bin/vcu/%d/ble-bootloader.bin", data.Model)
	filePath1 := fmt.Sprintf("./bin/vcu/%d/ble-application.bin", data.Model)
	// datasn1, datasn2, sn := SNKeyless(data)
	// Periksa apakah file ada
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("File gak ada")
		// Updatestatus(sn, data.Bord, "ERROR")
		Dialogeror("Gagal Flash" + filePath)
		return fmt.Errorf("file tidak ditemukan: %s", filePath)
	}

	// datasn3 := rearrangeBytes(datasn2)

	// Perintah untuk menjalankan OpenOCD

	cmd := exec.Command(".\\openocd\\bin\\openocd.exe",
		"-f", ".\\openocd\\share\\openocd\\scripts\\interface\\jlink.cfg",
		"-f", ".\\openocd\\share\\openocd\\scripts\\target\\rsl10.cfg",
		"-c", fmt.Sprintf("adapter speed 1000; flash init; init; halt; flash protect 0 0 2 off; flash erase_sector 0 0 191; program %s 0x00100000; halt; program %s 0x00107000; reset; exit;", filePath, filePath1))

	// Jalankan perintah dan ambil output
	output, err := cmd.CombinedOutput()
	if err != nil {

		// Updatestatus(sn, data.Bord, "ERROR")
		Dialogeror("Gagal Flash" + filePath)
		fmt.Println("Output:", string(output))
		fmt.Println(cmd.String())
		return fmt.Errorf("eksekusi OpenOCD gagal: %v. Output: %s", err, string(output))
	}

	fmt.Println("OpenOCD eksekusi berhasil!")
	fmt.Println("Output:", string(output))
	Dialoginfo("flash sukses RSL")
	fmt.Println("OpenOCD eksekusi berhasil!", filePath)
	// tambahbaris(data.Bord, sn, data.Modelver)
	// Updatestatus(sn, data.Bord, "SUCCESS")
	return nil
}
