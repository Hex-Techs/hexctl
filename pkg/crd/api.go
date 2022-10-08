package crd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Hex-Techs/hexctl/pkg/common/exec"
	"github.com/Hex-Techs/hexctl/pkg/common/file"
	"github.com/Hex-Techs/hexctl/pkg/crd/templates"
)

// InitAPI init some files used by crd API.
func InitAPI(gvk *GVK) error {
	// has no go.mod will raise error
	if !file.IsExists(goMOD) {
		return fmt.Errorf("%s was not found, may excute 'go mod init'", goMOD)
	}
	if err := os.Mkdir(hack, 0755); err != nil {
		return err
	}
	repo, err := getRepo()
	if err != nil {
		return err
	}
	gvk.Repo = repo
	file.Write(templates.Boilerplate, filepath.Join(hack, boilerplate+".go.txt"))
	file.Write(templates.ToolsTemp, filepath.Join(hack, tools))
	uc := filepath.Join(hack, script)
	if !file.IsExists(uc) {
		if err := file.WriteByTemp(script, templates.UpdateCodeGenTemp, uc, gvk); err != nil {
			return err
		}
	}
	if err := exec.RunCommand(fmt.Sprintf("chmod +x %s", filepath.Join(hack, script))); err != nil {
		return err
	}
	// NOTE: api makefile and controller makefile use diffrenet templates.
	if !file.IsExists(makefile) {
		file.Write(templates.APIMakefile, makefile)
	}
	if !file.IsExists(gitignore) {
		file.Write(templates.GitignoreTemplate, gitignore)
	}
	return exec.RunCommand("go mod tidy")
}

// CreateAPI 创建 API
func CreateAPI(gvk *GVK) error {
	repo, err := getRepo()
	if err != nil {
		return err
	}
	gvk.Repo = repo
	apiWorkDir := filepath.Join(api, gvk.Version)
	if !file.IsDir(api) {
		if err := os.Mkdir(api, 0755); err != nil {
			return err
		}
	}
	if !file.IsDir(apiWorkDir) {
		if err := os.Mkdir(apiWorkDir, 0755); err != nil {
			return err
		}
	}
	gvi := filepath.Join(apiWorkDir, "groupversion_info.go")
	if !file.IsExists(gvi) {
		if err := file.WriteByTemp(gvi, templates.Groupversion, gvi, gvk); err != nil {
			return err
		}
	}
	apiResource := API{
		Kind:       gvk.Kind,
		Version:    gvk.Version,
		Namespaced: gvk.UseNamespace,
		Status:     gvk.UseStatus,
		Group:      gvk.Group,
		Domain:     gvk.Repo,
	}
	typesName := filepath.Join(apiWorkDir, strings.ToLower(apiResource.Kind)+"_types.go")
	if file.IsExists(typesName) {
		if err := exec.RunCommand("rm -f " + typesName); err != nil {
			return err
		}
	}
	if err := file.WriteByTemp(typesName, templates.Types, typesName, apiResource); err != nil {
		return err
	}
	fmt.Printf("Generate success.\nPlease modify and execute the hack/update-codegen.sh script to generate the corresponding code if you need.\n")
	return exec.RunCommand("go mod tidy")
}
