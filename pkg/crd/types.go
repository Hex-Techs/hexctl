// Package crd help to create a operator controller by kubeilder or code-generate
package crd

import "strings"

// GVK this is the infomation of your CRD
type GVK struct {
	Repo         string
	Domain       string
	Group        string
	Version      string
	Kind         string
	UseNamespace bool
	UseStatus    bool
	// attempt to create resource even if it already exists
	Force bool
}

func (g *GVK) ToTitle(s string) string {
	return strings.Title(s)
}

func (g *GVK) ToLower(s string) string {
	return strings.ToLower(s)
}

// ProjectManager 项目结构
type ProjectManager struct {
	Domain string `yaml:"domain,omitempty"`
	Repo   string `yaml:"repo,omitempty"`
	Group  string `yaml:"group,omitempty"`
	API    []API  `yaml:"api,omitempty"`
}

// API api 相关的配置信息
type API struct {
	Kind       string `yaml:"kind,omitempty"`
	Version    string `yaml:"version,omitempty"`
	Namespaced bool   `yaml:"namespaced,omitempty"`
	Status     bool   `yaml:"status,omitempty"`
	Group      string `yaml:"group,omitempty"`
	Domain     string `yaml:"domain,omitempty"`
	Controller bool   `yaml:"controller,omitempty"`
	Webhook    bool   `yaml:"webhook,omitempty"`
}
