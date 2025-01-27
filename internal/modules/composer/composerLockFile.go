// SPDX-License-Identifier: Apache-2.0

package composer

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"strings"

	"spdx-sbom-generator/internal/helper"
	"spdx-sbom-generator/internal/models"
)

type ComposerLockFile struct {
	Packages    []ComposerLockPackage
	PackagesDev []ComposerLockPackage `json:"packages-dev"`
}

type ComposerLockPackage struct {
	Name        string
	Version     string
	Type        string
	Dist        ComposerLockPackageDist
	License     []string
	Description string
	Source      ComposerLockPackageSource
	Authors     []ComposerLockPackageAuthor
}
type ComposerLockPackageAuthor struct {
	Name  string
	Email string
}

type ComposerLockPackageSource struct {
	Type      string
	URL       string
	Reference string
}

type ComposerLockPackageDist struct {
	Type      string
	URL       string
	Reference string
	Shasum    string
}

func getComposerLockFileData() (ComposerLockFile, error) {

	raw, err := ioutil.ReadFile(COMPOSER_LOCK_FILE_NAME)
	if err != nil {
		return ComposerLockFile{}, err
	}

	var fileData ComposerLockFile
	err = json.Unmarshal(raw, &fileData)
	if err != nil {
		return ComposerLockFile{}, err
	}
	return fileData, nil
}

func (m *composer) getModulesFromComposerLockFile() ([]models.Module, error) {

	modules := make([]models.Module, 0)

	info, err := getComposerLockFileData()
	if err != nil {
		return nil, err
	}

	mainMod, err := m.getProjectInfo()
	if err != nil {
		return nil, err
	}

	modules = append(modules, mainMod)

	if len(info.Packages) > 0 {
		for _, pckg := range info.Packages {
			mod := convertLockPackageToModule(pckg)
			modules = append(modules, mod)
		}
	}

	if len(info.PackagesDev) > 0 {
		for _, pckg := range info.PackagesDev {
			mod := convertLockPackageToModule(pckg)
			modules = append(modules, mod)
		}
	}

	return modules, nil
}

func convertLockPackageToModule(dep ComposerLockPackage) models.Module {

	module := models.Module{
		Version:    normalizePackageVersion(dep.Version),
		Name:       getName(dep.Name),
		Root:       false,
		PackageURL: genUrlFromComposerPackage(dep),
		CheckSum: &models.CheckSum{
			Algorithm: models.HashAlgoSHA1,
			Value:     getCheckSumValue(dep),
		},
		Supplier:  getAuthor(dep),
		LocalPath: getLocalPath(dep),
		Modules:   map[string]*models.Module{},
	}
	path := getLocalPath(dep)
	licensePkg, err := helper.GetLicenses(path)
	if err == nil {
		module.LicenseDeclared = helper.BuildLicenseDeclared(licensePkg.ID)
		module.LicenseConcluded = helper.BuildLicenseConcluded(licensePkg.ID)
		module.Copyright = helper.GetCopyright(licensePkg.ExtractedText)
		module.CommentsLicense = licensePkg.Comments
	}

	return module
}

func getAuthor(dep ComposerLockPackage) models.SupplierContact {

	authors := dep.Authors
	if len(authors) == 0 {
		return models.SupplierContact{}
	}
	author := authors[0]
	pckAuthor := models.SupplierContact{
		Name:  author.Name,
		Email: author.Email,
	}
	return pckAuthor
}

func getName(moduleName string) string {
	s := strings.Split(moduleName, "/")

	if len(s) > 1 {
		return s[1]
	}

	return s[0]
}

func genUrlFromComposerPackage(dep ComposerLockPackage) string {
	URL := dep.Dist.URL
	if URL != "" {
		return URL
	}

	return genComposerUrl(dep.Name, dep.Version)
}
func genComposerUrl(name string, version string) string {
	return "https://github.com/" + name + ".git"
}

func normalizePackageVersion(version string) string {
	parts := strings.Split(version, "v")

	if parts[0] != "" {
		return version
	}

	if len(parts) > 1 {
		return parts[1]
	}

	return parts[0]
}

func getCheckSumValue(module ComposerLockPackage) string {
	value := module.Dist.Shasum
	if value != "" {
		return value
	}

	return readCheckSum(genUrlFromComposerPackage(module))
}

func readCheckSum(content string) string {
	if content == "" {
		return ""
	}
	h := sha1.New()
	h.Write([]byte(content))
	return hex.EncodeToString(h.Sum(nil))
}

func getLocalPath(module ComposerLockPackage) string {
	path := "./vendor/" + module.Name
	return path
}
