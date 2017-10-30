package config

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

type closedError struct {
	flockErr error
	fileErr  error
}

func (ce *closedError) Error() string {
	return fmt.Sprintf("%s, %s", ce.fileErr.Error(), ce.flockErr.Error())
}

func newClosedError(flockErr, fileErr error) error {
	if flockErr == nil && fileErr == nil {
		return nil
	}

	if fileErr == nil {
		fileErr = errors.New("no file errors")
	}

	if flockErr == nil {
		flockErr = errors.New("no lock errors")
	}

	return &closedError{flockErr, fileErr}
}

func createOrOpenLockedFile(name string) (file *os.File, err error) {
	if _, serr := os.Stat(name); os.IsNotExist(serr) {
		file, err = os.Create(name)
	} else {
		file, err = os.OpenFile(name, os.O_RDWR, 0644)
	}
	if err != nil {
		return
	}

	if ferr := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); ferr != nil {
		err = ferr
	}

	return
}

func closeLockedFile(file *os.File) error {
	flockErr := syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
	fileErr := file.Close()

	return newClosedError(flockErr, fileErr)
}
