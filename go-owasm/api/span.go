package api

// #include "bindings.h"
import "C"
import "unsafe"

func newSpan(size int) C.Span {
	return C.Span{
		ptr: (*C.uint8_t)(C.malloc(C.uintptr_t(size))),
		len: 0,
		cap: C.uintptr_t(size),
	}
}

func copySpan(data []byte) C.Span {
	return C.Span{
		ptr: (*C.uint8_t)(C.CBytes(data)),
		len: C.uintptr_t(len(data)),
		cap: C.uintptr_t(len(data)),
	}
}

func readSpan(span C.Span) []byte {
	return C.GoBytes(unsafe.Pointer(span.ptr), C.int(span.len))
}

func freeSpan(span C.Span) {
	C.free(unsafe.Pointer(span.ptr))
}
