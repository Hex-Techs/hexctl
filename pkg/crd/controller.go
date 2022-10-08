package crd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Hex-Techs/hexctl/pkg/common/display"
	"github.com/Hex-Techs/hexctl/pkg/common/exec"
	"github.com/Hex-Techs/hexctl/pkg/common/file"
	"github.com/Hex-Techs/hexctl/pkg/crd/templates"
)

const (
	// hack 目录
	boilerplate = "boilerplate"
	tools       = "tools.go"
	// build 文件
	makefile   = "Makefile"
	dockerfile = "Dockerfile"
	// update_codegen.sh 脚本
	script = "update-codegen.sh"

	gitignore    = ".gitignore"
	dockerignore = ".dockerignore"

	// 目录
	api           = "api"
	pkg           = "pkg"
	webhookDir    = "webhook"
	controllerDir = "controller"
)

// InitController init some files used by crd controller.
func InitController(gvk *GVK) error {
	// has no go.mod will raise error
	if !file.IsExists(goMOD) {
		return fmt.Errorf("%s was not found, may excute 'go mod init'", goMOD)
	}
	repo, err := getRepo()
	if err != nil {
		return err
	}
	gvk.Repo = repo
	if !file.IsExists(makefile) {
		file.Write(templates.ControllerMakefile, makefile)
	} else {
		return fmt.Errorf("%s was found in current directory, remove it before.", makefile)
	}
	if !file.IsExists(dockerfile) {
		file.Write(templates.Dockerfile, dockerfile)
	}
	if !file.IsExists(dockerignore) {
		file.Write(templates.DockerignoreTemplate, dockerignore)
	}
	if !file.IsExists(gitignore) {
		file.Write(templates.GitignoreTemplate, gitignore)
	}
	return nil
}

// 创建 controller
func CreateController(gvk *GVK) error {
	repo, err := getRepo()
	if err != nil {
		return err
	}
	gvk.Repo = repo
	d := display.NewTerminalDisplay("Are you sure to create controller", 0)
	ctrl := d.Confirm()
	if !ctrl {
		return nil
	}
	pkgWorkDir := filepath.Join(pkg, controllerDir)
	cwd := filepath.Join(pkgWorkDir, strings.ToLower(gvk.Kind))
	f := filepath.Join(cwd, fmt.Sprintf("%s_controller.go", strings.ToLower(gvk.Kind)))
	if file.IsExists(f) && !gvk.Force {
		// 如果 controler 已存在且没有强制覆盖，则退出
		return errors.New("the controller is already exist")
	}
	if !file.IsDir(pkg) {
		if err := os.Mkdir(pkg, 0755); err != nil {
			return err
		}
	}
	if !file.IsDir(pkgWorkDir) {
		if err := os.Mkdir(pkgWorkDir, 0755); err != nil {
			return err
		}
	}
	if !file.IsDir(cwd) {
		if err := os.Mkdir(cwd, 0755); err != nil {
			return err
		}
	}
	if err := file.WriteByTemp(f, templates.ControllerTemp, f, gvk); err != nil {
		return err
	}

	mainfile := "main.go"
	if !file.IsExists(mainfile) {
		return file.WriteByTemp(mainfile, templates.MainTemp, mainfile, gvk)
	}
	return exec.RunCommand("go mod tidy")
}

// CreateWebhook 创建 webhook
// 0.判断 webhook 目录，创建目录
// 1.创建 webhook 脚手架
func CreateWebhook(gvk *GVK) error {
	d := display.NewTerminalDisplay("Are you sure to create webhook", 0)
	if !d.Confirm() {
		return nil
	}
	hookDir := filepath.Join(pkg, webhookDir)
	cwd := filepath.Join(hookDir, strings.ToLower(gvk.Kind))
	f := filepath.Join(cwd, fmt.Sprintf("%s_webhook.go", strings.ToLower(gvk.Kind)))
	if file.IsExists(f) && !gvk.Force {
		return errors.New("the webhook is already exist")
	}
	if !file.IsDir(pkg) {
		if err := os.Mkdir(pkg, 0755); err != nil {
			return err
		}
	}
	if !file.IsDir(hookDir) {
		if err := os.Mkdir(hookDir, 0755); err != nil {
			return err
		}
	}
	if !file.IsDir(cwd) {
		if err := os.Mkdir(cwd, 0755); err != nil {
			return err
		}
	}
	return nil
}
