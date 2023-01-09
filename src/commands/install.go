package commands

import (
	"fmt"
	"fpm/src/types"
	utils "fpm/src/utils"
	"strings"
)

func prepareName(name string) (string, string, string) {
	var pkgName = name
	var version = ""
	var author = ""

	if strings.Contains(pkgName, "@") {
		pkgName = strings.ToLower(strings.Split(pkgName, "@")[0])
		version = strings.Split(pkgName, "@")[1]
	}

	if strings.Contains(pkgName, ":") {
		pkgName = strings.ToLower(strings.Split(pkgName, ":")[0])
		author = strings.ToLower(strings.Split(pkgName, ":")[1])
	}

	return pkgName, author, version
}

func InstallCommand(args []string) error {
	if len(args) == 0 {
		utils.Println("<yellow>Usage:</yellow> <green>%s</green> install <yellow><package/theme>(@<version>:<author>)</yellow>", utils.GetExecutableName())
		fmt.Println("")
		utils.Println("<yellow>Description:</yellow>")
		utils.Println("  Install a package or theme")
		fmt.Println("")
		utils.Println("<yellow>Examples:</yellow>")
		utils.Println("  <green>%s</green> install <yellow>git</yellow>", utils.GetExecutableName())
		utils.Println("  <green>%s</green> install <yellow>catpuccino@1.2.3</yellow>", utils.GetExecutableName())
		utils.Println("  <green>%s</green> install <yellow>catpuccino:catpuccino</yellow>", utils.GetExecutableName())
		utils.Println("  <green>%s</green> install <yellow>catpuccino@1.2.3:catpuccino</yellow>", utils.GetExecutableName())
		return nil
	}

	pkgName, author, version := prepareName(args[0])

	installedPackageCache, err := types.GetInstalledPackageCache()
	if err != nil {
		return err
	}

	installedPackage, ok := installedPackageCache.Packages[pkgName]
	if !ok {
		installedPackage, ok = installedPackageCache.Themes[pkgName]
		if ok && installedPackage == version {
			utils.Println("<red>Theme</red> <blue>%s</blue> <red>is already installed</red>", pkgName)
			return utils.ErrPackageAlreadyInstalled
		}
	} else if installedPackage == version {
		utils.Println("<red>Package</red> <blue>%s</blue> <red>is already installed</red>", pkgName)
		return utils.ErrPackageAlreadyInstalled
	}

	packageCache, err := types.GetPackageCache()
	if err != nil {
		return err
	}

	pkg, err := packageCache.GetPackage(pkgName, author)
	if err != nil {
		if err != utils.ErrPackageNotFound {
			// Fall through
		} else if err == utils.ErrPackageFoundMultiple {
			utils.Println("<red>Multiple packages found with name '%s'</red>", pkgName)
			utils.Println("Please specify the author with <yellow>%s</yellow> or <yellow>%s</yellow>",
				pkgName+":<author>", pkgName+"@<version>:<author>")
			return err
		} else {
			utils.Println("<red>Error occurred: %s</red>", err.Error())
			return err
		}
	}

	if pkg == nil {
		pkg, err = packageCache.GetTheme(pkgName, author)
		if err != nil {
			if err != utils.ErrPackageNotFound {
				utils.Println("<red>No package or theme found with name '%s'</red>", pkgName)
			} else if err == utils.ErrPackageFoundMultiple {
				utils.Println("<red>Multiple packages found with name '%s'</red>", pkgName)
				utils.Println("Please specify the author with <yellow>%s</yellow> or <yellow>%s</yellow>",
					pkgName+":<author>", pkgName+"@<version>:<author>")
			} else {
				utils.Println("<red>Error occurred: %s</red>", err.Error())
			}
			return err
		}
	}

	if !utils.AskYesNo(utils.ColorizeHtml("Install <blue>%s</blue>?", pkg.Name)) {
		utils.Println("<red>Aborting...</red>")
		return nil
	}

	err = pkg.Install(version)
	if err != nil {
		if err == utils.ErrPackageNotFound {
			utils.Println("<red>Package not found</red>")
		} else if err == utils.ErrPackageAlreadyInstalled {
			utils.Println("<red>Package</red> <blue>%s</blue> <red>is already installed</red>", pkg.Name)
		} else {
			utils.Println("<red>Error occurred: %s</red>", err.Error())
			return err
		}
	}

	return nil
}
