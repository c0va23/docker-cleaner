package main

import (
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/client"
	cli "github.com/jawher/mow.cli"
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
			cleanImages(cleanImagesOptions{
				sharedOptions: sharedOptions{
					client: client,
				},
				dryRun:     *dryRun,
				safePeriod: time.Duration(*safePeriod),
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
			cleanContainers(cleanContainersOptions{
				sharedOptions: sharedOptions{
					client: client,
				},
				dryRun:        *dryRun,
				safePeriod:    time.Duration(*safePeriod),
				removeVolumes: *removeVolumes,
				removeLinks:   *removeLinks,
			})
		}
	}
}

func commandVolumes(client client.CommonAPIClient) func(*cli.Cmd) {
	return func(cmd *cli.Cmd) {
		dryRun := getDryRun(cmd)
		force := cmd.BoolOpt("force F", false, "Force remove volumes")

		cmd.Action = func() {
			cleanVolumes(cleanVolumesOptions{
				sharedOptions: sharedOptions{
					client: client,
				},
				dryRun: *dryRun,
				force:  *force,
			})
		}
	}
}

type sharedOptions struct {
	client client.CommonAPIClient
}
