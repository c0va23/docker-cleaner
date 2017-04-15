package containers

import (
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
)

func genContainerID() string {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); nil != err {
		panic(err)
	}
	return fmt.Sprintf("%x", buf)
}

func TestFindUselessContainers(t *testing.T) {
	freshContainer := types.Container{
		ID:      genContainerID(),
		Created: time.Now().Unix(),
	}
	runningContainer := types.Container{
		ID:      genContainerID(),
		Created: time.Now().Add(-2 * time.Hour).Unix(),
		State:   "running",
	}
	uselessContainer := types.Container{
		ID:      genContainerID(),
		Created: time.Now().Add(-2 * time.Hour).Unix(),
		State:   "exited",
	}

	options := findUselessContainersOptions{
		timeLimit: time.Now().Add(-time.Hour),
		containers: []types.Container{
			freshContainer,
			runningContainer,
			uselessContainer,
		},
	}

	uselessContainers := findUselessContainers(options)

	if 1 == len(uselessContainers) && uselessContainers[0].ID == uselessContainer.ID {
		return
	}
	t.Errorf("findUselessContainers not return useless container")
}

func TestContainerName(t *testing.T) {
	nameless := containerName(types.Container{
		Names: []string{},
	})

	if "" != nameless {
		t.Errorf("Not return empty string for nameless container")
	}

	name := containerName(types.Container{
		Names: []string{
			"/golang",
		},
	})

	if "golang" != name {
		t.Errorf("containerName return invalid name")
	}
}
