package commands

import (
	"fpm/src/types"
	"fpm/src/utils"
)

func RemoveCommand(args []string) error {
	packageCache, err := types.GetPackageCache()
	if err != nil {
		return err
	}

	installedPackageCache, err := types.GetInstalledPackageCache()
	if err != nil {
		return err
	}

	installedPackages, err := getInstalledPackages(packageCache, installedPackageCache)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		utils.Println("<red>Missing argument: package name</red>")
		return nil
	}

	packageName := args[0]

	installedPackage := installedPackages[packageName]
	if installedPackage == nil {
		utils.Println("<red>Package is not installed.</red>")
		return nil
	}

	err = installedPackage.Uninstall()
	if err != nil {
		return err
	}

	return nil
}
