package main

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

type cleanVolumesOptions struct {
	sharedOptions
	dryRun bool
	force  bool
}

func cleanVolumes(options cleanVolumesOptions) {
	args := filters.NewArgs()
	args.Add("dangling", "true")
	volumesList, err := options.client.VolumeList(context.Background(), args)
	if nil != err {
		log.Fatal(err)
	}
	for _, warning := range volumesList.Warnings {
		log.Printf("Warning: %s", warning)
	}

	uselessVolumes := findUselessVolumes(findUselessVolumesOptions{
		volumes: volumesList.Volumes,
	})

	if !options.dryRun {
		removeVolumes(removeVolumesOptions{
			cleanVolumesOptions: options,
			volumes:             uselessVolumes,
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
	cleanVolumesOptions
	volumes []types.Volume
}

func removeVolumes(options removeVolumesOptions) {
	for _, volume := range options.volumes {
		err := options.client.VolumeRemove(context.Background(), volume.Name, options.force)
		if nil == err {
			log.Printf("Volume %s removed", volume.Name)
		} else {
			log.Printf("Error remove volume %s: %s", volume.Name, err)
		}
	}
}
