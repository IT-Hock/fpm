package commands

import (
	"fpm/src/types"
	"fpm/src/utils"
)

func SearchCommand(args []string) error {
	packageCache, err := types.GetPackageCache()
	if err != nil {
		return err
	}

	packages, err := packageCache.FindPackagesFast(args[0])
	if err != nil && err != utils.ErrPackageNotFound {
		return err
	}

	themes, err := packageCache.FindThemesFast(args[0])
	if err != nil && err != utils.ErrPackageNotFound {
		return err
	}

	if len(packages) == 0 && len(themes) == 0 {
		utils.Println("<red>Couldn't find any packages or themes matching the search query.</red>")
		return nil
	}

	if len(packages) == 1 && len(themes) == 0 {
		utils.Println("<green>Found 1 package matching the search query: %s</green>", args[0])
		if utils.AskYesNo("Would you like to install it?") {
			return InstallCommand([]string{packages[0].Name})
		}
		return nil
	} else if len(packages) == 0 && len(themes) == 1 {
		utils.Println("<yellow>Found 1 theme matching the search query: %s</yellow>", args[0])
		if utils.AskYesNo("Do you want to install it?") {
			return InstallCommand([]string{themes[0].Name})
		}
		return nil
	}

	if len(packages) > 0 {
		utils.Println("<blue>%d</blue> Packages matching <green>%s</green>:", len(packages), args[0])
		for _, pkg := range packages {
			if pkg.Description == "" {
				utils.Println("\t<green>%-50s</green> <white>No description</white>", pkg.Name)
			} else {
				utils.Println("\t<green>%-50s</green> <blue>%s</blue>", pkg.Name, pkg.Description)
			}
		}
	}

	if len(themes) > 0 {
		if len(packages) > 0 {
			utils.Println("")
		}

		utils.Println("<blue>%d</blue> Themes matching <green>%s</green>:", len(themes), args[0])
		for _, pkg := range themes {
			if pkg.Description == "" {
				utils.Println("\t<green>%-50s</green> <white>No description</white>", pkg.Name)
			} else {
				utils.Println("\t<green>%-50s</green> <blue>%s</blue>", pkg.Name, pkg.Description)
			}
		}
	}

	return nil
}
