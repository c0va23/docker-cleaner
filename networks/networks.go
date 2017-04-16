package networks

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// CleanOptions is options for function .Clean
type CleanOptions struct {
	DockerClient client.CommonAPIClient
	DryRun       bool
}

// Clean useless networks
func Clean(options CleanOptions) error {
	networks, err := findUseless(findOptions{
		dockerClient: options.DockerClient,
	})

	if nil != err {
		return err
	}

	if options.DryRun {
		return nil
	}

	return remove(removeOptions{
		dockerClient: options.DockerClient,
		networks:     networks,
	})
}

type findOptions struct {
	dockerClient client.CommonAPIClient
}

func findUseless(options findOptions) ([]types.NetworkResource, error) {
	allNetworks, err := options.dockerClient.NetworkList(
		context.Background(),
		types.NetworkListOptions{},
	)
	if nil != err {
		return nil, err
	}
	uselessNetworks := []types.NetworkResource{}
	for _, network := range allNetworks {
		if "none" == network.Name || "bridge" == network.Name || "host" == network.Name {
			log.Printf("Skip default network %s", network.Name)
			continue
		}
		log.Printf("Netowrk %s (%+v): %+v", network.Name, network.Labels, network.Containers)
		if 0 == len(network.Containers) {
			uselessNetworks = append(uselessNetworks, network)
		}
	}
	return uselessNetworks, nil
}

type removeOptions struct {
	dockerClient client.CommonAPIClient
	networks     []types.NetworkResource
}

func remove(options removeOptions) error {
	for _, network := range options.networks {
		err := options.dockerClient.NetworkRemove(
			context.Background(),
			network.ID,
		)
		if nil != err {
			return err
		}
	}
	return nil
}
