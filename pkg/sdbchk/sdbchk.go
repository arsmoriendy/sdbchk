package sdbchk

import (
	"crypto/sha1"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"slices"

	"github.com/arsmoriendy/sdbchk/internal/eprintf"
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
	invalidCount, missingCount := 0, 0
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
			if errors.Is(err, os.ErrNotExist) {
				missingCount++
				fmt.Println("missing")
				continue
			}

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

	fmt.Printf("Found %v invalid sums\nFound %v missing files", invalidCount, missingCount)
}

// Extract "Name", "Sha1 Hash" records from a csv.
// csvFn is the csv absolute file path.
// Returns a two dimensional string slice.
// The first dimension of the return slice is the record index.
// The second dimension of the return slice is the field index of each record.
func CsvRecs(csvFn string) [][]string {
	csvFile := eOpenFile(csvFn)
	defer csvFile.Close()

	r := csv.NewReader(csvFile)
	recs, err := r.ReadAll()

	if err != nil {
		eprintf.EPrintf("Failed parsing csv: %v", err.Error())
	}

	return recs
}

// Checks filenames in csvFn that's missing in chckDir.
// csvFn is the absolute path for the csv file.
// chckDir is the absolute dir path to check for files.
// Returns missing filenames
func CheckFiles(csvFn string, chckDir string) []string {
	recs := CsvRecs(csvFn)

	missingFns := []string{}
	for _, rec := range recs {
		filename := chckDir + "/" + rec[0]
		f, err := os.Open(filename)
		if err != nil {
			missingFns = append(missingFns, filename)
		}
		f.Close()
	}

	return missingFns
}

// Checks for extra dir/files that isn't on the csv file.
// csvFn is the absolute csv file path.
// chkDir is the absolute directory path on which to check extras for.
// Returns a slice of extra dir/filename strings
func CheckExtra(csvFn string, chkDir string) []string {
	// create a sorted array for all file names in the csv
	csvFileNames := []string{}
	for _, rec := range CsvRecs(csvFn) {
		csvFileNames = append(csvFileNames, rec[0])
	}
	slices.Sort(csvFileNames)

	// recurse through all files and subdirectories in absChkDir,
	// find extra files and append to efFns
	efFns := []string{}
	fs.WalkDir(
		os.DirFS(chkDir),
		".",
		func(path string, e fs.DirEntry, err error) error {
			if _, v := slices.BinarySearch(csvFileNames, path); !v {
				if e.IsDir() {
					path = path + "/"
				}
				efFns = append(efFns, path)
			}
			return nil
		},
	)

	// ignore first (0th index) since it's "." (the root directory itself)
	return efFns[1:]
}
