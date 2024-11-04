package decompress

import (
	"github.com/ulikunitz/xz"
	"io"
	"os"
	"strings"
)

func Decompress(inputFilePath string) (string, error) {
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
