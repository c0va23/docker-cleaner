package main

import (
	"context"
	"log"
	"time"

	"github.com/docker/docker/api/types"
)

type cleanContainersOptions struct {
	sharedOptions
	safePeriod    time.Duration
	dryRun        bool
	removeVolumes bool
	removeLinks   bool
}

func cleanContainers(options cleanContainersOptions) {
	containers, err := options.client.ContainerList(
		context.Background(),
		types.ContainerListOptions{All: true},
	)
	if nil != err {
		log.Fatal(err)
	}

	uselessContainers := findUselessContainers(findUselessContainersOptions{
		timeLimit:  time.Now().Add(-options.safePeriod * time.Second),
		containers: containers,
	})

	if !options.dryRun {
		removeContainers(removeContainerOptions{
			cleanContainersOptions: options,
			containers:             uselessContainers,
		})
	}
}

type findUselessContainersOptions struct {
	timeLimit  time.Time
	containers []types.Container
}

func containerName(container types.Container) string {
	if len(container.Names) == 0 || len(container.Names[0]) < 2 {
		return ""
	}
	return container.Names[0][1:]
}

func findUselessContainers(options findUselessContainersOptions) []types.Container {
	log.Printf("Time limit %s", options.timeLimit)

	uselessContainers := []types.Container{}

	for _, container := range options.containers {
		containerTime := time.Unix(container.Created, 0)

		if containerTime.After(options.timeLimit) {
			log.Printf("Container %s too freshs", containerName(container))
			continue
		}

		if container.State == "running" {
			log.Printf("Container %s is runnings", containerName(container))
			continue
		}

		uselessContainers = append(uselessContainers, container)
		log.Printf("Container %s useless", containerName(container))
	}

	return uselessContainers
}

type removeContainerOptions struct {
	cleanContainersOptions
	containers []types.Container
}

func removeContainers(options removeContainerOptions) {
	containerRemoveOptions := types.ContainerRemoveOptions{
		RemoveVolumes: options.removeVolumes,
		RemoveLinks:   options.removeLinks,
	}

	for _, container := range options.containers {
		err := options.client.ContainerRemove(
			context.Background(),
			container.ID,
			containerRemoveOptions,
		)
		if nil != err {
			log.Printf("Error remove container %s: %s", containerName(container), err)
		} else {
			log.Printf("Container %s removed", containerName(container))
		}
	}
}
