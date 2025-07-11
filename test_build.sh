#!/bin/bash

echo "Testing Arrow Go fixes for s390x endianness issues..."

# Test the union scalar fix
echo "Building test program..."
go build -o test_fixes test_fixes.go

if [ $? -eq 0 ]; then
    echo "Build successful, running test..."
    ./test_fixes
    rm -f test_fixes
else
    echo "Build failed"
    exit 1
fi

# Test specific packages that were failing
echo -e "\nTesting compute/exec package..."
go test -v ./arrow/compute/exec -run TestArraySpan_FillFromScalar

echo -e "\nTesting avro package..."
go test -v ./arrow/avro -run TestReader

echo "Done!"