package crd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Hex-Techs/hexctl/pkg/common/display"
	"github.com/Hex-Techs/hexctl/pkg/common/exec"
	"github.com/Hex-Techs/hexctl/pkg/common/file"
	"github.com/Hex-Techs/hexctl/pkg/crd/templates"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	goMOD = "go.mod"
	// hack 目录
	hack        = "hack"
	boilerplate = "boilerplate"
	tools       = "tools.go"
	// project 管理文件
	project = "PROJECT"
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
	cmd           = "cmd"
	controllerDir = "controller"
	signalsDir    = "signals"
	signals       = "signals.go"
)

// Init 初始化项目
// 0.检测项目 go.mod 文件
// 1.创建 hack 目录
// 2.渲染 boilerplate 文件
// 3.渲染 tools 文件
// 4.渲染 PROJECT 文件
// 5.makefile 文件
// 6.dockerfile, ingore 文件
func Init(gvk *GVK) {
	// 0
	initProject(gvk)
	// 1
	if err := os.Mkdir(hack, 0755); err != nil {
		cobra.CheckErr(err)
	}
	// 2
	file.Write(templates.Boilerplate, filepath.Join(hack, boilerplate+".go.txt"))
	// 3
	file.Write(templates.ToolsTemp, filepath.Join(hack, tools))
	// 4
	if gvk.Repo == "" {
		gvk.Repo = getRepo()
	}
	pm := ProjectManager{
		Domain: gvk.Domain,
		Repo:   gvk.Repo,
	}
	c, _ := yaml.Marshal(pm)
	file.Write(string(c), project)
	// 5
	if !file.IsExists(makefile) {
		file.Write(templates.Makefile, makefile)
	}
	// 6
	if !file.IsExists(dockerfile) {
		file.Write(templates.Dockerfile, dockerfile)
	}
	if !file.IsExists(dockerignore) {
		file.Write(templates.DockerignoreTemplate, dockerignore)
	}
	if !file.IsExists(gitignore) {
		file.Write(templates.GitignoreTemplate, gitignore)
	}
	exec.RunCommand("go mod tidy")
}

// CreateAPI 创建 API
// 0.判断 update_codegen 脚本如果不存在则进行渲染
// 1.创建 api 存放目录
// 2.创建 doc info 文件
// 3.渲染 api types 文件
// 4.渲染 PROJECT 文件
func CreateAPI(gvk *GVK) {
	pm := getProject()
	// 0
	gvk.Repo = pm.Repo
	uc := filepath.Join(hack, script)
	if !file.IsExists(uc) {
		file.WriteByTemp(script, templates.UpdateCodeGenTemp, uc, gvk)
	}
	exec.RunCommand(fmt.Sprintf("chmod +x %s", filepath.Join(hack, script)))
	// 1
	apiWorkDir := filepath.Join(api, gvk.Version)
	if !file.IsDir(api) {
		if err := os.Mkdir(api, 0755); err != nil {
			cobra.CheckErr(err)
		}
	}
	if !file.IsDir(apiWorkDir) {
		if err := os.Mkdir(apiWorkDir, 0755); err != nil {
			cobra.CheckErr(err)
		}
	}
	// 2
	doc := filepath.Join(apiWorkDir, "doc.go")
	if !file.IsExists(doc) {
		file.WriteByTemp(doc, templates.Doc, doc, gvk)
	}
	gvi := filepath.Join(apiWorkDir, "groupversion_info.go")
	if !file.IsExists(gvi) {
		file.WriteByTemp(gvi, templates.Groupversion, gvi, gvk)
	}
	// 3
	apiResource := API{
		Kind:       gvk.Kind,
		Version:    gvk.Version,
		Namespaced: gvk.UseNamespace,
		Status:     gvk.UseStatus,
		Group:      gvk.Group,
		Domain:     pm.Domain,
	}
	typesName := filepath.Join(apiWorkDir, strings.ToLower(apiResource.Kind)+"_types.go")
	if file.IsExists(typesName) {
		exec.RunCommand("rm -f " + typesName)
	}
	file.WriteByTemp(typesName, templates.Types, typesName, apiResource)
	// 4
	if len(pm.API) == 0 {
		pm.Group = gvk.Group
	}
	add := true
	for i, v := range pm.API {
		if v.Kind == apiResource.Kind && v.Group == apiResource.Group &&
			v.Version == apiResource.Version && v.Domain == apiResource.Domain && !gvk.Force {
			cobra.CheckErr(fmt.Sprintf("%s kind already exist", v.Kind))
		}
		if v.Kind == apiResource.Kind && v.Group == apiResource.Group &&
			v.Version == apiResource.Version && v.Domain == apiResource.Domain && gvk.Force {
			add = false
			pm.API[i] = apiResource
		}
	}
	if add {
		pm.API = append(pm.API, apiResource)
	}
	c, _ := yaml.Marshal(pm)
	file.Write(string(c), project)
	color.Yellow.Println("Please modify and execute the hack/update-codegen.sh script to generate the corresponding code")
}

