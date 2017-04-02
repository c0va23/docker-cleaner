package main

import (
	"os"
	"time"

	"github.com/docker/docker/client"
	cli "github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("declean", "Docker universal cleaner")
	safePeriod := app.IntOpt("safe-period", 0, "Save period")
	dryRun := app.BoolOpt("dry-run", false, "Dry run")

	client, err := client.NewEnvClient()
	if nil != err {
		panic(err)
	}

	app.Command("images", "Clean useless images", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			images(
				client,
				imagesOptions{
					sharedOptions{
						safePeriod: time.Duration(*safePeriod),
						dryRun:     *dryRun,
					},
				},
			)
		}
	})

	app.Command("containers", "Clean containers", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			cleanContainers(client, cleanContainersOptions{
				sharedOptions{
					safePeriod: time.Duration(*safePeriod),
					dryRun:     *dryRun,
				},
			})
		}
	})

	app.Run(os.Args)
}

type sharedOptions struct {
	safePeriod time.Duration
	dryRun     bool
}
