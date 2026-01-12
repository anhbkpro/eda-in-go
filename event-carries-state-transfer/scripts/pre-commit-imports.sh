#!/bin/bash
# pre-commit-imports.sh - Pre-commit hook to check import organization
#
# This hook runs before each commit to ensure staged Go files have
# properly organized imports. It uses goimports to check for any
# formatting issues that would change the file.
#
# For faster checking, consider using the Go implementation instead:
# make check-imports-go
#
# To set up: cp scripts/pre-commit-imports.sh .git/hooks/pre-commit
#            chmod +x .git/hooks/pre-commit

echo "🔍 Checking import organization..."

# Check if goimports is available
if ! command -v goimports &> /dev/null; then
    echo "⚠️  goimports not found. Install with: go install golang.org/x/tools/cmd/goimports@latest"
    echo "Skipping import check."
    exit 0
fi

# Get staged Go files
staged_files=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$')

if [[ -z "$staged_files" ]]; then
    echo "✅ No Go files staged for commit"
    exit 0
fi

echo "Checking imports in staged files..."
has_violations=false

for file in $staged_files; do
    if [[ ! -f "$file" ]]; then
        continue
    fi

    # Skip generated files
    if [[ "$file" =~ \.pb\.go$ ]] || [[ "$file" =~ _grpc\.pb\.go$ ]] || [[ "$file" =~ _pb\.go$ ]]; then
        echo "⏭️  Skipping generated file: $file"
        continue
    fi

    # Skip files with generated comments
    if head -5 "$file" | grep -q "Code generated\|DO NOT EDIT\|This file was auto-generated"; then
        echo "⏭️  Skipping generated file: $file"
        continue
    fi

    # Check if goimports would make changes
    if ! goimports -d -local eda-in-golang "$file" | grep -q .; then
        echo "✅ $file"
    else
        echo "❌ $file has import organization issues"
        echo "Run: goimports -w -local eda-in-golang $file"
        has_violations=true
    fi
done

if [[ "$has_violations" = true ]]; then
    echo ""
    echo "🚫 Import organization violations found!"
    echo "Fix with: make check-imports-go"
    echo "Or auto-fix with: goimports -w -local eda-in-golang ."
    exit 1
fi

echo "✅ All staged Go files have properly organized imports"
exit 0
