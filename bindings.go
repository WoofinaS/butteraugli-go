package butteraugli_go

// #cgo CFLAGS: -O3 -march=native
// #cgo LDFLAGS: -ljxl -shared
// #include <jxl/butteraugli.h>
import "C"

// API is a simple wrapper struct for the butteraugli api.
type API struct {
	jxlAPI *C.JxlButteraugliApi
}

// Result is a simple wrapper struct for results from the butteraugli api.
type Result struct {
	jxlResult *C.JxlButteraugliResult
}
