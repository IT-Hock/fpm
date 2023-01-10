package types

import (
	"context"
	"fmt"
	"fpm/src/utils"
	"github.com/google/go-github/v49/github"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Package struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Image        string   `json:"image"`
	Tags         []string `json:"tags"`
	Author       string   `json:"author"`
	Repository   string   `json:"repository"`
	Dependencies []string `json:"dependencies"`

	IsTheme bool `json:"-"`
}

func (p *Package) Equals(other *Package) bool {
	return p.Name == other.Name
}

func (p *Package) Fetch(version string) (string, string, error) {
	if strings.Contains(p.Repository, "gitlab.com") {
		return p.fetchFromGitLab(version)
	}

	if strings.Contains(p.Repository, "github.com") {
		v, t, err := p.fetchFromGitHub(version)
		if err != nil {
			return "", "", err
		}
		return v, t, nil
	}

	return "", "", utils.ErrPackageNotFound
}

// fetchFromGitLab fetches a package from GitLab returns the download url and the version
func (p *Package) fetchFromGitLab(version string) (string, string, error) {
	compile, err := regexp.Compile(`^https?://(www\.)?gitlab\.com/([\w-_]+)/([\w-_]+)$`)
	if err != nil {
		return "", "", fmt.Errorf("failed to compile regex: %s", err.Error())
	}

	matches := compile.FindStringSubmatch(p.Repository)
	if len(matches) < 3 {
		return "", "", fmt.Errorf("invalid repository url: %s", p.Repository)
	}

	owner := matches[2]
	repo := matches[3]

	client, err := utils.GetGitLabClient()
	if err != nil {
		return "", "", err
	}

	// Find the project
	project, _, err := client.Projects.GetProject(owner+"/"+repo, nil)
	if err != nil {
		return "", "", err
	}

	match, err := regexp.MatchString(`^v\d+\.\d+\.\d+$`, version)
	if err != nil || !match {
		return project.WebURL + "/-/archive/master/" + project.Name + "-" + project.DefaultBranch + ".zip", project.DefaultBranch, nil
	}

	// Find the latest release
	releases, _, err := client.Releases.ListReleases(project.ID, nil)
	if err != nil {
		return "", "", err
	}

	// Find the release with the given version
	for _, release := range releases {
		if release.TagName == version || release.TagName == version[1:] {
			return release.Assets.Links[0].DirectAssetURL, release.TagName, nil
		}
	}

	return "", "", utils.ErrPackageNotFound
}

// fetchFromGitHub fetches a package from GitHub returns the download url and the version
func (p *Package) fetchFromGitHub(version string) (string, string, error) {
	// With optional www
	compile, err := regexp.Compile(`^https?://(www\.)?github\.com/([\w-_]+)/([\w-_]+)$`)
	if err != nil {
		return "", "", fmt.Errorf("failed to compile regex: %s", err.Error())
	}
	matches := compile.FindStringSubmatch(p.Repository)
	if len(matches) < 3 {
		return "", "", fmt.Errorf("invalid repository url: %s", p.Repository)
	}
	owner := matches[2]
	repo := matches[3]

	client := utils.GetGithubClient()
	isLimited, err, _ := utils.GithubHasRateLimit(client)
	if isLimited {
		return "", "", err
	}

	repository, _, err := client.Repositories.Get(context.Background(), owner, repo)
	if err != nil {
		return "", "", err
	}

	match, err := regexp.MatchString(`^v\d+\.\d+\.\d+$`, version)
	if err != nil || !match {
		release, _, _ := client.Repositories.GetLatestRelease(context.Background(), owner, repo)
		if release != nil {
			return release.GetZipballURL(), release.GetTagName(), nil
		}

		if version == "" || repository.GetDefaultBranch() == version || version == "latest" {
			return "https://github.com/" + owner + "/" + repo + "/archive/refs/heads/" + *repository.DefaultBranch + ".zip", *repository.DefaultBranch, nil
		}

		return "https://github.com/" + owner + "/" + repo + "/archive/refs/heads/" + version + ".zip", version, nil
	}

	releases, _, err := client.Repositories.ListReleases(context.Background(), owner, repo, nil)
	if err != nil {
		return "", "", err
	}

	if version != "" {
		for _, release := range releases {
			if release.GetTagName() == version[1:] || release.GetName() == version {
				return release.GetZipballURL(), release.GetTagName(), nil
			}
		}
	}

	return "", "", utils.ErrPackageNotFound
}

