package volumes

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// CleanOptions is options for function .Clean
type CleanOptions struct {
	DockerClient client.CommonAPIClient
	DryRun       bool
	Force        bool
}

// Clean useless volumes
func Clean(options CleanOptions) {
	args := filters.NewArgs()
	args.Add("dangling", "true")
	volumesList, err := options.DockerClient.
		VolumeList(context.Background(), args)
	if nil != err {
		log.Fatal(err)
	}
	for _, warning := range volumesList.Warnings {
		log.Printf("Warning: %s", warning)
	}

	uselessVolumes := findUselessVolumes(findUselessVolumesOptions{
		volumes: volumesList.Volumes,
	})

	if !options.DryRun {
		removeVolumes(removeVolumesOptions{
			CleanOptions: options,
			volumes:      uselessVolumes,
		})
	}
}

type findUselessVolumesOptions struct {
	volumes []*types.Volume
}

func findUselessVolumes(options findUselessVolumesOptions) []types.Volume {
	uselessVolumes := []types.Volume{}

	for _, volume := range options.volumes {
		if volume == nil {
			log.Printf("Skip nil volume")
			continue
		}

		if nil != volume.UsageData && volume.UsageData.RefCount > 0 {
			log.Printf("Volume %s used %d times", volume.Name, volume.UsageData.RefCount)
			continue
		}

		log.Printf("Volume %s is useless", volume.Name)

		uselessVolumes = append(uselessVolumes, *volume)
	}

	return uselessVolumes
}

type removeVolumesOptions struct {
	CleanOptions
	volumes []types.Volume
}

func removeVolumes(options removeVolumesOptions) {
	for _, volume := range options.volumes {
		err := options.DockerClient.VolumeRemove(
			context.Background(),
			volume.Name,
			options.Force,
		)
		if nil == err {
			log.Printf("Volume %s removed", volume.Name)
		} else {
			log.Printf("Error remove volume %s: %s", volume.Name, err)
		}
	}
}
