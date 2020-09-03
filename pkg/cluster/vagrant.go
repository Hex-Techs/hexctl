package cluster

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/Fize/n/pkg/output"
)

const workDir = "~/.n"

type vagrantFileFiled struct {
	Num    []int
	Cpu    int
	Memory int
}

func renderVagrantFile() {
	tpl := template.Must(template.New("Vagrantfile").Parse(VagrantFileTmpl))
	f, err := os.Create(filepath.Join(workDir, "Vagrantfile"))
	if err != nil {
		output.Fatalln("cant create vagrant file:", err)
	}
	// now, cpu and memory are hardcoded by 2 and 2048, and 2 virtulabox virtual machine in use
	tpl.Execute(f, vagrantFileFiled{Num: []int{1, 2}, Cpu: 2, Memory: 2048})
}

func createWorkDir() {}
