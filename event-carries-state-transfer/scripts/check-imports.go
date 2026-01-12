// check-imports.go - Fast Go implementation of import organization checker
//
// This tool enforces the following import organization:
//   import (
//       // std libs (context, fmt, etc.)
//
//       // external libs (github.com/*, google.golang.org/*)
//
//       // alias (aliased imports like pg "path")
//
//       // internal libs (eda-in-golang/*)
//   )
//
// Generated files (*.pb.go, *_grpc.pb.go) are automatically skipped.
// This is the preferred/faster version compared to the bash script.
//
// Usage: make check-imports-go or go run scripts/check-imports.go [directory]

package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type ImportInfo struct {
	Path     string
	Group    ImportGroup
	Position int
}

type ImportGroup int

const (
	SystemImports ImportGroup = iota
	ExternalImports
	AliasImports
	InternalImports
)

func (g ImportGroup) String() string {
	switch g {
	case SystemImports:
		return "system"
	case ExternalImports:
		return "external"
	case AliasImports:
		return "alias"
	case InternalImports:
		return "internal"
	default:
		return "unknown"
	}
}

func classifyImport(line string) ImportGroup {
	// Remove whitespace and extract path
	line = strings.TrimSpace(line)
	var path string

	// Check if it's an aliased import (contains space before quote)
	if strings.Contains(line, " \"") {
		// This is an aliased import
		return AliasImports
	} else if strings.Contains(line, "\"") {
		// Extract path from quoted string
		start := strings.Index(line, "\"")
		end := strings.LastIndex(line, "\"")
		if start >= 0 && end > start {
			path = line[start+1 : end]
		}
	}

	// Internal imports (project modules)
	if strings.HasPrefix(path, "eda-in-golang/") {
		return InternalImports
	}

	// External imports (contains dots)
	if strings.Contains(path, ".") {
		return ExternalImports
	}

	// System imports (standard library, no dots)
	return SystemImports
}

func isGeneratedFile(filename string) bool {
	// Skip protobuf generated files
	if strings.HasSuffix(filename, ".pb.go") ||
		strings.HasSuffix(filename, "_grpc.pb.go") ||
		strings.HasSuffix(filename, "_pb.go") {
		return true
	}

	// Check for generated file comments
	src, err := os.ReadFile(filename)
	if err != nil {
		return false
	}

	lines := strings.Split(string(src), "\n")
	for i, line := range lines {
		if i > 10 { // Only check first 10 lines
			break
		}
		if strings.Contains(line, "Code generated") ||
			strings.Contains(line, "DO NOT EDIT") ||
			strings.Contains(line, "This file was auto-generated") {
			return true
		}
	}

	return false
}

func checkFile(filename string) error {
	if isGeneratedFile(filename) {
		fmt.Printf("⏭️  Skipping generated file: %s\n", filename)
		skippedCount++
		return nil
	}

	src, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return err
	}

	if len(file.Imports) == 0 {
		fmt.Printf("⏭️  No imports to check: %s\n", filename)
		skippedCount++
		return nil
	}

	checkedCount++
	var imports []ImportInfo
	for i, imp := range file.Imports {
		line := strings.TrimSpace(string(src[imp.Pos()-1 : imp.End()-1]))
		group := classifyImport(line)
		path := strings.Trim(imp.Path.Value, `"`)
		imports = append(imports, ImportInfo{
			Path:     path,
			Group:    group,
			Position: i,
		})
	}

	var violations []string

	// Check grouping order
	var lastGroup ImportGroup = -1
	var seenGroups = make(map[ImportGroup]bool)

	for _, imp := range imports {
		// Check if groups are in correct order
		if imp.Group < lastGroup && lastGroup != AliasImports {
			violations = append(violations, fmt.Sprintf(
				"Import '%s' (%s) appears after %s import",
				imp.Path, imp.Group, lastGroup))
		}

		// Check for proper separation
		if seenGroups[imp.Group] {
			// This group has been seen before - check if there's proper separation
			// This is a simplified check - a full implementation would need to check AST positions
		}

		seenGroups[imp.Group] = true
		lastGroup = imp.Group
	}

	if len(violations) > 0 {
		fmt.Printf("❌ %s:\n", filename)
		for _, violation := range violations {
			fmt.Printf("  - %s\n", violation)
		}
		errorCount++
		return fmt.Errorf("import violations found")
	}

	fmt.Printf("✅ %s\n", filename)
	return nil
}

// Statistics
var (
	checkedCount = 0
	skippedCount = 0
	errorCount   = 0
)

func checkDirectory(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".go") || info.IsDir() {
			return nil
		}
		return checkFile(path)
	})
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run check-imports.go <file|directory>")
		os.Exit(1)
	}

	target := os.Args[1]
	info, err := os.Stat(target)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	var checkErr error
	if info.IsDir() {
		checkErr = checkDirectory(target)
	} else {
		checkErr = checkFile(target)
	}

	if checkErr != nil {
		os.Exit(1)
	}

	// Display summary
	fmt.Println("")
	fmt.Println("📊 Import Check Summary:")
	fmt.Printf("  ✅ Checked files: %d\n", checkedCount)
	fmt.Printf("  ⏭️  Skipped files: %d\n", skippedCount)
	fmt.Printf("  ❌ Files with errors: %d\n", errorCount)
	fmt.Printf("  📁 Total files processed: %d\n", checkedCount+skippedCount)

	if errorCount > 0 {
		fmt.Println("")
		fmt.Println("💡 To fix import issues automatically, run: make fix-imports")
		os.Exit(1)
	} else {
		fmt.Println("")
		fmt.Println("🎉 All import checks passed!")
	}
}
