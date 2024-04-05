/*
Validates file hashes from steamdb

Usage:

	sdbchk <csv> <directory>

Where:

	<csv>

is the extracted csv file path from
"https://steamdb.info/depot/<id>/?show_hashes", with the format:
"Name,SHA1Hash\n".

At the time of writting this, the hashses shown there only emits the first and
last 10 characters of the hash, and the project is accomodated for that.

Also you can only get the latest checksums of the game through steamdb afaik

	<directory>

is the root of the game folder

(e.g. "E:\Games\Steam Library\steamapps\common\Half Life 3")
*/
package main

import (
	"github.com/arsmoriendy/sdbchk/pkg/sdbchk"
	"os"
)

func main() {
	sdbchk.SdbChk(os.Args[1], os.Args[2])
}
