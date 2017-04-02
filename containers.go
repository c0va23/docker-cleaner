package main

import (
	"context"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type cleanContainersOptions struct {
	sharedOptions
}

func cleanContainers(client *client.Client, options cleanContainersOptions) {
	containers, err := client.ContainerList(
		context.Background(),
		types.ContainerListOptions{All: true},
	)
	if nil != err {
		log.Fatal(err)
	}

	uselessContainerIDs := findUselessContainers(findUselessContainersOptions{
		timeLimit:  time.Now().Add(-options.safePeriod * time.Second),
		containers: containers,
	})

	if !options.dryRun {
		removeContainers(client, uselessContainerIDs)
	}
}

type findUselessContainersOptions struct {
	timeLimit  time.Time
	containers []types.Container
}

func containerName(container types.Container) string {
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

func removeContainers(client *client.Client, containerss []types.Container) {
	for _, container := range containerss {
		err := client.ContainerRemove(
			context.Background(),
			container.ID,
			types.ContainerRemoveOptions{},
		)
		if nil != err {
			log.Printf("Error remove container %s: %s", containerName(container), err)
		} else {
			log.Printf("Container %s removed", containerName(container))
		}
	}
}
