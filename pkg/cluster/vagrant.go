package cluster

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"text/template"

	"github.com/Fize/n/pkg/output"
	"github.com/Fize/n/pkg/utils"
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
		output.Fatalf("%s is alrady exist, please clean it.", homeDir())
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
