package cdhit

import (
	"errors"
	"io/ioutil"

	"github.com/prvst/philosopher/lib/err"

	"github.com/prvst/philosopher/lib/sys"
)

// Win64 ...
func Win64(win64 string) {

	bin, e := Asset("cd-hit.exe")
	e = ioutil.WriteFile(win64, bin, sys.FilePermission())

	if e != nil {
		err.ExecutingBinary(errors.New("CD-hit"), "trace")
	}

	return
}
