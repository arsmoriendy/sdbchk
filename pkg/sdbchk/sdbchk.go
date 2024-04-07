package sdbchk

import (
	"crypto/sha1"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"github.com/arsmoriendy/sdbchk/internal/eprintf"
	"os"
)

func eOpenFile(filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		eprintf.EPrintf("Failed to open %v: %v", filename, err.Error())
	}
	return file
}

func SdbChk(csvFileName string, chckDir string) {
	csvFile := eOpenFile(csvFileName)
	defer csvFile.Close()

	r := csv.NewReader(csvFile)
	r.Read() // skip first line (csv header)
	invalidCount := 0
	for row, err := r.Read(); row != nil; row, err = r.Read() {
		if err != nil {
			csvFile.Close()
			eprintf.EPrintf("Failed to read %v: %v", csvFileName, err.Error())
		}

		csvName, csvSum := row[0], row[1]

		if csvSum == "" {
			continue
		}

		fmt.Printf("Checking \"%v\"\t", csvName)

		filename := chckDir + "/" + csvName
		fileBytes, err := os.ReadFile(filename)
		if err != nil {
			csvFile.Close()
			eprintf.EPrintf("Failed to read %v: %v", filename, err.Error())
		}

		sum := sha1.Sum(fileBytes)
		hs := hex.EncodeToString(sum[:])

		if csvSum[:10] == hs[:10] && csvSum[13:] == hs[30:] {
			fmt.Println("valid")
		} else {
			fmt.Println("invalid")
			invalidCount++
		}
	}

	fmt.Printf("Found %v invalid sums", invalidCount)
}
