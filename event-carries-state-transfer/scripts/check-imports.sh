#!/bin/bash
# check-imports.sh - Check Go import organization rules
# Usage: ./scripts/check-imports.sh [file|directory]

set -e

TARGET="${1:-.}"

echo "Checking import organization in: $TARGET"

# Function to check imports in a single file
check_file() {
    local file="$1"
    if [[ ! -f "$file" ]] || [[ "$file" != *.go ]]; then
        return
    fi

    # Skip generated files
    if [[ "$file" =~ \.pb\.go$ ]] || [[ "$file" =~ _grpc\.pb\.go$ ]] || [[ "$file" =~ _pb\.go$ ]]; then
        echo "⏭️  Skipping generated file: $file"
        return
    fi

    # Skip files with generated comments
    if head -5 "$file" | grep -q "Code generated\|DO NOT EDIT\|This file was auto-generated"; then
        echo "⏭️  Skipping generated file: $file"
        return
    fi

    # Extract import block (remove import( and ) lines, keep only the actual imports)
    local imports=$(sed -n '/^import (/,/)/p' "$file" | sed '1d;$d' | grep -v '^$')

    if [[ -z "$imports" ]]; then
        echo "⏭️  No imports to check: $file"
        return
    fi

    echo "Checking: $file"

    # Check for proper grouping
    local has_system=false
    local has_external=false
    local has_internal=false
    local violations=()

    # Parse imports line by line
    local in_system=true
    local in_external=false
    local in_internal=false
    local prev_blank=false

    while IFS= read -r line; do
        # Skip empty lines and comments
        [[ -z "$line" ]] && continue
        [[ "$line" =~ ^[[:space:]]*// ]] && continue

        # Check import type
        # Extract the import path (remove quotes and whitespace)
        import_path=$(echo "$line" | sed 's/^[[:space:]]*[^"]*"\([^"]*\)".*/\1/')

        if [[ "$import_path" =~ ^eda-in-golang ]]; then
            # Internal lib (project modules)
            has_internal=true
            in_system=false
            in_external=false
        elif [[ "$import_path" =~ \. ]] && [[ ! "$import_path" =~ ^eda-in-golang ]]; then
            # External lib (contains dots and not internal)
            if $has_internal; then
                violations+=("External import '$line' found after internal imports")
            fi
            has_external=true
            in_system=false
        else
            # System lib (no dots, standard library)
            if $has_external || $has_internal; then
                violations+=("System import '$line' found after external/internal imports")
            fi
            has_system=true
        fi
    done <<< "$imports"

    if [[ ${#violations[@]} -gt 0 ]]; then
        echo "❌ Import violations in $file:"
        for violation in "${violations[@]}"; do
            echo "  - $violation"
        done
        return 1
    else
        echo "✅ $file"
    fi
}

# Check all .go files
find "$TARGET" -name "*.go" -type f | while read -r file; do
    check_file "$file" || exit 1
done

echo "Import check completed!"
