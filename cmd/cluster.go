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
	"github.com/Hex-Techs/hexctl/pkg/cluster"
	"github.com/Hex-Techs/hexctl/pkg/output"
	"github.com/spf13/cobra"
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "manage kubernetes cluster",
	Long: `cluster is a command for n that manage kubernetes cluster.
quickly create or destroy a kubernetes cluster.
quickly deploy a component for exist kubernetes cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "deploy a new or a component for kubernetes cluster",
	Long:  `cluster new is command for create a new kubernetes cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			if len(args) != 1 {
				output.Fatalln("too many args get, must given one local or remote")
			}
		} else {
			// default local
			args = []string{"local"}
		}
		if args[0] == "local" {
			kc := cluster.NewKubernetesCluster(&clusterCommand, "local", "n1", "init")
			cluster.StartKubernetesCluster(kc)
			return
		} else if args[0] == "remote" {
			kc := cluster.NewKubernetesCluster(&clusterCommand, "remote", clusterCommand.IP, "init")
			cluster.StartKubernetesCluster(kc)
			return
		} else {
			output.Errorf("unknown args %s, must given local or remote", args[0])
		}
		cmd.Help()
	},
}

var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "join a kubernetes cluster",
	Long:  `cluster join is command for join a exists kubernetes cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			if len(args) != 1 {
				output.Fatalln("too many args get, must given one local or remote")
			}
		} else {
			// default local
			args = []string{"local"}
		}

		if args[0] == "local" {
			kc := cluster.NewKubernetesCluster(&clusterCommand, "local", "n2", "join")
			cluster.StartKubernetesCluster(kc)
			return
		} else if args[0] == "remote" {
			kc := cluster.NewKubernetesCluster(&clusterCommand, "remote", clusterCommand.IP, "join")
			cluster.StartKubernetesCluster(kc)
			return
		} else {
			output.Errorf("unknown args %s, must given local or remote", args[0])
		}
		cmd.Help()
	},
}

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "destroy local environment",
	Run: func(cmd *cobra.Command, args []string) {
		cluster.DestroyVirtualMachine()
	},
}

var clusterCommand cluster.ClusterCommand

func init() {
	newCmd.Flags().StringVarP(&clusterCommand.PodCIDR, "pod-network-cidr", "", "10.244.0.0/16", "Specify range of IP addresses for the pod network. If set, the control plane will automatically allocate CIDRs for every node. (default 10.244.0.0/16).")
	newCmd.Flags().StringVarP(&clusterCommand.ServiceCIDR, "service-cidr", "", "10.96.0.0/12", "Use alternative range of IP address for service VIPs. (default 10.96.0.0/12).")
	newCmd.Flags().StringVarP(&clusterCommand.ServicePortRange, "service-port-range", "", "", "Specifies the range of node ports that the service can use.")
	newCmd.Flags().BoolVarP(&clusterCommand.CN, "cn", "", true, "set whether it is a cluster in gfw.")
	newCmd.Flags().StringSliceVarP(&clusterCommand.Ignore, "ignore", "i", []string{}, "A list of checks whose errors will be shown as warnings.")
	newCmd.Flags().StringVarP(&clusterCommand.User, "user", "u", "", "ssh user, default, root")
	// newCmd.Flags().StringVarP(&clusterCommand.Iface, "iface", "", "eth0", "network device, default: eth0.")
	newCmd.Flags().StringVarP(&clusterCommand.Repo, "image-repository", "", "", "image repository mirror. e.g. registry.cn-hangzhou.aliyuncs.com/google_containers")
	newCmd.Flags().StringSliceVarP(&clusterCommand.CertSANs, "cert-sans", "", []string{}, "Optional extra Subject Alternative Names (SANs) to use for the API Server serving certificate. Can be both IP addresses and DNS names.")

	joinCmd.Flags().StringVarP(&clusterCommand.User, "user", "u", "", "ssh user, default, root")
	joinCmd.Flags().StringVarP(&clusterCommand.Token, "token", "t", "", "cluster token.")
	joinCmd.Flags().StringVarP(&clusterCommand.CAHash, "ca-hash", "", "", "ca public key hash, already use '--discovery-token-unsafe-skip-ca-verification' option, you can skip this option")
	// joinCmd.Flags().BoolVarP(&clusterCommand.UnSafe, "unsafe", "", true, "allow join cluster without --ca.")
	joinCmd.Flags().StringVarP(&clusterCommand.Endpoint, "apiserver", "", "", "The IP address the API Server.")

	newCmd.Flags().StringVarP(&clusterCommand.Password, "password", "p", "", "password for host.")
	newCmd.Flags().StringVarP(&clusterCommand.Key, "key", "k", "", "private key for host.")
	newCmd.Flags().StringVarP(&clusterCommand.IP, "ip", "", "", "specify the ip address of the host, which is not required if it is in local mode.")
	joinCmd.Flags().StringVarP(&clusterCommand.Password, "password", "p", "", "password for host.")
	joinCmd.Flags().StringVarP(&clusterCommand.Key, "key", "k", "", "private key for host.")
	joinCmd.Flags().StringVarP(&clusterCommand.IP, "ip", "", "", "specify the ip address of the host, which is not required if it is in local mode.")

	clusterCmd.AddCommand(newCmd)
	clusterCmd.AddCommand(joinCmd)
	clusterCmd.AddCommand(destroyCmd)
	rootCmd.AddCommand(clusterCmd)
}
