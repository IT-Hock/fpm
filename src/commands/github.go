package commands

import (
	"fpm/src/build"
	"fpm/src/utils"
	"os"
)

func GithubCommand(args []string) error {
	if len(args) == 0 {
		utils.Println("Usage: <green>%s</green> github (login|logout|token)", utils.GetExecutableName())
		utils.Println("<white>This is an alternative to use environment variables.</white>")
		utils.Println("<white>It uses the secret storage of your OS to securely save the token.</white>")

		return nil
	}

	switch args[0] {
	case "login":
		if build.GithubClientId == "" || build.GithubClientSecret == "" {
			utils.Println("<red>GitHub login is not supported in this build.</red>")
			utils.Println("Please use <white>%s</white> <green>github</green> <yellow>token <token></yellow> to set your token manually.",
				utils.GetExecutableName())
			return nil
		}
		return GithubLoginCommand(args[1:])

	case "logout":
		return GithubLogoutCommand(args[1:])

	case "token":
		return GithubTokenCommand(args[1:])
	}
	return nil
}

func GithubTokenCommand(args []string) error {
	if len(args) == 0 {
		config := utils.GetConfig()
		token := config.GithubToken
		if token == "" {
			utils.Println("<red>You are not logged in.</red>")
			return nil
		}

		token = utils.ObfuscateString(token, len(token)/2)
		utils.Println("Your GitHub token is: <green>%s</green>", token)
		return nil
	}

	token := args[0]

	// Check if the token is valid
	_, err := utils.GithubGetUser(token)
	if err != nil {
		utils.Println("<red>The specified token is invalid.</red>")
		return nil
	}

	err = utils.SetGithubToken(token)
	if err != nil {
		return err
	}

	utils.Println("<green>You have been logged in.</green>")
	return nil
}

func GithubLogoutCommand([]string) error {
	config := utils.GetConfig()
	token := config.GithubToken
	if token == "" {
		utils.Println("<red>You are not logged in.</red>")
		return nil
	}

	err := utils.SetGithubToken("")
	if err != nil {
		return err
	}

	utils.Println("<green>You have been logged out.</green>")
	return nil
}

func GithubLoginCommand([]string) error {
	config := utils.GetConfig()
	if config.GithubToken != "" {
		utils.Println("<red>You are already logged in.</red>")
		return nil
	}

	code, err := utils.GithubGetDeviceCode()
	if err != nil {
		utils.Println("<red>Failed to get device code: %s</red>", err.Error())
		os.Exit(1)
	}

	utils.Println("<yellow>Open</yellow> <blue>%s</blue> <yellow>and enter the code</yellow> <blue>%s</blue>",
		"https://github.com/login/device", code.UserCode)

	tokenAnswer, err := utils.GithubGetToken(code)
	if err != nil {
		utils.Println("<red>Failed to get token: %s</red>", err.Error())
		os.Exit(1)
	}

	err = utils.SetGithubToken(tokenAnswer)
	if err != nil {
		return err
	}

	utils.Println("<green>You have been logged in.</green>")

	return nil
}
