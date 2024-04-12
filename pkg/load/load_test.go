package load

import (
	"os"
	"testing"
)

func TestFindStructsInPkg(t *testing.T) {
	tests := []struct {
		name      string
		file      string
		wantError bool
	}{
		{
			name:      "existing file with structs",
			file:      "./testdata/structs.go",
			wantError: false,
		},
		{
			name:      "existing file without structs",
			file:      "./testdata/no_structs.go",
			wantError: true,
		},
		{
			name:      "non-existing file",
			file:      "./invalid/path.go",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FindStructsInPkg(tt.file)
			if (err != nil) != tt.wantError {
				t.Errorf("FindStructsInPkg() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	os.MkdirAll("./testdata", os.ModePerm)
	writeTestFile("./testdata/structs.go", "package testdata\n\n type TestStruct struct{}")
	writeTestFile("./testdata/no_structs.go", "package testdata")
}

func teardown() {
	os.RemoveAll("./testdata")
}

func writeTestFile(path, content string) {
	file, _ := os.Create(path)
	defer file.Close()
	_, _ = file.WriteString(content)
}
