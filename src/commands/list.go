package commands

import (
	"fpm/src/types"
	"fpm/src/utils"
)

type InstalledTheme struct {
	Theme            *types.Package
	InstalledVersion string
}

type InstalledPackage struct {
	Package          *types.Package
	InstalledVersion string
}

func (p InstalledPackage) Uninstall() error {
	err := p.Package.Uninstall(p.InstalledVersion)
	if err != nil {
		return err
	}

	return nil
}

func ListCommand(args []string) error {
	packageCache, err := types.GetPackageCache()
	if err != nil {
		return err
	}

	typeFilter := 0
	installedOnly := false
	isQuiet := false
	isAll := false
	if len(args) > 0 {
		for _, arg := range args {
			if arg == "-t=package" || arg == "--type=package" {
				typeFilter = 1
			} else if arg == "-t=theme" || arg == "--type=theme" {
				typeFilter = 2
			} else if arg == "-a" || arg == "--all" {
				isAll = true
			} else if arg == "-i" || arg == "--installed" {
				installedOnly = true
			} else if arg == "-q" || arg == "--quiet" {
				isQuiet = true
			}
		}
		if isAll {
			if installedOnly {
				utils.Println("<red>Cannot use both -a and -i at the same time.</red>")
				return nil
			}
		}
	}

	if installedOnly || !isAll {
		installedPackageCache, err := types.GetInstalledPackageCache()
		if err != nil {
			return err
		}

		installedPackages, err := getInstalledPackages(packageCache, installedPackageCache)
		if err != nil {
			return err
		}

		installedThemes, err := getInstalledThemes(packageCache, installedPackageCache)
		if err != nil {
			return err
		}

		if typeFilter == 0 || typeFilter == 1 {
			if !isQuiet {
				utils.Println("<yellow>Installed Packages</yellow>")
			}
			stringBuffer := ""

			i := 0
			for _, pkg := range installedPackages {
				if isQuiet {
					stringBuffer += pkg.Package.Name + "\t" + pkg.InstalledVersion
					// Check if it's the last element.
					if i != len(installedThemes)-1 {
						stringBuffer += "\n"
					}
				} else {
					description := pkg.Package.Description
					if description == "" {
						description = utils.ColorizeHtml("<white>No description</white>")
					}
					utils.Println("\t%-50s - %s",
						utils.ColorizeHtml("<blue>%s</blue>@<yellow>%s</yellow>", pkg.Package.Name, pkg.InstalledVersion),
						description)
				}
				i++
			}

			if isQuiet {
				utils.Println(stringBuffer)
			}
		}

		if typeFilter == 0 || typeFilter == 2 {
			if len(installedThemes) > 0 {
				if len(installedPackages) > 0 {
					utils.Println("")
				}

				if !isQuiet {
					utils.Println("<yellow>Installed Themes</yellow>")
				}

				stringBuffer := ""

				i := 0
				for _, theme := range installedThemes {
					if isQuiet {
						stringBuffer += theme.Theme.Name + "\t" + theme.InstalledVersion
						// Check if it's the last element.
						if i != len(installedThemes)-1 {
							stringBuffer += "\n"
						}
					} else {
						description := theme.Theme.Description
						if description == "" {
							description = utils.ColorizeHtml("<white>No description</white>")
						}
						utils.Println("\t%-50s - %s",
							utils.ColorizeHtml("<blue>%s</blue>@<yellow>%s</yellow>", theme.Theme.Name, theme.InstalledVersion),
							description)
					}
					i++
				}

				if isQuiet {
					utils.Println(stringBuffer)
				}
			}
		}
	} else if isAll {

		if typeFilter == 0 || typeFilter == 1 {
			if !isQuiet {
				utils.Println("<yellow>Packages</yellow>")
			}

			stringBuffer := ""

			i := 0
			for _, pkg := range packageCache.Packages {
				description := pkg.Description
				if isQuiet {
					if description == "" {
						description = "No description"
					}
					stringBuffer += pkg.Name + "\t" + description
					// Check if it's the last element.
					if i != len(packageCache.Packages)-1 {
						stringBuffer += "\n"
					}
				} else {
					if description == "" {
						description = utils.ColorizeHtml("<white>No description</white>")
					}
					utils.Println("\t%-50s - %s",
						utils.ColorizeHtml("<blue>%s</blue>", pkg.Name),
						description)
				}
				i++
			}

			if isQuiet {
				utils.Println(stringBuffer)
			}
		}

		if typeFilter == 0 || typeFilter == 2 {
			if len(packageCache.Themes) > 0 {
				if len(packageCache.Packages) > 0 {
					utils.Println("")
				}

				if !isQuiet {
					utils.Println("<yellow>Themes</yellow>")
				}

				stringBuffer := ""

				i := 0
				for _, theme := range packageCache.Themes {
					description := theme.Description
					if isQuiet {
						if description == "" {
							description = "No description"
						}
						stringBuffer += theme.Name + "\t" + description
						// Check if it's the last element.
						if i != len(packageCache.Themes)-1 {
							stringBuffer += "\n"
						}
					} else {
						if description == "" {
							description = utils.ColorizeHtml("<white>No description</white>")
						}

						utils.Println("\t%-50s - %s",
							utils.ColorizeHtml("<blue>%s</blue>", theme.Name),
							description)
					}
					i++
				}

				if isQuiet {
					utils.Println(stringBuffer)
				}
			}
		}
	}

	return nil
}

func getInstalledPackages(packageCache *types.Packages, installedPackageCache *types.InstalledPackages) (map[string]*InstalledPackage, error) {
	installedPackages := installedPackageCache.Packages

	var packages = make(map[string]*InstalledPackage)

	for name, version := range installedPackages {
		pkg, err := packageCache.GetPackage(name, "")
		if err != nil {
			return nil, err
		}

		packages[name] = &InstalledPackage{
			Package:          pkg,
			InstalledVersion: version,
		}
	}

	return packages, nil
}

func getInstalledThemes(packageCache *types.Packages, installedPackageCache *types.InstalledPackages) (map[string]*InstalledTheme, error) {
	installedThemes := installedPackageCache.Themes

	var themes = make(map[string]*InstalledTheme)

	for name, version := range installedThemes {
		theme, err := packageCache.GetTheme(name, "")
		if err != nil {
			return nil, err
		}

		themes[name] = &InstalledTheme{
			Theme:            theme,
			InstalledVersion: version,
		}
	}

	return themes, nil
}
