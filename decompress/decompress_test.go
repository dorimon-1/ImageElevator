package decompress

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/ulikunitz/xz"
)

func setupDecompress(t *testing.T, testData string) (string, func()) {
	tempXZFile, err := os.Create("testfile.xz")
	if err != nil {
		t.Fatal("Failed to create temporary XZ file:", err)
	}

	var buffer bytes.Buffer
	xzWriter, err := xz.NewWriter(&buffer)
	if err != nil {
		t.Fatal("compression of test-file failed:", err)
	}
	if _, err := io.WriteString(xzWriter, testData); err != nil {
		// TODO: write unique error msg
		t.Fatal("compression of test-file failed:", err)
	}
	if err := xzWriter.Close(); err != nil {
		t.Fatal("failed to close xzWriter, error:", err)
	}

	if _, err := tempXZFile.Write(buffer.Bytes()); err != nil {
		t.Fatal("failed to write compress data to file, error:", err)
	}

	cleanup := func() {
		os.Remove(tempXZFile.Name())
		log.Info().Msgf("Temp test file has been deleted")
	}

	log.Info().Msgf("Decompress setup had finished sucessfuly")
	return tempXZFile.Name(), cleanup
}

func TestDecompress(t *testing.T) {
	testData := "Test data for decompression"
	inputFile, cleanup := setupDecompress(t, testData)
	defer cleanup()

	decompressor := new(TarDecompressor)

	decompressedFilePath, err := decompressor.Decompress(inputFile)
	if err != nil {
		t.Fatal("Decompression failed:", err)
	}

	decompressedData, err := os.ReadFile(decompressedFilePath)
	if err != nil {
		t.Fatal("Failed to read decompressed file:", err)
	}

	if string(decompressedData) != string(testData) {
		t.Errorf("Decompressed data does not match original test data. Expected: %s, Recieved: %s", testData, decompressedData)
	}

	if decompressedFilePath != "testfile" {
		t.Errorf("inputFilePath do not much expectation. Expected: %s, Recieved: %s", testData, decompressedData)
	}
}
