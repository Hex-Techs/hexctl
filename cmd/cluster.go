/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"strings"

	"github.com/Fize/n/pkg/cluster"
	"github.com/Fize/n/pkg/output"
	"github.com/Fize/n/pkg/utils"
	"github.com/spf13/cobra"
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "deploy a new or a component for kubernetes cluster",
	Long: `cluster is a command for n that manage kubernetes cluster.
quickly create a new kubernetes cluster.
quickly deploy a component for exist kubernetes cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		if clusterCommand.Master != "" {
			if err := utils.ValidataIP(clusterCommand.Master); err != nil {
				output.Fatalln(err)
			}
			if err := cluster.PreCheck(clusterCommand.PodCIDR, clusterCommand.ServiceCIDR,
				clusterCommand.ServicePortRange, clusterCommand.Master); err != nil {
				output.Fatalln(err)
			}
			// init cluster
			// cluster.Run(password, key, "22", repo, volume, podcidr, servicecidr, apiserver, token, ca, iface, master...)
			return
		}
		if len(clusterCommand.Node) != 0 {
			if clusterCommand.APIServer == "" {
				output.Fatalln("apiserver can not empty when join a cluster.")
			}
			if clusterCommand.Token == "" {
				output.Fatalln("token can not empty when join a cluster.")
			}
			if clusterCommand.CAHash == "" {
				if !clusterCommand.UnSafe {
					output.Fatalln("ca can not empty when join a cluster use safe, unless --unsafe")
				}
			}
			if err := utils.ValidataIP(clusterCommand.Node...); err != nil {
				output.Fatalln(err)
			}
			if err := utils.ValidataURL(clusterCommand.APIServer); err != nil {
				output.Fatalln(err)
			}
			// join cluster
		}
		if validateArgsLocal(args) {
			output.Noteln("install local cluster")
			cluster.Run(&clusterCommand)
		} else {
			output.Noteln("install remote cluster")
		}
		// cmd.Help()
	},
}

var clusterCommand cluster.ClusterCommand

func init() {
	clusterCmd.Flags().StringVarP(&clusterCommand.Master, "master", "m", "", "kubernetes master address.")
	clusterCmd.Flags().StringSliceVarP(&clusterCommand.Node, "node", "n", []string{}, "kubernetes node address.")
	clusterCmd.Flags().StringSliceVarP(&clusterCommand.CertHost, "cert-host", "", []string{}, "these host will joins the certificate.")
	clusterCmd.Flags().StringVarP(&clusterCommand.Repo, "repo", "", "", "image repo address.")
	clusterCmd.Flags().StringVarP(&clusterCommand.Volume, "volume", "v", "", "the disk character used to format a disk. if not, do not format")
	clusterCmd.Flags().StringVarP(&clusterCommand.APIServer, "apiserver", "", "", "cluster apiserver url.")
	clusterCmd.Flags().StringVarP(&clusterCommand.PodCIDR, "pod-network-cidr", "", "", "Specify range of IP addresses for the pod network. If set, the control plane will automatically allocate CIDRs for every node.")
	clusterCmd.Flags().StringVarP(&clusterCommand.ServiceCIDR, "service-cidr", "", "", "Use alternative range of IP address for service VIPs. (default 10.96.0.0/12).")
	clusterCmd.Flags().StringVarP(&clusterCommand.ServicePortRange, "service-port-range", "", "", "Specifies the range of node ports that the service can use.")
	clusterCmd.Flags().StringVarP(&clusterCommand.Token, "token", "t", "", "cluster token.")
	clusterCmd.Flags().StringVarP(&clusterCommand.CAHash, "ca-hash", "", "", "ca public key hash.")
	clusterCmd.Flags().BoolVarP(&clusterCommand.UnSafe, "unsafe", "", false, "allow join cluster without --ca.")
	clusterCmd.Flags().BoolVarP(&clusterCommand.IPVS, "ipvs", "", false, "set kube-proxy mode is ipvs.")
	clusterCmd.Flags().BoolVarP(&clusterCommand.ControlPlane, "control-plane", "", false, "set node as master join the exists cluster.")
	clusterCmd.Flags().StringSliceVarP(&clusterCommand.Ignore, "ignore", "i", []string{}, "A list of checks whose errors will be shown as warnings.")
	clusterCmd.Flags().StringVarP(&clusterCommand.Password, "password", "p", "", "password for host.")
	clusterCmd.Flags().StringVarP(&clusterCommand.Key, "key", "k", "", "private key for host.")
	clusterCmd.Flags().StringVarP(&clusterCommand.Iface, "iface", "", "eth0", "network device, default: eth0.")

	rootCmd.AddCommand(clusterCmd)
}

func validateArgsLocal(args []string) bool {
	if len(args) > 1 {
		output.Fatalf("Unknown args %v, only one arg is accepted.\n", args)
	}
	if len(args) == 0 {
		return false
	}
	if strings.ToLower(args[0]) != "local" {
		output.Fatalf("Unknown args %v, only one arg is accepted.\n", args)
	}
	return true
}
