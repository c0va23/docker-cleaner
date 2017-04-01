package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
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

type imagesOptions struct {
	sharedOptions
}

func images(imagesOptions imagesOptions) {
	cli, err := client.NewEnvClient()
	if nil != err {
		panic(err)
	}

	containers, err := cli.ContainerList(
		context.Background(),
		types.ContainerListOptions{All: true},
	)
	if nil != err {
		panic(err)
	}

	fmt.Printf("Containers count %d\n", len(containers))

	images, err := cli.ImageList(
		context.Background(),
		types.ImageListOptions{All: true},
	)
	if nil != err {
		panic(err)
	}

	fmt.Printf("Images count %d\n", len(images))

	uselessImageIds := []string{}

	timeLimit := time.Now().Truncate(imagesOptions.safePeriod)
	fmt.Printf("Time limit %s\n", timeLimit)

	for _, image := range images {
		imageCreated := time.Unix(image.Created, 0)

		if imageCreated.After(timeLimit) {
			fmt.Printf("Image %s too fresh\n", image.ID)
			continue
		}

		imageUsed := false
		for _, container := range containers {
			if container.ImageID == image.ID {
				imageUsed = true
				fmt.Printf("Image %s used by container %s\n", image.ID, container.ID)
				break
			}
		}

		for _, childImage := range images {
			if childImage.ParentID == image.ID {
				imageUsed = true
				fmt.Printf("Image %s used by image %s\n", image.ID, childImage.ParentID)
				break
			}
		}

		if !imageUsed {
			fmt.Printf("Image %s useless\n", image.ID)
			uselessImageIds = append(uselessImageIds, image.ID)
		}
	}

	for _, imageID := range uselessImageIds {
		response, err := cli.ImageRemove(
			context.Background(),
			imageID,
			types.ImageRemoveOptions{},
		)
		if nil != err {
			fmt.Printf("Err: %s\n", err)
		}

		fmt.Printf("Response: %+v\n", response)
	}
}
