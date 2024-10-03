package runner

import "testing"

func TestWriteToCache(t *testing.T) {
	fileName := "test_runner.json"
	files := make([]string, 5)
	files[0] = "Hello"
	files[1] = "HelloWorld"
	files[2] = "WorldHello"
	files[3] = "ReaLWorld"
	files[4] = "lolHello"

	if err := saveCache(fileName, files); err != nil {
		t.Error(err)
	}

	readFiles := loadCache(fileName)
	if len(readFiles) != len(files) {
		t.Error("Files doesnt match")
	}
}
