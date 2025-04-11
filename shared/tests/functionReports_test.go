package shared_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"cryptobotmanager.com/cbm-backend/microservices/backTesting/functions"
)

func TestExtractFunctionNames(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Create some dummy .go files with function declarations
	fileContents := `package main

func Func1() {}
func Func2() {}
`
	err := ioutil.WriteFile(filepath.Join(tmpDir, "file1.go"), []byte(fileContents), 0644)
	if err != nil {
		t.Fatalf("failed to create file1.go: %v", err)
	}

	fileContents = `package main

func Func3() {}
`
	err = ioutil.WriteFile(filepath.Join(tmpDir, "file2.go"), []byte(fileContents), 0644)
	if err != nil {
		t.Fatalf("failed to create file2.go: %v", err)
	}

	// Call the function being tested
	functionNames, err := functions.ExtractFunctionNames(tmpDir)
	if err != nil {
		t.Fatalf("ExtractFunctionNames returned an error: %v", err)
	}

	// Define the expected function names
	expected := []string{"Func1", "Func2", "Func3"}

	// Check if the result matches the expected function names
	if !reflect.DeepEqual(functionNames, expected) {
		t.Errorf("ExtractFunctionNames returned %v, expected %v", functionNames, expected)
	}
}

func TestExtractFuncNamesFromTests(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Create the "tests" directory
	testsDir := filepath.Join(tmpDir, "tests")
	err := os.Mkdir(testsDir, 0755)
	if err != nil {
		t.Fatalf("failed to create tests directory: %v", err)
	}

	// Create some dummy _test.go files with test function declarations
	fileContents := `package main_test

func TestFunc1(t *testing.T) {}
func TestFunc2(t *testing.T) {}
`
	err = ioutil.WriteFile(filepath.Join(testsDir, "file1_test.go"), []byte(fileContents), 0644)
	if err != nil {
		t.Fatalf("failed to create file1_test.go: %v", err)
	}

	fileContents = `package main_test

func TestFunc3(t *testing.T) {}
`
	err = ioutil.WriteFile(filepath.Join(testsDir, "file2_test.go"), []byte(fileContents), 0644)
	if err != nil {
		t.Fatalf("failed to create file2_test.go: %v", err)
	}

	// Call the function being tested
	functionNames, err := functions.ExtractFuncNamesFromTests(testsDir)
	if err != nil {
		t.Fatalf("ExtractFuncNamesFromTests returned an error: %v", err)
	}

	// Define the expected function names
	expected := []string{"Func1", "Func2", "Func3"}

	// Check if the result matches the expected function names
	if !reflect.DeepEqual(functionNames, expected) {
		t.Errorf("ExtractFuncNamesFromTests returned %v, expected %v", functionNames, expected)
	}
}

func TestCalculateCoverage(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		funcNames      []string
		testFuncNames  []string
		expectedResult float64
	}{
		{
			name:           "Case 1: No functions or tests",
			funcNames:      []string{},
			testFuncNames:  []string{},
			expectedResult: 100.0,
		},
		{
			name:           "Case 2: Some functions without test coverage",
			funcNames:      []string{"func1", "func2", "func3"},
			testFuncNames:  []string{"TestFunc1", "TestFunc2"},
			expectedResult: 33.3,
		},
		{
			name:           "Case 3: All functions have test coverage",
			funcNames:      []string{"func1", "func2", "func3"},
			testFuncNames:  []string{"TestFunc1", "TestFunc2", "TestFunc3"},
			expectedResult: 0.0,
		},
	}

	// Run test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := functions.CalculateCoverage(test.funcNames, test.testFuncNames)
			if result != test.expectedResult {
				t.Errorf("Test case '%s' failed: expected %f, got %f", test.name, test.expectedResult, result)
			}
		})
	}
}
