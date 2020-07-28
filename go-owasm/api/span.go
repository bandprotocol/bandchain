package api

// #include "bindings.h"
// #include <string.h>
import "C"
import "unsafe"

// newSpan creates a span with the given capacity and zero length, ready to be written.
func newSpan(size int) C.Span {
	return C.Span{
		ptr: (*C.uint8_t)(C.malloc(C.uintptr_t(size))),
		len: 0,
		cap: C.uintptr_t(size),
	}
}

// copySpan creates a full span populated with the given data.
func copySpan(data []byte) C.Span {
	return C.Span{
		ptr: (*C.uint8_t)(C.CBytes(data)),
		len: C.uintptr_t(len(data)),
		cap: C.uintptr_t(len(data)),
	}
}

// readSpan returns a copy of data in the span.
func readSpan(span C.Span) []byte {
	return C.GoBytes(unsafe.Pointer(span.ptr), C.int(span.len))
}

// freeSpan deallocates the memory data of the span. Must be called for every newSpan/copySpan.
func freeSpan(span C.Span) {
	C.free(unsafe.Pointer(span.ptr))
}

// writeSpan writes the given data into the span. Returns error if capacity is not enough.
func writeSpan(span *C.Span, data []byte) C.Error {
	if int(span.cap) < len(data) {
		return C.Error_SpanTooSmallError
	}
	if len(data) > 0 {
		C.memcpy(unsafe.Pointer(span.ptr), unsafe.Pointer(&data[0]), C.size_t(len(data)))
	}
	span.len = C.uintptr_t(len(data))
	return C.Error_NoError
}
