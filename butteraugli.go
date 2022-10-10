package butteraugli_go

// #cgo CFLAGS: -O3 -march=native
// #cgo LDFLAGS: -ljxl -shared
// #include <stdint.h>
// #include <stdlib.h>
// #include <jxl/butteraugli.h>
import "C"
import (
	"errors"
	"image"
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
func (a *API) SetIntensityTarget(intensity int) {
	C.JxlButteraugliApiSetIntensityTarget(a.jxlAPI, C.float(intensity))
}

// SetHFAsymmetry sets the butteraugli Asymmetry.
func (a *API) SetHFAsymmetry(asymmetry float32) {
	C.JxlButteraugliApiSetHFAsymmetry(a.jxlAPI, C.float(asymmetry))
}

// Compute takes in two image.Image frames and computes the butteraugli scores
// and returns a butteraugli result that is wraped in a  Result.
func (a *API) Compute(ref, dis image.Image) (Result, error) {
	rwidth, rheight, chans := ref.Bounds().Dx(), ref.Bounds().Dy(), 4
	if rwidth != dis.Bounds().Dx() || rheight != dis.Bounds().Dy() {
		return Result{}, errors.New("ref/dist sizes dont match")
	}
	numBytes := rwidth * rheight * chans * int(unsafe.Sizeof(uint16(0)))
	lenBytes := C.ulong(rwidth * rheight * chans)
	refPoint := C.malloc(C.ulong(numBytes))
	disPoint := C.malloc(C.ulong(numBytes))
	refBytes := unsafe.Slice((*uint16)(refPoint), lenBytes)
	disBytes := unsafe.Slice((*uint16)(disPoint), lenBytes)

	index := 0
	for x := 0; x < rwidth; x++ {
		for y := 0; y < rheight; y++ {
			rr, rg, rb, ra := ref.At(x, y).RGBA()
			dr, dg, db, da := dis.At(x, y).RGBA()
			refBytes[index+0] = uint16(rr)
			refBytes[index+1] = uint16(rg)
			refBytes[index+2] = uint16(rb)
			refBytes[index+3] = uint16(ra)
			disBytes[index+0] = uint16(dr)
			disBytes[index+1] = uint16(dg)
			disBytes[index+2] = uint16(db)
			disBytes[index+3] = uint16(da)
			index += 4
		}
	}

	pixelfmt := C.JxlPixelFormat{4, 3, 1, 0}

	result := C.JxlButteraugliCompute(a.jxlAPI, C.uint32_t(rheight), C.uint32_t(rwidth),
		&pixelfmt, refPoint, C.ulong(numBytes), &pixelfmt, disPoint, C.ulong(numBytes))
	if result == nil {
		return Result{}, errors.New("failed to compute result")
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
func (r *Result) GetDistance(pnorm int) float32 {
	return float32(C.JxlButteraugliResultGetDistance(r.jxlResult, C.float(pnorm)))
}
