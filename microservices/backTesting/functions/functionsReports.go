package functions

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"cryptobotmanager.com/cbm-backend/shared"
	"github.com/rs/zerolog/log"
)

// ExtractFunctionNames traverses through the files in the given package directory,
// parses them, and extracts the names of functions.
func ExtractFunctionNames(packageNames string) ([]string, error) {
	var functionNames []string

	err := filepath.Walk(packageNames, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
		if err != nil {
			return fmt.Errorf("error parsing file %s: %v", path, err)
		}

		// Traverse the AST to find function declarations
		for _, decl := range node.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			functionNames = append(functionNames, fn.Name.Name)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return functionNames, nil
}

// ExtractFuncNames traverses through the files in the given directory,
// parses test files, and extracts function names without the "Test" prefix.
func ExtractFuncNamesFromTests(rootDir string) ([]string, error) {
	var functionNames []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		if strings.HasSuffix(path, "_test.go") {
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
			if err != nil {
				return fmt.Errorf("error parsing file %s: %v", path, err)
			}

			// Traverse the AST to find function declarations
			for _, decl := range node.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}
				// Remove the "Test" prefix from the test function name
				fnName := strings.TrimPrefix(fn.Name.Name, "Test")
				functionNames = append(functionNames, fnName)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return functionNames, nil
}

// CalculateCoverage calculates the percentage of functions without test coverage.
func CalculateCoverage(funcNames, testFuncNames []string) float64 {
	totalFuncs := float64(len(funcNames))
	testedFuncs := float64(len(testFuncNames))
	if totalFuncs == 0 {
		return 100.0 // All functions are considered without test coverage if no functions exist
	}
	percentage := (totalFuncs - testedFuncs) / totalFuncs * 100.0
	return shared.RoundFloatToDecimal(percentage, 1)
}

// PrintFunctionsWithoutTestCoverage prints the names of functions without test coverage.
func PrintFunctionsWithoutTestCoverage() {
	// Load config file
	cfg := shared.GetDefaultCfg()

	var funcNames []string

	fmt.Println("\n All Functions Names:")

	// Process regular files
	for _, packageName := range cfg.PackageNames {
		packageFuncNames, _ := ExtractFunctionNames(packageName)
		fmt.Printf("%s Package \n", packageName)
		for i, funcName := range packageFuncNames {
			fmt.Printf("%d. %s\n", i+1, funcName)
		}
		fmt.Println("")
		funcNames = append(funcNames, packageFuncNames...)
	}

	fmt.Println("\n Functions Names with Test Coverage:")

	testFuncNames, _ := ExtractFuncNamesFromTests("tests")

	for i, funcName := range testFuncNames {
		fmt.Printf("%d. %s\n", i+1, funcName)
	}

	fmt.Println("\n Excempt Functions Names:")
	for i, funcName := range cfg.TestExemptFuncs {
		fmt.Printf("%d. %s\n", i+1, funcName)
	}

	uniqueNames := shared.FindUniqueStrings(funcNames, testFuncNames)
	uniqueNames = shared.FindUniqueStrings(uniqueNames, cfg.TestExemptFuncs)

	coverage := CalculateCoverage(testFuncNames, uniqueNames)

	// Print the names of functions without test coverage
	if len(uniqueNames) == 0 {
		log.Info().Msg("All functions have test coverage!")
	} else {
		fmt.Println("\n Functions without test coverage")
		for i, funcName := range uniqueNames {
			fmt.Printf("%d. %s\n", i+1, funcName)
		}
	}

	fmt.Println("")

	// Print the percentage of functions without test coverage
	log.Warn().Float64("coverage", coverage).Msg("Percentage of functions with test")
	log.Info().Int("functions: ", len(funcNames)).Msg("Total No# of")
	log.Info().Int("functions without test coverage: ", len(uniqueNames)).Msg("No# of")
	log.Info().Int("tested functionss is: ", len(testFuncNames)).Msg("No# of")
	log.Info().Int("exempt functions: ", len(cfg.TestExemptFuncs)).Msg("No# of")

}
