package modul

import (
	"hash/crc32"
	"io"
	"os"
)

func CalculateCrcAndLen(path string) (uint32, uint32, error) {
	var crc uint32 = 0
	var length uint32 = 0

	f, err := os.Open(path)

	defer func() {
		f.Close()
	}()

	if err != nil {
		return 0, 0, err
	}

	b2 := make([]byte, 256)

	for {
		n2, err := f.Read(b2)
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, 0, err
		}
		length += uint32(n2)

		if n2 == 0 {
			break
		}
		crc = crc32.Update(crc, crc32.IEEETable, b2[:n2])
		if n2 < 256 {
			break
		}
	}
	return crc, length, nil
}
