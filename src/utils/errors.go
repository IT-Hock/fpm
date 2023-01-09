package utils

import (
	"errors"
	"runtime"
)

var (
	ErrPackageCacheNotFound = errors.New("package cache not found")

	ErrPackageNotFound      = errors.New("package is not found")
	ErrPackageFoundMultiple = errors.New("multiple packages found")

	ErrConfigNotFound = errors.New("config file not found")
	ErrConfigExists   = errors.New("config file already exists")

	ErrSecretNotFound = errors.New("secret not found")

	ErrPackageAlreadyInstalled = errors.New("package is already installed")
	ErrPackageLeftoverFiles    = errors.New("package has leftover files")
	ErrPackageNotInstalled     = errors.New("package is not installed")
)

func StackTrace() string {
	stack := make([]byte, 1024)
	stack = stack[:runtime.Stack(stack, false)]
	return string(stack)
}
