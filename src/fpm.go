package main

import (
	"flag"
	"fpm/src/build"
	commands2 "fpm/src/commands"
	"fpm/src/types"
	utils "fpm/src/utils"
	"fpm/src/utils/resources"
	"os"
)

func null([]string) error {
	return nil
}

var ConfiguredCommands = []types.Command{
	{
		// This is just here to list it in help
		Command:     "help",
		Aliases:     []string{"h"},
		Name:        "Help",
		Description: "Show this help message",
		Usage:       "",
		Function:    func([]string) error { return nil },
	},
	{
		Command:     "version",
		Aliases:     []string{"v"},
		Name:        "Version",
		Description: "Show the version of fpm",
		Usage:       "",
		Function:    func([]string) error { utils.Println("fpm version %s", build.VersionString); return nil },
	},
	{
		Command:     "install",
		Aliases:     []string{"i"},
		Name:        "Install",
		Description: "Install a package/theme",
		Usage:       "<package/theme>",
		Function:    commands2.InstallCommand,
	},
	{
		Command:     "update",
		Aliases:     []string{"u"},
		Name:        "Update",
		Description: "Update a package/theme",
		Usage:       "<package/theme>",
		Function:    null,
	},
	{
		Command:     "remove",
		Aliases:     []string{"r"},
		Name:        "Remove",
		Description: "Remove a package/theme",
		Usage:       "<package/theme>",
		Function:    commands2.RemoveCommand,
	},
	{
		Command:     "search",
		Aliases:     []string{"s"},
		Name:        "Search",
		Description: "Search for a package/theme",
		Usage:       "<package/theme>",
		Function:    commands2.SearchCommand,
	},
	{
		Command:     "github",
		Aliases:     []string{"gh"},
		Name:        "Github",
		Description: "Manage Github token",
		Usage:       "<login|logout|token>",
		Function:    commands2.GithubCommand,
	},
	{
		Command:     "list",
		Aliases:     []string{"ls"},
		Name:        "List",
		Description: "List installed packages/themes",
		Usage:       "",
		Function:    commands2.ListCommand,
	},
}

func main() {
	utils.GetConfig()
	err := resources.InstallFpmScripts()
	if err != nil {
		panic(err)
	}

	flag.Parse()

	if flag.NArg() == 0 || flag.Arg(0) == "help" || flag.Arg(0) == "h" {
		printHelp()

		os.Exit(1)
	}

	command := flag.Arg(0)
	for _, c := range ConfiguredCommands {
		if c.MatchesCommand(command) {
			var args = flag.Args()
			args = args[1:]
			err := c.Function(args)
			if err != nil {
				utils.Println("<red>%s</red>", err.Error())
				os.Exit(1)
			}
			os.Exit(0)
		}
	}

	utils.Println("<red>Unknown command: %s</red>", command)
	printHelp()

	os.Exit(1)
}

func printHelp() {
	utils.Println("<yellow>Usage:</yellow> <green>%s</green> <yellow><command> [options]</yellow>", utils.GetExecutableName())
	utils.Println("")
	utils.Println("<blue>Commands:</blue>")
	// Print configuredCommands with equal indentation
	for _, command := range ConfiguredCommands {
		utils.Println("  <yellow>%-10s</yellow> <blue>%-20s</blue> - <white>%s</white>", command.Command, command.Usage, command.Description)
	}
	utils.Println("")
	utils.Println("<blue>Options:</blue>")

	flag.VisitAll(func(f *flag.Flag) {
		utils.Println("  <yellow>%-10s</yellow> <blue>%-35s</blue> <white>(Default: %s)</white>", f.Name, f.Usage, f.DefValue)
	})
}
