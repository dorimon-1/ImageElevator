package decompress

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/ulikunitz/xz"
)

type Decompressor interface {
	Decompress(string) (string, error)
}

type XZDecompressor struct {
}

type TarDecompressor struct {
}

func (TarDecompressor) Decompress(inputFilePath string) (string, error) {
	if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("archive file does not exist: %s", inputFilePath)
	}

	outputFilePath := strings.TrimSuffix(inputFilePath, ".xz")
	// Prepare the tar command
	cmd := exec.Command("tar", "-xJf", inputFilePath)

	// Capture the output for debugging
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to decompress archive: %v", err)
	}

	return outputFilePath, nil
}

func (XZDecompressor) Decompress(inputFilePath string) (string, error) {
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return "", err
	}
	defer inputFile.Close()

	outputFilePath := strings.TrimSuffix(inputFilePath, ".xz")
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return "", err
	}
	defer outputFile.Close()

	xzReader, err := xz.NewReader(inputFile)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(outputFile, xzReader)
	if err != nil {
		return "", err
	}

	return outputFilePath, nil
}
