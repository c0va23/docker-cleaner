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
	safePeriod := app.IntOpt("safe-period", 0, "Save period (seconds)")
	dryRun := app.BoolOpt("dry-run", false, "Dry run")

	client, err := client.NewEnvClient()
	if nil != err {
		panic(err)
	}

	getSharedOptions := func() sharedOptions {
		return sharedOptions{
			client:     client,
			safePeriod: time.Duration(*safePeriod),
			dryRun:     *dryRun,
		}
	}

	app.Command("images", "Clean useless images", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			cleanImages(cleanImagesOptions{
				sharedOptions: getSharedOptions(),
			})
		}
	})

	app.Command("containers", "Clean containers", func(cmd *cli.Cmd) {
		removeVolumes := cmd.BoolOpt("remove-volumes V", false, "Remove volumes")
		removeLinks := cmd.BoolOpt("remove-links L", false, "Remove links")

		cmd.Action = func() {
			cleanContainers(cleanContainersOptions{
				sharedOptions: getSharedOptions(),
				removeVolumes: *removeVolumes,
				removeLinks:   *removeLinks,
			})
		}
	})

	app.Command("volumes", "Clean useless volumes", func(cmd *cli.Cmd) {
		force := cmd.BoolOpt("force F", false, "Force remove volumes")
		cmd.Action = func() {
			cleanVolumes(cleanVolumesOptions{
				sharedOptions: getSharedOptions(),
				force:         *force,
			})
		}
	})

	app.Command("version", "Print version", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			fmt.Printf("Version: %s\nBuild time: %s\nCommit: %s\n", version, buildTime, commit)
		}
	})

	app.Run(os.Args)
}

type sharedOptions struct {
	client     *client.Client
	safePeriod time.Duration
	dryRun     bool
}
