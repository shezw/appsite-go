#!/bin/bash

# Usage: ./coverage.sh [module_name]
# Example: ./coverage.sh core
# Example: ./coverage.sh services
# Example: ./coverage.sh commerce (matches internal/services/commerce)

TARGET=$1

if [ -z "$TARGET" ]; then
    echo "No target specified. Running all tests by default."
    echo "Usage: ./coverage.sh [all|services|core|pkg|commerce|user|...]"
    
    # Run all tests
    go test -coverpkg=./internal/...,./pkg/... -coverprofile=coverage.out ./tests/...
    go tool cover -func=coverage.out
    exit 0
fi

# Define paths based on target
COVER_PKG=""
TEST_PATH=""

case $TARGET in
    "all")
        COVER_PKG="./internal/...,./pkg/..."
        TEST_PATH="./tests/..."
        ;;
    "services")
        COVER_PKG="./internal/services/..."
        TEST_PATH="./tests/internal/services/..."
        ;;
    "core")
        COVER_PKG="./internal/core/..."
        TEST_PATH="./tests/internal/core/..."
        ;;
    "pkg")
        COVER_PKG="./pkg/..."
        TEST_PATH="./tests/pkg/..."
        ;;
    "model")
        COVER_PKG="./internal/core/model/..."
        TEST_PATH="./tests/internal/core/model/..."
        ;;
    *)
        # Check if it's a specific service
        if [ -d "./internal/services/$TARGET" ]; then
            COVER_PKG="./internal/services/$TARGET/..."
            TEST_PATH="./tests/internal/services/$TARGET/..."
        # Check if it's a specific core module
        elif [ -d "./internal/core/$TARGET" ]; then
            COVER_PKG="./internal/core/$TARGET/..."
            TEST_PATH="./tests/internal/core/$TARGET/..."
        # Check if it's a specific pkg module
        elif [ -d "./pkg/$TARGET" ]; then
            COVER_PKG="./pkg/$TARGET/..."
            # Note: tests might be in tests/pkg/TARGET or just pkg/TARGET if it was a unit test inside pkg
            # But based on project structure, tests seem to be in tests/
            TEST_PATH="./tests/pkg/$TARGET/..."
        # Check if it's a specific pkg extra module
        elif [ -d "./pkg/extra/$TARGET" ]; then
            COVER_PKG="./pkg/extra/$TARGET/..."
            TEST_PATH="./tests/pkg/extra/$TARGET/..."
        else
            echo "Error: Unknown module directory for '$TARGET'"
            exit 1
        fi
        ;;
esac

echo "========================================"
echo "Target: $TARGET"
echo "Pkg:    $COVER_PKG"
echo "Tests:  $TEST_PATH"
echo "========================================"

# Run test
go test -v -coverpkg="$COVER_PKG" -coverprofile=coverage.out "$TEST_PATH"
RESULT=$?

if [ $RESULT -eq 0 ]; then
    echo "----------------------------------------"
    go tool cover -func=coverage.out
else
    echo "Tests Failed!"
    exit $RESULT
fi
