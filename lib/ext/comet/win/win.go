package comet

import (
	"errors"
	"io/ioutil"
)

// WinParameterFile ...
func WinParameterFile(winParam string) error {

	param, err := Asset("comet.params.txt")
	err = ioutil.WriteFile(winParam, param, 0644)

	if err != nil {
		return errors.New("Cannot deploy Comet parameter file")
	}

	return nil
}

// Win32 ...
func Win32(win32 string) error {

	bin, err := Asset("comet.2016012.win32.exe")
	err = ioutil.WriteFile(win32, bin, 0755)

	if err != nil {
		return errors.New("Cannot deploy Comet parameter file")
	}

	return nil
}

// Win64 ...
func Win64(win64 string) error {

	bin, err := Asset("comet.2016012.win64.exe")
	err = ioutil.WriteFile(win64, bin, 0755)

	if err != nil {
		return errors.New("Cannot deploy Comet parameter file")
	}

	return nil
}