// CreateWebhook 创建 webhook
// 0.判断 webhook 目录，创建目录
// 1.创建 webhook 脚手架
func CreateWebhook(gvk *GVK) {
	pm := getProject()
	gvk.Repo = pm.Repo
	gvk.Domain = pm.Domain
	gvk.Group = pm.Group
	for i, v := range pm.API {
		if v.Group == gvk.Group && v.Kind == gvk.Kind {
			gvk.Version = v.Version
			gvk.UseNamespace = v.Namespaced
			gvk.UseStatus = v.Status
			pm.API[i].Webhook = true
		} else {
			// 如果在 PROJECT 配置中没找到 kind 则要求先创建 API 并进行相应的 update-codegen
			if i == (len(pm.API)-1) && gvk.Kind != v.Kind {
				cobra.CheckErr(fmt.Errorf("not found this kind %s in this project, you can use generate api to create a API first and generate code", gvk.Kind))
			}
		}
	}
	d := display.NewTerminalDisplay("Are you sure to create webhook", 0)
	if !d.Confirm() {
		return
	}
	hookDir := filepath.Join(pkg, webhookDir)
	cwd := filepath.Join(hookDir, strings.ToLower(gvk.Kind))
	f := filepath.Join(cwd, fmt.Sprintf("%s_webhook.go", strings.ToLower(gvk.Kind)))
	if file.IsExists(f) && !gvk.Force {
		cobra.CheckErr("the webhook is already exist")
	}
	if !file.IsDir(pkg) {
		if err := os.Mkdir(pkg, 0755); err != nil {
			cobra.CheckErr(err)
		}
	}
	if !file.IsDir(hookDir) {
		if err := os.Mkdir(hookDir, 0755); err != nil {
			cobra.CheckErr(err)
		}
	}
	if !file.IsDir(cwd) {
		if err := os.Mkdir(cwd, 0755); err != nil {
			cobra.CheckErr(err)
		}
	}
	generateSignalsFile()
}

// 创建 controller
// 0.创建目录
// 1.渲染 signals
// 2.渲染 controller 文件
// 3.创建 cmd 等结构和文件
// 4.创建 main.go 文件
// 5.配置 PROJECT 文件
func CreateController(gvk *GVK) {
	// 做一些初始化的工作
	pm := getProject()
	gvk.Repo = pm.Repo
	gvk.Domain = pm.Domain
	gvk.Group = pm.Group
	for i, v := range pm.API {
		if v.Group == gvk.Group && v.Kind == gvk.Kind {
			gvk.Version = v.Version
			gvk.UseNamespace = v.Namespaced
			gvk.UseStatus = v.Status
			pm.API[i].Controller = true
		} else {
			// 如果在 PROJECT 配置中没找到 kind 则要求先创建 API 并进行相应的 update-codegen
			if i == (len(pm.API)-1) && gvk.Kind != v.Kind {
				cobra.CheckErr(fmt.Errorf("not found this kind %s in this project, you can use generate api to create a API first and generate code", gvk.Kind))
			}
		}
	}
	d := display.NewTerminalDisplay("Are you sure to create controller", 0)
	ctrl := d.Confirm()
	if !ctrl {
		return
	}
	// 0
	pkgWorkDir := filepath.Join(pkg, controllerDir)
	cwd := filepath.Join(pkgWorkDir, strings.ToLower(gvk.Kind))
	f := filepath.Join(cwd, fmt.Sprintf("%s_controller.go", strings.ToLower(gvk.Kind)))
	if file.IsExists(f) && !gvk.Force {
		// 如果 controler 已存在且没有强制覆盖，则退出
		cobra.CheckErr("the controller is already exist")
	}
	if !file.IsDir(pkg) {
		if err := os.Mkdir(pkg, 0755); err != nil {
			cobra.CheckErr(err)
		}
	}
	if !file.IsDir(pkgWorkDir) {
		if err := os.Mkdir(pkgWorkDir, 0755); err != nil {
			cobra.CheckErr(err)
		}
	}
	if !file.IsDir(cwd) {
		if err := os.Mkdir(cwd, 0755); err != nil {
			cobra.CheckErr(err)
		}
	}
	// 1
	generateSignalsFile()
	file.WriteByTemp(f, templates.ControllerTemp, f, gvk)
	// 3
	createCmd(gvk)
	// 4
	mainfile := "main.go"
	if !file.IsExists(mainfile) {
		file.WriteByTemp(mainfile, templates.MainTemp, mainfile, gvk)
	}
	// 5
	c, _ := yaml.Marshal(pm)
	file.Write(string(c), project)
}

