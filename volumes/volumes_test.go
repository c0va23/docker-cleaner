package volumes

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
)

func genVolumeID() string {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); nil != err {
		panic(err)
	}
	return fmt.Sprintf("%x", buf)
}

func TestFindUselessVolumes(t *testing.T) {
	usedVolume := types.Volume{
		Name: genVolumeID(),
		UsageData: &types.VolumeUsageData{
			RefCount: 1,
		},
	}

	notUsedVolume := types.Volume{
		Name: genVolumeID(),
		UsageData: &types.VolumeUsageData{
			RefCount: 0,
		},
	}

	withoutUsageVolume := types.Volume{
		Name:      genVolumeID(),
		UsageData: nil,
	}

	options := findUselessVolumesOptions{
		volumes: []*types.Volume{
			nil,
			&usedVolume,
			&notUsedVolume,
			&withoutUsageVolume,
		},
	}
	uselessVolumes := findUselessVolumes(options)

	if 2 != len(uselessVolumes) ||
		uselessVolumes[0].Name != notUsedVolume.Name ||
		uselessVolumes[1].Name != withoutUsageVolume.Name {
		t.Errorf("findUselessVolumes return invalid result (%+v)", uselessVolumes)
	}
}
