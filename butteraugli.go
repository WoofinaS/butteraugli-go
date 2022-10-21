package butteraugli_go

// #cgo CFLAGS:
// #cgo LDFLAGS: -ljxl -shared
// #include <stdint.h>
// #include <jxl/butteraugli.h>
import "C"
import (
	"errors"
	"reflect"
	"runtime"
	"unsafe"
)

// ApiCreate creates a new API structure that is used to interact with the
// butteraugli api.
func ApiCreate() API {
	a := API{false, C.JxlButteraugliApiCreate(nil)}

	// frees underlying butteraugli api if the user didn't when the struct is
	// garbage collected. This avoids a potential memory leak.
	runtime.SetFinalizer(&a, func(a *API) {
		a.Destroy()
	})

	return a
}

// Destroy destroys the underlying butteraugli api and frees memory. This is
// implicitly called when the Result is garbage collected. However when the
// Result is garbage collected is undeterminable.
func (a *API) Destroy() {
	if a.freed {
		return
	}
	C.JxlButteraugliApiDestroy(a.jxlAPI)
	a.freed = true
}

// SetIntensityTarget sets the butteraugli Intensity Target. This should be set
// to the target contents or display max brightness.
func (a *API) SetIntensityTarget(intensity float32) {
	C.JxlButteraugliApiSetIntensityTarget(a.jxlAPI, C.float(intensity))
}

// SetHFAsymmetry controls bias for penalizing high frequency artifacts over
// blurring. 1.0 = neutral / default
func (a *API) SetHFAsymmetry(asymmetry float32) {
	C.JxlButteraugliApiSetHFAsymmetry(a.jxlAPI, C.float(asymmetry))
}

// Compute takes a ComputeTask as a input and returns a Result struct and error
// If there is a error the user can safely assume that the Result does not need
// to be destroyed.
func (a *API) Compute(t ComputeTask) (Result, error) {
	refPixFmt := C.JxlPixelFormat{
		C.uint32_t(t.RefPixFmt.NumChannels),
		C.JxlDataType(t.RefPixFmt.DataType),
		C.JxlEndianness(t.RefPixFmt.Endianness),
		C.ulong(t.RefPixFmt.Align),
	}
	disPixFmt := C.JxlPixelFormat{
		C.uint32_t(t.DisPixFmt.NumChannels),
		C.JxlDataType(t.DisPixFmt.DataType),
		C.JxlEndianness(t.DisPixFmt.Endianness),
		C.ulong(t.DisPixFmt.Align),
	}

	// We can not pass go slices to c as they contain go pointers as well as
	// being as incompatible data type. Therefore we must get a unsafe pointer
	// to the underlying array if we want to avoid making a copy in c allocated
	// memory. This also means we can avoid importing stdlib.h for malloc/free.

	// Creates a unsafe pointer to both byte slices
	refHeader := (*reflect.SliceHeader)(unsafe.Pointer(&t.RefBytes))
	disHeader := (*reflect.SliceHeader)(unsafe.Pointer(&t.DisBytes))

	// Gets unsafe pointer to underlying data for both slices
	refData := unsafe.Pointer(refHeader.Data)
	disData := unsafe.Pointer(disHeader.Data)

	result := C.JxlButteraugliCompute(a.jxlAPI, C.uint32_t(t.Height),
		C.uint32_t(t.Width), &refPixFmt, refData, C.ulong(len(t.RefBytes)),
		&disPixFmt, disData, C.ulong(len(t.DisBytes)))

	// Prevents slices from being garbage collected while in use. GC does not
	// keep track of unsafe pointers. Removing this means the GC might free the
	// byte slices during Butteraugli calculations.
	_ = t.RefBytes[0]
	_ = t.DisBytes[0]

	if result == nil {
		return Result{}, errors.New("failed to compute butteraugli scores")
	}

	r := Result{false, result}

	// frees underlying butteraugli result if the user didn't when the struct
	// is garbage collected. This avoids a potential memory leak.
	runtime.SetFinalizer(&r, func(r *Result) {
		r.Destroy()
	})

	return r, nil
}

// Destroy destroys the underlying butteraugli result and frees memory. This
// is implicitly called when the Result is garbage collected. However when the
// Result is garbage collected is undeterminable.
func (r *Result) Destroy() {
	if r.freed {
		return
	}
	C.JxlButteraugliResultDestroy(r.jxlResult)
	r.freed = true
}

// GetMaxDistance returns the highest distance of the result.
func (r *Result) GetMaxDistance() float32 {
	return float32(C.JxlButteraugliResultGetMaxDistance(r.jxlResult))
}

// GetDistance returns the average butteraugli distance of each pixel
// averaged by the given pnorm. More information about pnorm can be seen here.
// https://en.wikipedia.org/wiki/Norm_(mathematics)#p-norm
func (r *Result) GetDistance(pnorm float32) float32 {
	return float32(C.JxlButteraugliResultGetDistance(r.jxlResult,
		C.float(pnorm)))
}
