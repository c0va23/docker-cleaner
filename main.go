package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
)

func main() {
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

	for _, image := range images {
		imageUsed := false
		for _, container := range containers {
			if container.ImageID == image.ID {
				imageUsed = true
				break
			}
		}

		if !imageUsed {
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