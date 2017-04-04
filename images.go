package main

import (
	"context"
	"log"
	"sort"
	"time"

	"github.com/docker/docker/api/types"
)

type cleanImagesOptions struct {
	sharedOptions
	safePeriod time.Duration
	dryRun     bool
}

func cleanImages(options cleanImagesOptions) {
	containers, err := options.client.ContainerList(
		context.Background(),
		types.ContainerListOptions{All: true},
	)
	if nil != err {
		log.Fatal(err)
	}

	log.Printf("Containers count %d", len(containers))

	images, err := options.client.ImageList(
		context.Background(),
		types.ImageListOptions{All: true},
	)
	if nil != err {
		log.Fatal(err)
	}

	log.Printf("Images count %d", len(images))

	timeLimit := time.Now().Add(-options.safePeriod)
	log.Printf("Time limit %s", timeLimit)

	uselessImageIDs := findUselessImages(findUselessImagesOptions{
		timeLimit:  timeLimit,
		images:     images,
		containers: containers,
	})

	if !options.dryRun {
		removeImages(removeImagesOptions{
			cleanImagesOptions: options,
			imageIDs:           uselessImageIDs,
		})
	}
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
			log.Printf("Image %s too fresh", image.ID)
			continue
		}

		imageUsed := false
		for _, container := range options.containers {
			if container.ImageID == image.ID {
				imageUsed = true
				log.Printf("Image %s used by container %s", image.ID, container.ID)
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
					log.Printf("Image %s used by image %s", image.ID, childImage.ParentID)
					break
				}
			}
		}

		if !imageUsed {
			log.Printf("Image %s useless", image.ID)
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
			log.Printf("Error remove image %s: %s", imageID, err)
		}

		for _, deleteItem := range response {
			if "" != deleteItem.Deleted {
				log.Printf("Delete %s", deleteItem.Deleted)
			} else {
				log.Printf("Untagged %s", deleteItem.Untagged)
			}
		}
	}
}
