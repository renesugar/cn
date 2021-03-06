package cmd

import (
	"log"

	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

// CliClusterStatus is the Cobra CLI call
func CliClusterStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Stat object storage server",
		Args:  cobra.ExactArgs(1),
		Run:   statusNano,
	}

	return cmd
}

// statusNano shows Ceph Nano status
func statusNano(cmd *cobra.Command, args []string) {
	ContainerName := ContainerNamePrefix + args[0]

	notExistCheck(ContainerName)
	notRunningCheck(ContainerName)
	echoInfo(ContainerName)
}

// containerStatus checks container status
// the parameter corresponds to the type listOptions and its entry all
func containerStatus(ContainerName string, allList bool, containerState string) bool {
	listOptions := types.ContainerListOptions{
		All:   allList,
		Quiet: true,
	}
	containers, err := getDocker().ContainerList(ctx, listOptions)
	if err != nil {
		log.Fatal(err)
	}

	// run the loop on both indexes, it's fine they have the same length
	for _, container := range containers {
		for i := range container.Names {
			if container.Names[i] == "/"+ContainerName && container.State == containerState {
				return true
			}
		}
	}
	return false
}
