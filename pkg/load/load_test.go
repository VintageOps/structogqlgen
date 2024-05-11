package load

import (
	"os"
	"testing"
)

func TestFindStructsInPkg(t *testing.T) {
	var tests = []struct {
		sourceFilePath string
		expectedError  string
		expectedLen    int
	}{
		{"valid.go", "", 1},
		{"empty.go", "no structs found", 0},
		{"invalid.go", "failed to parsed the file, error was: invalid.go:1:1: expected 'package', found invalid", 0},
	}

	for _, testcase := range tests {
		t.Run(testcase.sourceFilePath, func(t *testing.T) {
			// Run the GetStructsFromSourceFile function
			result, err := GetStructsFromSourceFile(testcase.sourceFilePath)

			if err != nil && err.Error() != testcase.expectedError {
				t.Errorf("expected error '%s', got '%s'", testcase.expectedError, err)
			}

			if len(result) != testcase.expectedLen {
				t.Errorf("expected %d structs, got %d", testcase.expectedLen, len(result))
			}
		})
	}
}

func TestMain(m *testing.M) {
	// generate test data files
	_ = os.WriteFile("valid.go", []byte("package foo; type Bar struct { Counter int }"), 0600)
	_ = os.WriteFile("empty.go", []byte("package foo; type Foo int"), 0600)
	_ = os.WriteFile("invalid.go", []byte("invalid syntax"), 0600)

	// Run the test suite
	retCode := m.Run()

	// clean up test data files
	_ = os.Remove("valid.go")
	_ = os.Remove("empty.go")
	_ = os.Remove("invalid.go")

	// pass on the exit code
	os.Exit(retCode)
}
