package main

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type imagesOptions struct {
	sharedOptions
}

func images(client *client.Client, imagesOptions imagesOptions) {
	containers, err := client.ContainerList(
		context.Background(),
		types.ContainerListOptions{All: true},
	)
	if nil != err {
		panic(err)
	}

	fmt.Printf("Containers count %d\n", len(containers))

	images, err := client.ImageList(
		context.Background(),
		types.ImageListOptions{All: true},
	)
	if nil != err {
		panic(err)
	}

	fmt.Printf("Images count %d\n", len(images))

	timeLimit := time.Now().Add(-imagesOptions.safePeriod)
	fmt.Printf("Time limit %s\n", timeLimit)

	uselessImageIDs := findUselessImages(findUselessImagesOptions{
		timeLimit:  timeLimit,
		images:     images,
		containers: containers,
	})

	removeImages(client, uselessImageIDs)
}

type findUselessImagesOptions struct {
	timeLimit  time.Time
	images     []types.ImageSummary
	containers []types.Container
}

func findUselessImages(options findUselessImagesOptions) []string {
	uselessImageIDs := []string{}

	for _, image := range options.images {
		imageCreated := time.Unix(image.Created, 0)

		if imageCreated.After(options.timeLimit) {
			fmt.Printf("Image %s too fresh\n", image.ID)
			continue
		}

		imageUsed := false
		for _, container := range options.containers {
			if container.ImageID == image.ID {
				imageUsed = true
				fmt.Printf("Image %s used by container %s\n", image.ID, container.ID)
				break
			}
		}

		for _, childImage := range options.images {
			if childImage.ParentID == image.ID {
				imageUsed = true
				fmt.Printf("Image %s used by image %s\n", image.ID, childImage.ParentID)
				break
			}
		}

		if !imageUsed {
			fmt.Printf("Image %s useless\n", image.ID)
			uselessImageIDs = append(uselessImageIDs, image.ID)
		}
	}
	return uselessImageIDs
}

func removeImages(client *client.Client, imageIDs []string) {
	for _, imageID := range imageIDs {
		response, err := client.ImageRemove(
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
