package resources

import (
	"fpm/src/build"
	"fpm/src/utils"
	"io/fs"
	"os"
	"path"
	"strconv"
	"strings"
)

func GetScriptInstallPath() string {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		dataHome = os.Getenv("HOME") + "/.local/share/fpm"
	} else {
		dataHome += "/fpm"
	}

	return dataHome
}

func ExtractSingleFile(file fs.DirEntry) error {
	if _, err := os.Stat(GetScriptInstallPath() + "/" + file.Name()); !os.IsNotExist(err) {
		err = os.Remove(GetScriptInstallPath() + "/" + file.Name())
		if err != nil {
			return err
		}
	}

	fileContent, err := fishFiles.ReadFile("files/" + file.Name())
	if err != nil {
		return err
	}
	err = os.WriteFile(GetScriptInstallPath()+"/"+file.Name(), fileContent, 0644)
	if err != nil {
		return err
	}
	return nil
}

func InstalledScriptsValid() bool {
	// List files in fishFiles
	dir, err := fishFiles.ReadDir(".")
	if err != nil {
		return false
	}

	if _, err := os.Stat(path.Join(GetScriptInstallPath(), ".version")); os.IsNotExist(err) {
		return false
	}

	file, err := os.ReadFile(path.Join(GetScriptInstallPath(), ".version"))
	if err != nil {
		return false
	}

	splittedVersion := strings.Split(string(file), ".")
	if len(splittedVersion) < 3 {
		return false
	} else {
		major, _ := strconv.Atoi(splittedVersion[0])
		minor, _ := strconv.Atoi(splittedVersion[1])
		patch, _ := strconv.Atoi(splittedVersion[2])

		if major < build.VersionMajor ||
			(major == build.VersionMajor && minor < build.VersionMinor) ||
			(major == build.VersionMajor && minor == build.VersionMinor && patch < build.VersionPatch) {
			return false
		}
	}

	for _, resourceFile := range dir {
		if resourceFile.IsDir() {
			continue
		}

		if _, err := os.Stat(GetScriptInstallPath() + "/" + resourceFile.Name()); os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func InstallFpmScripts() error {
	if InstalledScriptsValid() {
		return nil
	}

	dir, err := fishFiles.ReadDir("files")
	if err != nil {
		return err
	}

	err = os.MkdirAll(GetScriptInstallPath(), 0755)
	if err != nil && err != os.ErrExist {
		return err
	}

	// Extract files
	for _, resourceFile := range dir {
		if resourceFile.IsDir() {
			continue
		}

		err = ExtractSingleFile(resourceFile)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(path.Join(GetScriptInstallPath(), ".version"), []byte(build.VersionString), 0644)
	if err != nil {
		return err
	}

	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		configHome = os.Getenv("HOME") + "/.config/fish/conf.d"
	} else {
		configHome += "/fish/conf.d"
	}

	if _, err := os.Stat(configHome); os.IsNotExist(err) {
		err = os.MkdirAll(configHome, 0755)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(path.Join(configHome, "fpm.fish"), []byte("source "+GetScriptInstallPath()+"/init.fish"), 0644)
	if err != nil {
		return err
	}

	// Copy myself to XDG_DATA_HOME / fpm
	executablePath, err := utils.GetExecutablePath()
	if err != nil {
		return err
	}

	fileContent, err := os.ReadFile(path.Join(executablePath, utils.GetExecutableName()))
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(GetScriptInstallPath(), "fpm"), fileContent, 0755)
	if err != nil {
		return err
	}
	if build.Debug {
		utils.PrintDebug("Installed scripts v%s to %s", build.VersionString, GetScriptInstallPath())
	}

	return nil
}
