package main

import (
	"os"
	"path/filepath"

	"github.com/electric-saw/kafta/pkg/cmd/kafta"
	"github.com/electric-saw/kafta/pkg/cmd/util"
)

func main() {
	baseName := filepath.Base(os.Args[0])

	err := kafta.NewKaftaCommand(baseName).
		Execute()
	util.CheckErr(err)

}
