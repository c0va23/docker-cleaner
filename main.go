package main

import (
	"os"
	"time"

	cli "github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("declean", "Docker universal cleaner")
	safePeriod := app.IntOpt("safe-period", 0, "Save period")

	app.Command("images", "Clean useless images", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			images(imagesOptions{
				sharedOptions{
					safePeriod: time.Duration(*safePeriod),
				},
			})
		}
	})

	app.Run(os.Args)
}

type sharedOptions struct {
	safePeriod time.Duration
}
