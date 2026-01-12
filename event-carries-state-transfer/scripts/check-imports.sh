#!/bin/bash
# check-imports.sh - Check Go import organization rules (bash implementation)
# Usage: ./scripts/check-imports.sh [file|directory]
#
# This script enforces the following import organization:
#   import (
#       // std libs (context, fmt, etc.)
#
#       // external libs (github.com/*, google.golang.org/*)
#
#       // alias (aliased imports like pg "path")
#
#       // internal libs (eda-in-golang/*)
#   )
#
# Generated files (*.pb.go, *_grpc.pb.go) are automatically skipped.
# Note: For better performance, consider using the Go version: make check-imports-go

set -e

TARGET="${1:-.}"

echo "Checking import organization in: $TARGET"

# Statistics
checked_count=0
skipped_count=0
error_count=0

# Function to check imports in a single file
check_file() {
    local file="$1"
    if [[ ! -f "$file" ]] || [[ "$file" != *.go ]]; then
        return
    fi

    # Skip generated files
    if [[ "$file" =~ \.pb\.go$ ]] || [[ "$file" =~ _grpc\.pb\.go$ ]] || [[ "$file" =~ _pb\.go$ ]]; then
        echo "⏭️  Skipping generated file: $file"
        ((skipped_count++))
        return
    fi

    # Skip files with generated comments
    if head -5 "$file" 2>/dev/null | grep -q "Code generated\|DO NOT EDIT\|This file was auto-generated"; then
        echo "⏭️  Skipping generated file: $file"
        ((skipped_count++))
        return
    fi

    # Extract import statements (both block and single-line formats)
    local block_imports=""
    local single_imports=""

    # Extract block imports
    if grep -q '^import (' "$file" 2>/dev/null; then
        block_imports=$(sed -n '/^import (/,/)/p' "$file" 2>/dev/null | sed '1d;$d' 2>/dev/null | grep -v '^$' 2>/dev/null | grep -v '^[[:space:]]*//' 2>/dev/null || true)
    fi

    # Extract single-line imports
    if grep -q '^import "' "$file" 2>/dev/null; then
        single_imports=$(grep '^import "' "$file" 2>/dev/null | sed 's/^import "//; s/".*//' 2>/dev/null || true)
    fi

    local imports=""
    if [[ -n "$block_imports" ]]; then
        imports="$block_imports"
    fi
    if [[ -n "$single_imports" ]]; then
        if [[ -n "$imports" ]]; then
            imports="$imports
$single_imports"
        else
            imports="$single_imports"
        fi
    fi

    if [[ -z "$imports" ]]; then
        echo "⏭️  No imports to check: $file"
        ((skipped_count++))
        return
    fi

    echo "Checking: $file"
    ((checked_count++))

    # Check for proper grouping
    local has_system=false
    local has_external=false
    local has_alias_imports=false
    local has_internal=false
    local violations=()

    # Parse imports line by line
    local in_system=true
    local in_external=false
    local in_alias=false
    local in_internal=false
    local prev_blank=false

    while IFS= read -r line; do
        # Skip empty lines and comments
        [[ -z "$line" ]] && continue
        [[ "$line" =~ ^[[:space:]]*// ]] && continue

        # Check import type
        # Extract the import path and check for alias
        import_line=$(echo "$line" | sed 's/^[[:space:]]*//; s/[[:space:]]*$//')
        import_path=$(echo "$line" | sed 's/^[[:space:]]*[^"]*"\([^"]*\)".*/\1/')
        has_alias=false
        if [[ "$import_line" =~ ^[a-zA-Z_][a-zA-Z0-9_]*\  ]]; then
            has_alias=true
        fi

        if $has_alias; then
            # Alias import
            if $has_internal; then
                violations+=("Aliased import '$line' found after internal imports")
            fi
            has_alias_imports=true
            in_system=false
            in_external=false
            in_alias=true
        elif [[ "$import_path" =~ ^eda-in-golang ]]; then
            # Internal lib (project modules)
            # Internal imports can come after aliases - this is correct
            has_internal=true
            in_system=false
            in_external=false
            in_alias=false
        elif [[ "$import_path" =~ \. ]] && [[ ! "$import_path" =~ ^eda-in-golang ]]; then
            # External lib (contains dots and not internal)
            if $has_alias_imports || $has_internal; then
                violations+=("External import '$line' found after alias/internal imports")
            fi
            has_external=true
            in_system=false
            in_alias=false
            in_alias=false
        else
            # System lib (no dots, standard library)
            if $has_external || $has_alias_imports || $has_internal; then
                violations+=("System import '$line' found after external/alias/internal imports")
            fi
            has_system=true
        fi
    done <<< "$imports"

    if [[ ${#violations[@]} -gt 0 ]]; then
        echo "❌ Import violations in $file:"
        for violation in "${violations[@]}"; do
            echo "  - $violation"
        done
        ((error_count++))
        return 1
    else
        echo "✅ $file"
    fi
}

# Check all .go files
if [[ -f "$TARGET" && "$TARGET" == *.go ]]; then
    # Single file check
    check_file "$TARGET"
else
    # Directory check - process files using a for loop
    # This avoids subshell issues with process substitution
    for file in $(find "$TARGET" -name "*.go" -type f 2>/dev/null || true); do
        [[ -z "$file" ]] && continue  # Skip empty entries
        if ! check_file "$file"; then
            # check_file returned 1 (import violations), but we continue processing
            true
        fi
    done
fi

# Display summary
echo ""
echo "📊 Import Check Summary:"
echo "  ✅ Checked files: $checked_count"
echo "  ⏭️  Skipped files: $skipped_count"
echo "  ❌ Files with errors: $error_count"
echo "  📁 Total files processed: $((checked_count + skipped_count))"

if [[ $error_count -gt 0 ]]; then
    echo ""
    echo "💡 To fix import issues automatically, run: make fix-imports"
    exit 1
else
    echo ""
    echo "🎉 All import checks passed!"
    exit 0
fi
