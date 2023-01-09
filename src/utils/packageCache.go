package utils

import (
	"os"
)

func CheckPackageCache() error {
	// Check cache for package
	cacheDir, err := GetCacheDirectory()
	if err != nil {
		return ErrPackageCacheNotFound
	}

	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		return ErrPackageCacheNotFound
	}

	if _, err := os.Stat(cacheDir + "/fpm"); os.IsNotExist(err) {
		return ErrPackageCacheNotFound
	}

	if _, err := os.Stat(cacheDir + "/fpm/packages.json"); os.IsNotExist(err) {
		return ErrPackageCacheNotFound
	}

	return nil
}
