package npm

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"path/filepath"
	"spdx-sbom-generator/internal/helper"
	"spdx-sbom-generator/internal/models"
	"spdx-sbom-generator/internal/reader"
	"strings"
	"testing"
)

func TestNPM(t *testing.T) {
	t.Run("test is valid", TestIsValid)
	t.Run("test has modules installed", TestHasModulesInstalled)
	t.Run("test get module", TestGetModule)
	t.Run("test list modules", TestListModules)
	t.Run("test list all modules", TestListAllModules)
}

func TestIsValid(t *testing.T) {
	n := New()
	path := fmt.Sprintf("%s/test", getPath())

	valid := n.IsValid(path)
	invalid := n.IsValid(getPath())

	// Assert
	assert.Equal(t, true, valid)
	assert.Equal(t, false, invalid)
}

func TestHasModulesInstalled(t *testing.T) {
	n := New()
	path := fmt.Sprintf("%s/test", getPath())

	installed := n.HasModulesInstalled(path)
	assert.NoError(t, installed)
	uninstalled := n.HasModulesInstalled(getPath())
	assert.Error(t, uninstalled)
}

func TestGetModule(t *testing.T) {
	n := New()
	path := fmt.Sprintf("%s/test", getPath())
	mod, err := n.GetRootModule(path)

	assert.NoError(t, err)
	assert.Equal(t, "e-commerce", mod.Name)
	assert.Equal(t, "ahmed saber", mod.Supplier.Name)
	assert.Equal(t, "1.0.0", mod.Version)

}

func TestListModules(t *testing.T) {
	n := New()
	path := fmt.Sprintf("%s/test", getPath())
	mods, err := n.ListUsedModules(path)

	assert.NoError(t, err)

	count := 0
	for _, mod := range mods {

		if mod.Name == "bcryptjs" {
			assert.Equal(t, "bcryptjs", mod.Name)
			assert.Equal(t, "2.4.3", mod.Version)
			count++
			continue
		}

		if mod.Name == "body-parser" {
			assert.Equal(t, "body-parser", mod.Name)
			assert.Equal(t, "1.18.3", mod.Version)
			count++
			continue
		}
		if mod.Name == "shortid" {
			assert.Equal(t, "shortid", mod.Name)
			assert.Equal(t, "2.2.13", mod.Version)
			count++
			continue
		}

		if mod.Name == "validator" {
			assert.Equal(t, "validator", mod.Name)
			assert.Equal(t, "10.7.1", mod.Version)
			count++
			continue
		}
	}

	assert.Equal(t, 4, count)
}

func TestListAllModules(t *testing.T) {
	n := New()
	path := fmt.Sprintf("%s/test", getPath())
	mods, err := n.ListModulesWithDeps(path)

	assert.NoError(t, err)

	count := 0
	for _, mod := range mods {
		if mod.Name == "validator-10.11.0" {
			assert.Equal(t, "10.11.0", mod.Version)
			assert.Equal(t, "https://registry.npmjs.org/validator/-/validator-10.11.0.tgz", mod.PackageURL)
			assert.Equal(t, models.HashAlgorithm("sha512"), mod.CheckSum.Algorithm)
			assert.Equal(t, "X/p3UZerAIsbBfN/IwahhYaBbY68EN/UQBWHtsbXGT5bfrH/p4NQzUCG1kF/rtKaNpnJ7jAu6NGTdSNtyNIXMw==", mod.CheckSum.Value)
			assert.Equal(t, "Copyright (c) 2018 Chris O'Hara <cohara87@gmail.com>", mod.Copyright)
			assert.Equal(t, "MIT", mod.LicenseDeclared)
			count++
			continue
		}
		if mod.Name == "shortid-2.2.16" {
			assert.Equal(t, "2.2.16", mod.Version)
			assert.Equal(t, "https://registry.npmjs.org/shortid/-/shortid-2.2.16.tgz", mod.PackageURL)
			assert.Equal(t, models.HashAlgorithm("sha512"), mod.CheckSum.Algorithm)
			assert.Equal(t, "Ugt+GIZqvGXCIItnsL+lvFJOiN7RYqlGy7QE41O3YC1xbNSeDGIRO7xg2JJXIAj1cAGnOeC1r7/T9pgrtQbv4g==", mod.CheckSum.Value)
			assert.Equal(t, "Copyright (c) Dylan Greene", mod.Copyright)
			assert.Equal(t, "MITNFA", mod.LicenseDeclared)
			count++
			continue
		}
		if mod.Name == "body-parser-1.19.0" {
			assert.Equal(t, "1.19.0", mod.Version)
			assert.Equal(t, "https://registry.npmjs.org/body-parser/-/body-parser-1.19.0.tgz", mod.PackageURL)
			assert.Equal(t, models.HashAlgorithm("sha512"), mod.CheckSum.Algorithm)
			assert.Equal(t, "dhEPs72UPbDnAQJ9ZKMNTP6ptJaionhP5cBb541nXPlW60Jepo9RV/a4fX4XWW9CuFNK22krhrj1+rgzifNCsw==", mod.CheckSum.Value)
			assert.Equal(t, "Copyright (c) 2014 Jonathan Ong <me@jongleberry.com>", mod.Copyright)
			assert.Equal(t, "MIT", mod.LicenseDeclared)
			count++
			continue
		}
		if mod.Name == "bcryptjs-2.4.3" {
			fmt.Println("bcrypt: ", mod)
			assert.Equal(t, "2.4.3", mod.Version)
			assert.Equal(t, "https://registry.npmjs.org/bcryptjs/-/bcryptjs-2.4.3.tgz", mod.PackageURL)
			assert.Equal(t, models.HashAlgorithm("sha1"), mod.CheckSum.Algorithm)
			assert.Equal(t, "mrVie5PmBiH/fNrF2pczAn3x0Ms=", mod.CheckSum.Value)
			assert.Equal(t, "Copyright (c) 2012 Nevins Bartolomeo <nevins.bartolomeo@gmail.com>", mod.Copyright)
			assert.Equal(t, "MIT", mod.LicenseDeclared)
			count++
			continue
		}
	}

	assert.Equal(t, 4, count)
}


func TestGetCopyright(t *testing.T) {
	path := fmt.Sprintf("%s/test", getPath())
	licensePath := filepath.Join(path, "node_modules", "bcryptjs", "LICENSE")
	if helper.Exists(licensePath) {
		r := reader.New(licensePath)
		s := r.StringFromFile()
		res := helper.GetCopyright(s)
		assert.Equal(t, "Copyright (c) 2012 Nevins Bartolomeo <nevins.bartolomeo@gmail.com>", res)
	}

	licensePath2 := filepath.Join(path, "node_modules", "shortid", "LICENSE")
	if helper.Exists(licensePath2) {
		r := reader.New(licensePath2)
		s := r.StringFromFile()
		res := helper.GetCopyright(s)
		assert.Equal(t, "Copyright (c) Dylan Greene", res)
	}
}

func getPath() string {
	cmd := exec.Command("pwd")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	path := strings.TrimSuffix(string(output), "\n")

	return path
}