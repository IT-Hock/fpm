package types

func (p *Packages) GetPackages() ([]Package, error) {
	var packages []Package
	for _, pkg := range p.Packages {
		packages = append(packages, pkg)
	}
	return packages, nil
}

func (p *Packages) GetThemes() ([]Package, error) {
	var packages []Package
	for _, pkg := range p.Themes {
		packages = append(packages, pkg)
	}
	return packages, nil
}
