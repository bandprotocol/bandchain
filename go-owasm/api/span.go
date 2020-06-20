package api

// #include "bindings.h"
// #include <string.h>
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

func writeSpan(span *C.Span, data []byte) C.GoResult {
	if int(span.cap) < len(data) {
		return C.GoResult_SpanExceededCapacity
	}
	C.memcpy(unsafe.Pointer(span.ptr), unsafe.Pointer(&data[0]), C.size_t(len(data)))
	span.len = C.uintptr_t(len(data))
	return C.GoResult_Ok
}
