package ftp_test

import (
	"bytes"
	"io"
	"os"
	"slices"
	"testing"

	"github.com/Kjone1/imageElevator/ftp"
	"github.com/rs/zerolog/log"
	"github.com/secsy/goftp"
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

	decompressedFilePath, err := ftp.Decompress(inputFile)
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

var client *goftp.Client

func SetupClient(t *testing.T) *goftp.Client {
	if client == nil {
		ftpClient, err := ftp.Connect("localhost", "testuser", "testpassword")
		if err != nil {
			t.Errorf("Failed With FTP client setup => %v", err)
		}
		client = ftpClient
	}
	return client
}
func TestConnect(t *testing.T) {
	client := SetupClient(t)
	if client == nil {
		t.Error("Failed to Connect to FTP server")
	}
}

func TestList(t *testing.T) {
	client := SetupClient(t)
	testdir, err := client.Mkdir("test-dir")
	if err != nil {
		t.Errorf("Failed creating testdir for TestList => %v", err)
	}
	test_list, err := ftp.List(client, "/", "test.*")
	if err != nil {
		t.Errorf("Reading FTP directory failed with error => %s", err)
	}
	if !slices.Contains(test_list, testdir) {
		t.Errorf("Test dir not found in FTP server")
	}
}
