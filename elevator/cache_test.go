package elevator

import "testing"

func TestWriteToCache(t *testing.T) {
	fileName := "test_elevator.json"
	files := make(map[string]bool)
	files["hello"] = true
	files["helloWorld"] = true
	files["helloworld"] = true
	files["Worldhello"] = true
	files["Hleohello"] = true

	if err := saveCache(fileName, files); err != nil {
		t.Error(err)
	}

	readFiles := loadCache(fileName)
	if len(readFiles) != len(files) {
		t.Error("Files doesnt match")
	}
}
