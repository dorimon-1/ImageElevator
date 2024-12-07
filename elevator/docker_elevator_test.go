package elevator

import (
	"context"
	"testing"

	"github.com/Kjone1/imageElevator/config"
	"github.com/Kjone1/imageElevator/docker"
	"github.com/Kjone1/imageElevator/mocks"
	"github.com/stretchr/testify/assert"
)

func setupElevator() (*DockerElevator, *mocks.MockFTPClient, *mocks.MockRegistry) {
	ftpClient := new(mocks.MockFTPClient)
	dockerRegistry := new(mocks.MockRegistry)
	elevatorConfig := &config.ElevatorConfiguration{
		SampleRateInMinutes: 1,
		MaxUploadsPerRun:    5,
	}

	baseElevator := NewBaseElevator(elevatorConfig.SampleRateInMinutes, ftpClient, "", "", nil, elevatorConfig.MaxUploadsPerRun)
	elevator := NewDockerElevator(context.Background(), baseElevator, dockerRegistry, elevatorConfig)
	return elevator, ftpClient, dockerRegistry
}

func TestUploadImagesNoNewImages(t *testing.T) {
	expectedCount := 0
	elevator, ftpClient, _ := setupElevator()
	ftpClient.On("List").Return(nil, nil)
	count, err := elevator.uploadImages()

	assert.Nil(t, err)
	assert.Equal(t, expectedCount, count)
}

func TestUploadImagesViaTrigger(t *testing.T) {

	foundFiles := make([]string, 1)
	foundFiles[0] = "testfile.tar"

	elevator, ftpClient, registry := setupElevator()
	ftpClient.On("List").Return(foundFiles, nil)
	ftpClient.On("Pull").Return(foundFiles, nil)
	registry.On("PushTar").Return(nil)
	registry.On("Sync").Return(nil)

	count, err := elevator.uploadImages()

	assert.Nil(t, err)
	assert.Equal(t, len(foundFiles), count)
}

type TarsToImageTest struct {
	name     string
	tarFiles []string
	want     []docker.Image
}

func NewTarToImageTest(name string, tarFiles []string, want []docker.Image) TarsToImageTest {
	return TarsToImageTest{
		name:     name,
		tarFiles: tarFiles,
		want:     want,
	}
}

func Test_tarsToImages(t *testing.T) {
	tests := []TarsToImageTest{
		NewTarToImageTest("correct_image", []string{"cms-client-5.1.1-hf.2-docker.tar"}, []docker.Image{{
			Name:    "cms-client",
			Tag:     "5.1.1-hf.2",
			TarPath: "cms-client-5.1.1-hf.2-docker.tar",
		}}),
		NewTarToImageTest("correct_image_int", []string{"int-msp-5.1.1-hf.2-docker.tar"}, []docker.Image{{
			Name:    "msp",
			Tag:     "5.1.1-hf.2",
			TarPath: "int-msp-5.1.1-hf.2-docker.tar",
		}}),
		NewTarToImageTest("correct_image_int", []string{"msp-5.1.1-hf.2.tar"}, []docker.Image{{
			Name:    "msp",
			Tag:     "5.1.1-hf.2",
			TarPath: "msp-5.1.1-hf.2.tar",
		}}),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			received := tarsToImages(tt.tarFiles)
			assert.Equal(t, tt.want, received)
		})
	}
}
