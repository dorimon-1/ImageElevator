package runner

import (
	"testing"

	"github.com/Kjone1/imageElevator/docker"
	"github.com/stretchr/testify/assert"
)

type TestData struct {
	name     string
	tarFiles []string
	want     []docker.Image
}

func NewTest(name string, tarFiles []string, want []docker.Image) TestData {
	return TestData{
		name:     name,
		tarFiles: tarFiles,
		want:     want,
	}
}

func Test_tarsToImages(t *testing.T) {
	tests := []TestData{
		NewTest("correct_image", []string{"cms-client-5.1.1-hf.2-docker.tar"}, []docker.Image{{
			Name:    "cms-client",
			Tag:     "5.1.1-hf.2",
			TarPath: "cms-client-5.1.1-hf.2-docker.tar",
		}}),
		NewTest("correct_image_int", []string{"int-msp-5.1.1-hf.2-docker.tar"}, []docker.Image{{
			Name:    "msp",
			Tag:     "5.1.1-hf.2",
			TarPath: "int-msp-5.1.1-hf.2-docker.tar",
		}}),
		NewTest("correct_image_int", []string{"msp-5.1.1-hf.2.tar"}, []docker.Image{{
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