func (p *Package) IsInstalled(version string) bool {
	if version == "" {
		return p.IsInstalled("latest")
	}

	if version == "latest" {
		_, s, err := p.Fetch(version)
		if err != nil {
			if rateLimit, ok := err.(*github.RateLimitError); ok {
				// Convert err.Error() to github.RateLimitError
				var resetTime = rateLimit.Rate.Reset.Time.Format("2006-01-02 15:04:05")
				var secondsLeft = strconv.FormatFloat(time.Until(rateLimit.Rate.Reset.Time).Seconds(), 'f', 0, 64)
				configDirectory, err := utils.GetConfigDirectory()
				if err != nil {
					configDirectory = "~/.config"
				}

				utils.Println("<red>GitHub rate limit exceeded.</red>")
				utils.Println("<yellow>Possible ways to fix this:</yellow>")
				utils.Println("\t<yellow>- Add a GitHub token (https://github.com/settings/tokens/new) to your %s/fpm/.fpmrc</yellow>", configDirectory)
				utils.Println("\t<yellow>- Add a GitHub token (https://github.com/settings/tokens/new) to your environment variables (FPM_GITHUB_TOKEN)</yellow>")
				utils.Println("\t<yellow>- Wait %s second(s) (%s) for the rate limit to reset</yellow>", secondsLeft, resetTime)

				if utils.AskYesNo("<yellow>Do you want to generate a token now?</yellow>") {
					code, err := utils.GithubGetDeviceCode()
					if err != nil {
						utils.Println("<red>Failed to get device code: %s</red>", err.Error())
						os.Exit(1)
					}

					utils.Println("<yellow>Open</yellow> <blue>%s</blue> <yellow>and enter the code</yellow> <blue>%s</blue>",
						"https://github.com/login/device", code.UserCode)

					_, err = utils.GithubGetToken(code)
					if err != nil {
						utils.Println("<red>Failed to get token: %s</red>", err.Error())
						os.Exit(1)
					}
				}

				os.Exit(1)
			}
			panic(err)
			return false
		}

		version = s
	}

	// First check if there is a folder
	installPath, err := p.GetInstallPath(version)
	if err != nil {
		return false
	}

	if utils.Exists(installPath) {
		return true
	}

	return false
}

func (p *Package) Empty() bool {
	return p.Name == "" && p.Description == "" && p.Image == "" && p.Tags == nil && p.Author == "" && p.Repository == "" && p.Dependencies == nil
}

// GetInstallPath returns the install path of the package (e.g. ~/.fpm/packages/<name>/<version>)
func (p *Package) GetInstallPath(version string) (string, error) {
	packageDirectory, err := utils.GetPackageDirectory()
	if err != nil {
		return "", err
	}

	if version == "" || version == "latest" {
		return path.Join(packageDirectory, p.Name, "latest"), nil
	}

	return path.Join(packageDirectory, p.Name, version), nil
}

func (p *Package) Install(version string) error {
	if version == "" {
		return p.Install("latest")
	}

	if version == "latest" {
		_, s, err := p.Fetch(version)
		if err != nil {
			panic(err)
		}

		version = s
	}

	if p.IsInstalled(version) {
		return utils.ErrPackageAlreadyInstalled
	}

	installPath, err := p.GetInstallPath(version)
	if err != nil {
		panic(err)
	}

	if utils.Exists(installPath) {
		utils.Println("<yellow>Package</yellow> <blue>%s</blue> <yellow>was not uninstalled properly</yellow>", p.Name)
		return utils.ErrPackageLeftoverFiles
	}

	zipballURL, _, err := p.Fetch(version)
	if err != nil {
		utils.Println("<red>Failed to fetch package</red> <blue>%s</blue> <red>version</red> <blue>%s</blue>", p.Name, version)
		os.Exit(1)
	}

	installPath = path.Join(installPath, "/")

	if !utils.Exists(installPath) {
		err := utils.CreateDirectory(installPath, 0755)
		if err != nil {
			panic(err)
		}
	}
	filePath := path.Join(installPath, "package.zip")

	// Download the zipball
	utils.Println("Downloading <blue>%s</blue>@%s", p.Name, version)
	err = utils.DownloadFile(zipballURL, filePath)
	if err != nil {
		utils.Println("<red>Failed to download package</red>")
		utils.Println("<red>%s</red>", err.Error())
		os.Exit(1)
	}

	// Unzip the zipball
	utils.Println("Unzipping downloaded file..")
	_, err = utils.Unzip(filePath, installPath)
	if err != nil {
		utils.Println("<red>Failed to unzip downloaded file: %s</red>", err.Error())
		os.Exit(1)
	}

	// Remove the zipball
	utils.Println("Removing temporary file..")
	err = os.Remove(filePath)
	if err != nil {
		utils.Println("<red>Failed to remove temporary file: %s</red>", err.Error())
		utils.Println("<yellow>Please delete the temporary file:</yellow> <blue>%s</blue>", filePath)
		os.Exit(1)
	}

	installedPackageCache, err := GetInstalledPackageCache()
	if err != nil {
		utils.Println("<red>Failed to get installed package cache: %s</red>", err.Error())
		os.Exit(1)
	}

	if p.IsTheme {
		installedPackageCache.Themes[p.Name] = version
	} else {
		installedPackageCache.Packages[p.Name] = version
	}

	err = installedPackageCache.Save()
	if err != nil {
		return err
	}

	return nil
}

func (p *Package) Uninstall(version string) error {
	if version == "" {
		return p.Uninstall("latest")
	}

	if version == "latest" {
		_, s, err := p.Fetch(version)
		if err != nil {
			panic(err)
		}

		version = s
	}

	if !p.IsInstalled(version) {
		return utils.ErrPackageNotInstalled
	}

	installPath, err := p.GetInstallPath(version)
	if err != nil {
		panic(err)
	}

	err = os.RemoveAll(installPath)
	if err != nil {
		utils.Println("<red>Failed to remove package directory: %s</red>", err.Error())
		os.Exit(1)
	}

	installedPackageCache, err := GetInstalledPackageCache()
	if err != nil {
		utils.Println("<red>Failed to get installed package cache: %s</red>", err.Error())
		os.Exit(1)
	}

	if p.IsTheme {
		delete(installedPackageCache.Themes, p.Name)
	} else {
		delete(installedPackageCache.Packages, p.Name)
	}

	err = installedPackageCache.Save()
	if err != nil {
		return err
	}

	return nil
}
