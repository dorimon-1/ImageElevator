package ftp

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/ulikunitz/xz"
)



func setupDecompress(t *testing.T, testData string) (string, func()){
    tempXZFile, err := os.Create("testfile.xz")
    if err != nil {
        t.Fatal("Failed to create temporary XZ file:", err)
    }
    defer os.Remove(tempXZFile.Name())

    var buffer bytes.Buffer
    xzWriter, err := xz.NewWriter(&buffer)
    if err != nil {
        t.Fatal("compression of test-file failed:", err)
    }
    if _, err := io.WriteString(xzWriter, testData); err != nil {
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
    }

    return tempXZFile.Name(), cleanup
}

func TestDecompress(t *testing.T) {
    testData := string("Test data for decompression")
    inputFile, cleanup := setupDecompress(t,testData)
    defer cleanup()

    decompressedFilePath, err := decompress(inputFile)
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