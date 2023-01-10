package types

import (
	"encoding/json"
	"fpm/src/utils"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

type InstalledPackages struct {
	Packages map[string]string `json:"packages"`
	Themes   map[string]string `json:"themes"`
}

type Packages struct {
	Packages map[string]Package `json:"packages"`
	Themes   map[string]Package `json:"themes"`
}

var packageCache *Packages
var installedPackageCache *InstalledPackages

func GetInstalledPackageCache() (*InstalledPackages, error) {
	if installedPackageCache == nil {
		installedPackageCache = &InstalledPackages{}
		err := installedPackageCache.load()
		if err != nil {
			return nil, err
		}
	}

	return installedPackageCache, nil
}

func (p *InstalledPackages) Save() error {
	cacheDir, _ := utils.GetCacheDirectory()

	file, err := os.Create(cacheDir + "/fpm/.installed.json")
	if err != nil {
		return err
	}

	err = json.NewEncoder(file).Encode(p)
	if err != nil {
		return err
	}

	// Write updated packages.fish
	err = p.LinkPackages()
	if err != nil {
		return err
	}

	return nil
}

func (p *InstalledPackages) load() error {
	err := utils.CheckPackageCache()
	if err != nil {
		return err
	}
	cacheDir, _ := utils.GetCacheDirectory()

	file, err := os.Open(cacheDir + "/fpm/.installed.json")
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(cacheDir + "/fpm/.installed.json")
			if err != nil {
				return err
			}
			err := json.NewEncoder(file).Encode(p)
			if err != nil {
				return err
			}
			p.Themes = make(map[string]string)
			p.Packages = make(map[string]string)
			return nil
		}
		return err
	}

	err = json.NewDecoder(file).Decode(&p)
	if err != nil {
		return err
	}

	return nil
}

func (p *InstalledPackages) LinkPackages() error {
	directory, err := utils.GetFpmDirectory()
	if err != nil {
		return err
	}

	confDir := path.Join(directory, "enabled")
	if _, err := os.Stat(confDir); os.IsNotExist(err) {
		err = os.Mkdir(confDir, 0755)
		if err != nil {
			return err
		}
	}

	for pkg, version := range p.Packages {
		packageDir := path.Join(directory, "packages", pkg, version)
		if _, err := os.Stat(path.Join(confDir, pkg+"@"+version)); !os.IsNotExist(err) {
			continue
		}

		err := os.Symlink(packageDir, path.Join(confDir, pkg+"@"+version))
		if err != nil {
			panic(err)
			return err
		}
	}

	return nil
}

func GetPackageCache() (*Packages, error) {
	if packageCache == nil {
		if err := utils.CheckPackageCache(); err != nil {
			return nil, err
		}

		packageCache = &Packages{}
		err := packageCache.load()
		if err != nil {
			return nil, err
		}
	}

	return packageCache, nil
}

func (p *Packages) load() error {
	err := utils.CheckPackageCache()
	if err != nil {
		return err
	}
	cacheDir, _ := utils.GetCacheDirectory()

	file, err := os.Open(cacheDir + "/fpm/packages.json")
	if err != nil {
		return err
	}

	err = json.NewDecoder(file).Decode(&p)
	if err != nil {
		return err
	}

	return nil
}

func (p *Packages) FindPackage(name string) (*Package, error) {
	if pkg, ok := p.Packages[name]; ok {
		return &pkg, nil
	}

	if pkg, ok := p.Themes[name]; ok {
		return &pkg, nil
	}

	return nil, utils.ErrPackageNotFound
}

func (p *Packages) FindPackagesFast(query string) ([]Package, error) {
	var packages []Package

	for k, pkg := range p.Packages {
		if strings.HasPrefix(k, query) {
			packages = append(packages, pkg)
		}
	}

	if len(packages) == 0 {
		return []Package{}, utils.ErrPackageNotFound
	}

	return packages, nil
}

func (p *Packages) FindThemesFast(query string) ([]Package, error) {
	var packages []Package

	for k, pkg := range p.Themes {
		if strings.HasPrefix(k, query) {
			packages = append(packages, pkg)
		}
	}

	if len(packages) == 0 {
		return []Package{}, utils.ErrPackageNotFound
	}

	return packages, nil
}

// FindPackages slower than FindPackagesFast, but it's more accurate
func (p *Packages) FindPackages(expression *regexp.Regexp) ([]Package, error) {
	var packages []Package

	for k, pkg := range p.Packages {
		if expression.MatchString(k) {
			packages = append(packages, pkg)
		}
	}

	for k, pkg := range p.Themes {
		if expression.MatchString(k) {
			packages = append(packages, pkg)
		}
	}

	if len(packages) == 0 {
		return []Package{}, utils.ErrPackageNotFound
	}

	return packages, nil
}

func (p *Packages) GetPackage(name string, author string) (*Package, error) {
	if pkg, ok := p.Packages[name]; ok {
		if author == "" || pkg.Author == author {
			pkg.IsTheme = false
			return &pkg, nil
		}
	}

	return nil, utils.ErrPackageNotFound
}

func (p *Packages) GetTheme(name string, author string) (*Package, error) {
	if pkg, ok := p.Themes[name]; ok {
		pkg.IsTheme = true
		if author == "" || pkg.Author == author {
			return &pkg, nil
		}
	}

	return nil, utils.ErrPackageNotFound
}

func (p *Packages) Update() error {
	cacheDir, _ := utils.GetCacheDirectory()

	config := utils.GetConfig()
	if config == nil {
		return utils.ErrConfigNotFound
	}

	// Get packages
	resp, err := http.Get(config.PackageRepository)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &p)
	if err != nil {
		return err
	}

	// Write to packages.json
	file, err := os.Create(cacheDir + "/fpm/packages.json")
	if err != nil {
		return err
	}

	defer func(Close func() error) {
		err := Close()
		if err != nil {
			panic(err)
		}
	}(file.Close)

	err = json.NewEncoder(file).Encode(p)
	if err != nil {
		return err
	}

	return nil
}
