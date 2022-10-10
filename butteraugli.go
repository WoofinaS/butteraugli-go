package butteraugli_go

// #cgo CFLAGS:
// #cgo LDFLAGS: -ljxl -shared
// #include <stdint.h>
// #include <stdlib.h>
// #include <jxl/butteraugli.h>
import "C"
import (
	"errors"
	"unsafe"
)

// ApiCreate creates a new API structure that wraps around the butteraugli api.
func ApiCreate() API {
	return API{C.JxlButteraugliApiCreate(nil)}
}

// Destroy destorys the underlying butteraugli api and frees memory. This must
// be called when done using the api, else it creates a memory leak.
func (a *API) Destroy() {
	C.JxlButteraugliApiDestroy(a.jxlAPI)
}

// SetIntensityTarget sets the butteraugli Intensity Target. This should not be
// set to the brightness of the display and instead the larget luminance of the
// color space. For sRGB this is 80.
func (a *API) SetIntensityTarget(intensity float32) {
	C.JxlButteraugliApiSetIntensityTarget(a.jxlAPI, C.float(intensity))
}

// SetHFAsymmetry sets the butteraugli Asymmetry.
func (a *API) SetHFAsymmetry(asymmetry float32) {
	C.JxlButteraugliApiSetHFAsymmetry(a.jxlAPI, C.float(asymmetry))
}

// Compute takes a ComputeTask as a input and calculates the score of the image
// within it.
func (a *API) Compute_new(t ComputeTask) (Result, error) {
	refPixFmt := C.JxlPixelFormat{
		C.uint32_t(t.refPixFmt.NumChannels),
		C.JxlDataType(t.refPixFmt.DataType),
		C.JxlEndianness(t.refPixFmt.Endianness),
		C.ulong(t.refPixFmt.Align),
	}
	disPixFmt := C.JxlPixelFormat{
		C.uint32_t(t.disPixFmt.NumChannels),
		C.JxlDataType(t.disPixFmt.DataType),
		C.JxlEndianness(t.disPixFmt.Endianness),
		C.ulong(t.disPixFmt.Align),
	}

	refPoint := C.malloc(C.ulong(len(t.refBytes)))
	disPoint := C.malloc(C.ulong(len(t.disBytes)))
	refBytes := unsafe.Slice((*byte)(refPoint), len(t.refBytes))
	disBytes := unsafe.Slice((*byte)(disPoint), len(t.disBytes))

	copy(refBytes, t.refBytes)
	copy(disBytes, t.disBytes)

	result := C.JxlButteraugliCompute(a.jxlAPI, C.uint32_t(t.height),
		C.uint32_t(t.width), &refPixFmt, refPoint, C.ulong(len(t.refBytes)),
		&disPixFmt, disPoint, C.ulong(len(t.disBytes)))
	if result == nil {
		return Result{}, errors.New("failed to compute butteraugli scores")
	}

	C.free(refPoint)
	C.free(disPoint)

	return Result{result}, nil
}

// Destroy destorys the underlying butteraugli result and frees memory. This
// must be called when done using the result, else it creates a memory leak.
func (r *Result) Destroy() {
	C.JxlButteraugliResultDestroy(r.jxlResult)
}

// GetMaxDistance returns the max butteraugli score from the result.
func (r *Result) GetMaxDistance() float32 {
	return float32(C.JxlButteraugliResultGetMaxDistance(r.jxlResult))
}

// GetDistance returns the average butteraugli score from the result averaged
// by the given pnorm.
func (r *Result) GetDistance(pnorm float32) float32 {
	return float32(C.JxlButteraugliResultGetDistance(r.jxlResult, C.float(pnorm)))
}
