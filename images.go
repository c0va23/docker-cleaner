package main

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/docker/docker/api/types"
)

type cleanImagesOptions struct {
	sharedOptions
}

func cleanImages(options cleanImagesOptions) {
	containers, err := options.client.ContainerList(
		context.Background(),
		types.ContainerListOptions{All: true},
	)
	if nil != err {
		panic(err)
	}

	fmt.Printf("Containers count %d\n", len(containers))

	images, err := options.client.ImageList(
		context.Background(),
		types.ImageListOptions{All: true},
	)
	if nil != err {
		panic(err)
	}

	fmt.Printf("Images count %d\n", len(images))

	timeLimit := time.Now().Add(-options.safePeriod)
	fmt.Printf("Time limit %s\n", timeLimit)

	uselessImageIDs := findUselessImages(findUselessImagesOptions{
		timeLimit:  timeLimit,
		images:     images,
		containers: containers,
	})

	removeImages(removeImagesOptions{
		cleanImagesOptions: options,
		imageIDs:           uselessImageIDs,
	})
}

type findUselessImagesOptions struct {
	timeLimit  time.Time
	images     []types.ImageSummary
	containers []types.Container
}

func findUselessImages(options findUselessImagesOptions) []string {
	uselessImageIDs := []string{}

	sort.Slice(options.images, func(i, j int) bool {
		imageI := options.images[i]
		imageJ := options.images[j]
		return imageI.Created > imageJ.Created ||
			(imageI.Created == imageJ.Created && imageI.ParentID == imageJ.ID)
	})

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
				childImageUseless := false
				for _, uselessImageID := range uselessImageIDs {
					if childImage.ID == uselessImageID {
						childImageUseless = true
						break
					}
				}
				if !childImageUseless {
					imageUsed = true
					fmt.Printf("Image %s used by image %s\n", image.ID, childImage.ParentID)
					break
				}
			}
		}

		if !imageUsed {
			fmt.Printf("Image %s useless\n", image.ID)
			uselessImageIDs = append(uselessImageIDs, image.ID)
		}
	}
	return uselessImageIDs
}

type removeImagesOptions struct {
	cleanImagesOptions
	imageIDs []string
}

func removeImages(options removeImagesOptions) {
	for _, imageID := range options.imageIDs {
		response, err := options.client.ImageRemove(
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
