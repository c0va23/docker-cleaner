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
