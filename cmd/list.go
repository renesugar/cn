package cmd

import (
	"fmt"
	"log"
	"regexp"

	"github.com/apcera/termtables"
	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

// CliClusterList is the Cobra CLI call
func CliClusterList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "Print the list of Ceph cluster(s)",
		Args:  cobra.NoArgs,
		Run:   listNano,
	}
	return cmd
}

// listNano prints running Ceph cluster(s)
func listNano(cmd *cobra.Command, args []string) {
	showNanoClusters()
}

func showNanoClusters() {
	listOptions := types.ContainerListOptions{
		All:   true,
		Quiet: true,
	}
	containers, err := getDocker().ContainerList(ctx, listOptions)
	if err != nil {
		log.Fatal(err)
	}

	table := termtables.CreateTable()
	table.AddHeaders("NAME", "STATUS", "IMAGE", "IMAGE RELEASE", "IMAGE CREATION TIME")

	// run the loop on both indexes, it's fine they have the same length
	for _, container := range containers {
		for i := range container.Names {
			match, _ := regexp.MatchString(ContainerNamePrefix, container.Names[i])
			if match {
				// remove 7 first char since container.ImageID is in the form of sha256:<ID>
				containerImgTag := inspectImage(container.ImageID[7:], "tag")
				containerImgCreated := inspectImage(container.ImageID[7:], "created")
				containerImgRelease := inspectImage(container.ImageID[7:], "release")
				ContainerNameToShow := container.Names[i][len(ContainerNamePrefix):]
				// We trim again so we can remove the '/' since container name returned is /ceph-nano
				table.AddRow(ContainerNameToShow[1:], container.State, containerImgTag, containerImgRelease, containerImgCreated)
			}
		}
	}
	fmt.Println(table.Render())
}
