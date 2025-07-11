package main

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/compute/exec"
	"github.com/apache/arrow-go/v18/arrow/scalar"
)

func main() {
	fmt.Printf("Architecture: %s\n", runtime.GOARCH)
	fmt.Printf("Endianness: %s\n", getEndianness())
	
	// Test the union scalar fix
	testUnionScalar()
}

func getEndianness() string {
	var i int32 = 0x01020304
	u := unsafe.Pointer(&i)
	pb := (*byte)(u)
	b := *pb
	if b == 0x04 {
		return "little"
	}
	return "big"
}

func testUnionScalar() {
	fmt.Println("\nTesting union scalar fix...")
	
	// Create a dense union scalar
	dt := arrow.UnionOf(arrow.DenseMode, []arrow.Field{
		{Name: "string", Type: arrow.BinaryTypes.String, Nullable: true},
		{Name: "number", Type: arrow.PrimitiveTypes.Uint64, Nullable: true},
		{Name: "other_number", Type: arrow.PrimitiveTypes.Uint64, Nullable: true},
	}, []arrow.UnionTypeCode{3, 42, 43})
	
	unionScalar := scalar.NewDenseUnionScalar(scalar.MakeScalar(uint64(25)), 42, dt.(*arrow.DenseUnionType))
	
	// Create ArraySpan and fill from scalar
	var span exec.ArraySpan
	span.FillFromScalar(unionScalar)
	
	fmt.Printf("Dense Union Scratch[0]: %d (expected: 42)\n", span.Scratch[0])
	fmt.Printf("Dense Union Scratch[1]: %d (expected: 1)\n", span.Scratch[1])
	
	// Create a sparse union scalar
	dt2 := arrow.UnionOf(arrow.SparseMode, []arrow.Field{
		{Name: "string", Type: arrow.BinaryTypes.String, Nullable: true},
		{Name: "number", Type: arrow.PrimitiveTypes.Uint64, Nullable: true},
		{Name: "other_number", Type: arrow.PrimitiveTypes.Uint64, Nullable: true},
	}, []arrow.UnionTypeCode{3, 42, 43})
	
	sparseScalar := scalar.NewSparseUnionScalarFromValue(scalar.MakeScalar(uint64(25)), 1, dt2.(*arrow.SparseUnionType))
	
	var span2 exec.ArraySpan
	span2.FillFromScalar(sparseScalar)
	
	fmt.Printf("Sparse Union Scratch[0]: %d (expected: 42)\n", span2.Scratch[0])
	fmt.Printf("Sparse Union Scratch[1]: %d (expected: 0)\n", span2.Scratch[1])
}