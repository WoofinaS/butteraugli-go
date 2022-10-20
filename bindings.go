package butteraugli_go

// #cgo CFLAGS:
// #cgo LDFLAGS: -ljxl -shared
// #include <jxl/butteraugli.h>
// #include <jxl/types.h>
import "C"

// NUM_CHANNELS represents the number of channels/subpixels per pixel.
type NUM_CHANNELS int

const (
	GRAYSCALE  NUM_CHANNELS = 1
	GRAY_ALPHA NUM_CHANNELS = 2
	RGB        NUM_CHANNELS = 3
	RGBA       NUM_CHANNELS = 4
)

// DATATYPE represets the type for every pixel/sub pixel of a raw frame.
type DATATYPE int

const (
	TYPE_FLOAT   DATATYPE = C.JXL_TYPE_FLOAT
	TYPE_UINT8   DATATYPE = C.JXL_TYPE_UINT8
	TYPE_UINT16  DATATYPE = C.JXL_TYPE_UINT16
	TYPE_FLOAT16 DATATYPE = C.JXL_TYPE_FLOAT16
)

// ENDIANNESS represets the byte level ordering of the raw frame.
type ENDIANNESS int

const (
	NATIVE_ENDIAN ENDIANNESS = C.JXL_NATIVE_ENDIAN
	LITTLE_ENDIAN ENDIANNESS = C.JXL_LITTLE_ENDIAN
	BIG_ENDIAN    ENDIANNESS = C.JXL_BIG_ENDIAN
)

// API is a wrapper struct for interacting with the butteraugli api.
type API struct {
	freed  bool // keeps track of if the jxlAPI has been freed yet.
	jxlAPI *C.JxlButteraugliApi
}

// PixelFormat represents the pixel format properties of a frame.
type PixelFormat struct {
	NumChannels NUM_CHANNELS
	DataType    DATATYPE
	Endianness  ENDIANNESS
	Align       uint32
}

// Result is a wrapper struct for interacting with butteraugli results.
type Result struct {
	freed     bool // keeps track of if the jxlResult has been freed yet
	jxlResult *C.JxlButteraugliResult
}

// ComputeTask represents all the inputs used for the Compute function. Width
// and Height represent the width and height of the source and distored frame.
// RefBytes and DisBytes represents the raw bytes of both frames going from the
// top left of the frame to the bottom right. RefPixFmt and DisPixFmt represent
// the pixel properties of both frames.
type ComputeTask struct {
	// Must be the same for both frames
	Width, Height uint32

	// the raw image must first be converted into a linear format before being
	// passed. If this is not done then the results can vary from slightly to
	// significantly inaccurate.
	RefBytes, DisBytes []byte

	// If the pixel format is incorrect Compute will return a error.
	RefPixFmt, DisPixFmt PixelFormat
}
