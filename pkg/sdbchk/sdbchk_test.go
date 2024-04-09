package sdbchk_test

import (
	"github.com/arsmoriendy/sdbchk/pkg/sdbchk"
	"os"
	fp "path/filepath"
	"testing"
)

func TestCheckExtra(t *testing.T) {
	testCheckExtra(t, "../../test/data/foo.csv", "../../test/data/foo", []string{})
	testCheckExtra(t, "../../test/data/bar.csv", "../../test/data/bar", []string{"foo/baz"})
	testCheckExtra(t, "../../test/data/baz.csv", "../../test/data/baz", []string{})
}

func testCheckExtra(t *testing.T, csvFileName, chkDir string, expectedFns []string) {
	// define vars
	wd, _ := os.Getwd()
	absCsvFileName := fp.Join(wd, fp.ToSlash(csvFileName))
	absChkDir := fp.Join(wd, fp.ToSlash(chkDir))

	// run check extra
	fns := sdbchk.CheckExtra(absCsvFileName, absChkDir)

	// checks
	fnsLen, expectedFnsLen := len(fns), len(expectedFns)
	if fnsLen != expectedFnsLen {
		t.Fatalf(
			"Wrong number of extra files, expected %v, found %v\n%v",
			expectedFnsLen,
			fnsLen,
			fns,
		)
	}

	extraCount := 0
	for i, fn := range fns {
		expectedFn := expectedFns[i]
		if expectedFn != fn {
			t.Logf(
				"Extra file \"%v\" is not at the expected index: %v. Expected \"%v\"",
				fn,
				i,
				expectedFn,
			)
			extraCount++
		}
	}
	if extraCount > 0 {
		t.FailNow()
	}
}
