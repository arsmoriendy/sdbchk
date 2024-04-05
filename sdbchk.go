package main

import (
	"crypto/sha1"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	sdbchk(os.Args[1], os.Args[2])
}

func sdbchk(csvFileName string, chckDir string) {
	csvFile, err := os.Open(csvFileName)
	if err != nil {
		ePrintf("Failed to open %v: %v", csvFileName, err.Error())
	}
	defer csvFile.Close()

	r := csv.NewReader(csvFile)
	r.Read() // skip first line (csv header)
	invalidCount := 0
	for row, err := r.Read(); row != nil; row, err = r.Read() {
		if err != nil {
			csvFile.Close()
			ePrintf("Failed to read %v: %v", csvFileName, err.Error())
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
			ePrintf("Failed to read %v: %v", filename, err.Error())
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

/**
 * Printf to stderr and exit with -1
 */
func ePrintf(format string, a ...any) {
	fmt.Fprintln(os.Stderr, fmt.Errorf(format, a...))
	os.Exit(-1)
}
