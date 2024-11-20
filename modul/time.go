package modul

import (
	"fmt"

	"github.com/beevik/ntp"
)

func Timecek() int {
	// Mendapatkan waktu saat ini dari NTP
	currentTime, err := ntp.Time("pool.ntp.org")
	if err != nil {
		fmt.Println("Error getting time from NTP:", err)
		return 0
	}

	// Mendapatkan bulan saat ini dalam format desimal
	month := int(currentTime.Month())

	// Mencetak bulan saat ini dalam format desimal
	return month
}