func createCmd(gvk *GVK) {
	cwd := filepath.Join(cmd, strings.ToLower(gvk.Kind))
	if !file.IsDir(cmd) {
		if err := os.Mkdir(cmd, 0755); err != nil {
			cobra.CheckErr(err)
		}
	}
	if !file.IsDir(cwd) {
		if err := os.Mkdir(cwd, 0755); err != nil {
			cobra.CheckErr(err)
		}
	}
	f := filepath.Join(cwd, strings.ToLower(gvk.Kind)+".go")
	file.WriteByTemp(f, templates.RunTemp, f, gvk)
}

// 初始化项目
func initProject(gvk *GVK) {
	if gvk.Repo != "" && !file.IsExists(goMOD) {
		cmd := "go mod init " + gvk.Repo
		exec.RunCommand(cmd)
	}
	if gvk.Repo != "" && file.IsExists(goMOD) {
		if !validateProject(gvk.Repo) {
			cobra.CheckErr("--repo is error, go.mod is exist but repo mismatch")
		}
	}
	if gvk.Repo == "" && !file.IsExists(goMOD) {
		cobra.CheckErr("go.mod is not exist, must given a repo")
	}
}

// 验证 project 是否已经使用 go mod 初始化
func validateProject(repo string) bool {
	content := file.Read(goMOD)
	s := fmt.Sprintf("^module %s\n", repo)
	matchStr := regexp.MustCompile(s)
	return matchStr.MatchString((content))
}

// 获取 go.mod 里的 repo 信息
func getRepo() string {
	fd, err := os.Open(goMOD)
	if err != nil {
		cobra.CheckErr(err)
	}
	defer fd.Close()
	br := bufio.NewReader(fd)
	a, _, err := br.ReadLine()
	if err != nil {
		cobra.CheckErr(err)
	}
	repos := strings.Split(string(a), " ")
	if len(repos) != 2 {
		cobra.CheckErr("go.mod read error")
	}
	return repos[1]
}

// 获取项目配置文件
func getProject() *ProjectManager {
	if !file.IsExists(project) {
		cobra.CheckErr("PROJECT is not exist")
	}
	c := file.Read(project)
	var pm ProjectManager
	err := yaml.Unmarshal([]byte(c), &pm)
	if err != nil {
		cobra.CheckErr(err)
	}
	return &pm
}

// 生成信号处理程序
func generateSignalsFile() {
	sigDir := filepath.Join(pkg, signalsDir)
	if !file.IsDir(pkg) {
		if err := os.Mkdir(pkg, 0755); err != nil {
			cobra.CheckErr(err)
		}
	}
	if !file.IsDir(sigDir) {
		if err := os.Mkdir(sigDir, 0755); err != nil {
			cobra.CheckErr(err)
		}
	}
	f := filepath.Join(sigDir, signals)
	if !file.IsExists(f) {
		file.Write(templates.SignalTemp, f)
	}
}
