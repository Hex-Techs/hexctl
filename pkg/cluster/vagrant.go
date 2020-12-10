package cluster

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Hex-Techs/hexctl/pkg/output"
	"github.com/Hex-Techs/hexctl/pkg/utils"
)

const workDir = ".n"

func homeDir() string {
	u, _ := user.Current()
	return filepath.Join(u.HomeDir, workDir)
}

type vagrantFileFiled struct {
	Num    []int
	Cpu    int
	Memory int
	Pub    string
}

func renderVagrantFile(ncmd *ClusterCommand) {
	tpl := template.Must(template.New("Vagrantfile").Parse(VagrantFileTmpl))
	f, err := os.Create(filepath.Join(homeDir(), "Vagrantfile"))
	if err != nil {
		output.Fatalln("cant create vagrant file:", err)
	}
	// cpu and memory are hardcoded by 2 and 2048, and 2 virtulabox virtual machine in use
	tpl.Execute(f, vagrantFileFiled{Num: []int{1, 2}, Cpu: 2, Memory: 2048, Pub: ncmd.Key})
}

func handleWorkDir() {
	if utils.IsExists(homeDir()) {
		output.Errorf("%s is alrady exist\n", homeDir())
		output.Notef("Do you want continue? (yes/no) ")
		reader := bufio.NewReader(os.Stdin)
		cmdStr, err := reader.ReadString('\n')
		if err != nil {
			output.Fatalln(err)
		}
		choose := strings.ToLower(cmdStr)
		if choose == "yes\n" || choose == "y\n" || choose == "\n" {
			output.Progressln("Continue Progress...")
			return
		} else if choose == "no\n" || choose == "n\n" {
			os.Exit(0)
		} else {
			output.Fatalln("unknow option")
		}
	}
	if err := os.Mkdir(homeDir(), os.ModePerm); err != nil {
		output.Fatalln(err)
	}
}

func startVirtualMachine() {
	if err := os.Chdir(homeDir()); err != nil {
		output.Fatalln(err)
	}
	utils.RunCommand("vagrant up")
}

func DestroyVirtualMachine() {
	pwd, _ := os.Getwd()
	if err := os.Chdir(homeDir()); err != nil {
		output.Fatalln(err)
	}
	utils.RunCommand("vagrant destroy")
	os.Chdir(pwd)
	utils.RunCommand(fmt.Sprintf("rm -rf %s", homeDir()))
}
