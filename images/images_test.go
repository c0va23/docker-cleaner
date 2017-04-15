package images

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
)

func genImageID() string {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); nil != err {
		panic(err)
	}
	return fmt.Sprintf("sha256:%x", buf)
}

func checkID(IDs []string, id string) bool {
	for _, otherID := range IDs {
		if id == otherID {
			return true
		}
	}
	return false
}

func TestFindUselessImages_TooFresh(t *testing.T) {
	timeLimit := time.Now().Add(-5 * time.Minute)
	freshImage := types.ImageSummary{
		ID:      genImageID(),
		Created: time.Now().Unix(),
	}
	uselessImageIDs := findUselessImages(findUselessImagesOptions{
		timeLimit: timeLimit,
		images: []types.ImageSummary{
			freshImage,
		},
	})

	if checkID(uselessImageIDs, freshImage.ID) {
		t.Errorf("Useless image IDs contain too fresh image")
	}
}

func TestFindUselessImages_ContainerExist(t *testing.T) {
	timeLimit := time.Now().Add(-5 * time.Minute)
	usedImage := types.ImageSummary{
		ID:      genImageID(),
		Created: timeLimit.Add(-5 * time.Minute).Unix(),
	}
	container := types.Container{
		ImageID: usedImage.ID,
	}
	uselessImageIDs := findUselessImages(findUselessImagesOptions{
		timeLimit: timeLimit,
		images: []types.ImageSummary{
			usedImage,
		},
		containers: []types.Container{
			container,
		},
	})

	if checkID(uselessImageIDs, usedImage.ID) {
		t.Errorf("Useless image IDs contain image  used by container")
	}
}

func TestFindUselessImages_ChildContainerExist(t *testing.T) {
	timeLimit := time.Now().Add(-5 * time.Minute)
	parentImage := types.ImageSummary{
		ID:      genImageID(),
		Created: timeLimit.Add(-5 * time.Minute).Unix(),
	}
	childImage := types.ImageSummary{
		ID:       genImageID(),
		ParentID: parentImage.ID,
		Created:  timeLimit.Add(-5 * time.Minute).Unix(),
	}
	container := types.Container{
		ImageID: childImage.ID,
	}
	uselessImageIDs := findUselessImages(findUselessImagesOptions{
		timeLimit: timeLimit,
		images: []types.ImageSummary{
			parentImage,
			childImage,
		},
		containers: []types.Container{
			container,
		},
	})

	if checkID(uselessImageIDs, parentImage.ID) {
		t.Errorf("Useless image IDs contain image used by other image used by container")
	}
}

func TestFindUselessImages_FullyUseless(t *testing.T) {
	timeLimit := time.Now().Add(-5 * time.Minute)
	image := types.ImageSummary{
		ID:      genImageID(),
		Created: timeLimit.Add(-5 * time.Minute).Unix(),
	}
	uselessImageIDs := findUselessImages(findUselessImagesOptions{
		timeLimit: timeLimit,
		images: []types.ImageSummary{
			image,
		},
	})

	if !checkID(uselessImageIDs, image.ID) {
		t.Errorf("Useless image IDs not contain fully useless image")
	}
}

func TestFindUselessImages_UsedByUselessImage(t *testing.T) {
	timeLimit := time.Now().Add(-5 * time.Minute)
	image := types.ImageSummary{
		ID:      genImageID(),
		Created: timeLimit.Add(-5 * time.Minute).Unix(),
	}
	childImage := types.ImageSummary{
		ID:       genImageID(),
		ParentID: image.ID,
		Created:  timeLimit.Add(-4 * time.Minute).Unix(),
	}
	uselessImageIDs := findUselessImages(findUselessImagesOptions{
		timeLimit: timeLimit,
		images: []types.ImageSummary{
			image,
			childImage,
		},
	})

	if len(uselessImageIDs) != 2 {
		t.Errorf("Useless image IDs not contain both image and  child image")
	} else if uselessImageIDs[0] != childImage.ID {
		t.Errorf("Useless image IDs not contain child image before image")
	} else if uselessImageIDs[1] != image.ID {
		t.Errorf("Useless image IDs not contain image after child image")
	}
}

func TestFindUselessImages_UsedByUselessImageWithEqualTime(t *testing.T) {
	timeLimit := time.Now().Add(-5 * time.Minute)
	image := types.ImageSummary{
		ID:      genImageID(),
		Created: timeLimit.Add(-5 * time.Minute).Unix(),
	}
	childImage := types.ImageSummary{
		ID:       genImageID(),
		ParentID: image.ID,
		Created:  image.Created,
	}
	uselessImageIDs := findUselessImages(findUselessImagesOptions{
		timeLimit: timeLimit,
		images: []types.ImageSummary{
			image,
			childImage,
		},
	})

	if len(uselessImageIDs) != 2 {
		t.Errorf("Useless image IDs not contain both image and child image")
	} else if uselessImageIDs[0] != childImage.ID {
		t.Errorf("Useless image IDs not contain child image before image")
	} else if uselessImageIDs[1] != image.ID {
		t.Errorf("Useless image IDs not contain image after child image")
	}
}
