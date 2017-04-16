package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/docker/docker/client"
	cli "github.com/jawher/mow.cli"

	"github.com/c0va23/duclean/containers"
	"github.com/c0va23/duclean/images"
	"github.com/c0va23/duclean/networks"
	"github.com/c0va23/duclean/volumes"
)

var (
	version   string
	buildTime string
	commit    string
)

func main() {
	app := cli.App("declean", "Docker universal cleaner")

	client, err := client.NewEnvClient()
	if nil != err {
		panic(err)
	}

	defer client.Close()

	app.Command("images", "Clean useless images", commandImages(client))

	app.Command("containers", "Clean containers", commandContainers(client))

	app.Command("volumes", "Clean useless volumes", commandVolumes(client))

	app.Command("networks", "Clean useless networks", commandNetworks(client))

	app.Command("version", "Print version", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			fmt.Printf("Version: %s\nBuild time: %s\nCommit: %s\n", version, buildTime, commit)
		}
	})

	app.Run(os.Args)
}

func getSafePeriod(cmd *cli.Cmd) *int {
	return cmd.IntOpt("safe-period", 0, "Save period (seconds)")

}

func getDryRun(cmd *cli.Cmd) *bool {
	return cmd.BoolOpt("dry-run", false, "Dry run")
}

func commandImages(client client.CommonAPIClient) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		safePeriod := getSafePeriod(cmd)
		dryRun := getDryRun(cmd)

		cmd.Action = func() {
			images.Clean(images.CleanOptions{
				DockerClient: client,
				DryRun:       *dryRun,
				SafePeriod:   time.Duration(*safePeriod),
			})
		}
	}
}

func commandContainers(client client.CommonAPIClient) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		safePeriod := getSafePeriod(cmd)
		dryRun := getDryRun(cmd)
		removeVolumes := cmd.BoolOpt("remove-volumes V", false, "Remove volumes")
		removeLinks := cmd.BoolOpt("remove-links L", false, "Remove links")

		cmd.Action = func() {
			containers.Clean(containers.CleanOptions{
				DockerClient:  client,
				DryRun:        *dryRun,
				SafePeriod:    time.Duration(*safePeriod),
				RemoveVolumes: *removeVolumes,
				RemoveLinks:   *removeLinks,
			})
		}
	}
}

func commandNetworks(client client.CommonAPIClient) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		dryRun := getDryRun(cmd)

		cmd.Action = func() {
			err := networks.Clean(networks.CleanOptions{
				DockerClient: client,
				DryRun:       *dryRun,
			})

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func commandVolumes(client client.CommonAPIClient) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		dryRun := getDryRun(cmd)
		force := cmd.BoolOpt("force F", false, "Force remove volumes")

		cmd.Action = func() {
			volumes.Clean(volumes.CleanOptions{
				DockerClient: client,
				DryRun:       *dryRun,
				Force:        *force,
			})
		}
	}
}

type sharedOptions struct {
	client client.CommonAPIClient
}
