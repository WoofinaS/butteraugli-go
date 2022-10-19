package butteraugli_go

// #cgo CFLAGS:
// #cgo LDFLAGS: -ljxl -shared
// #include <jxl/butteraugli.h>
// #include <jxl/types.h>
import "C"

type NUM_CHANNELS int

const (
	GRAYSCALE  NUM_CHANNELS = 1
	GRAY_ALPHA NUM_CHANNELS = 2
	RGB        NUM_CHANNELS = 3
	RGBA       NUM_CHANNELS = 4
)

type DATATYPE int

const (
	TYPE_FLOAT   DATATYPE = C.JXL_TYPE_FLOAT
	TYPE_UINT8   DATATYPE = C.JXL_TYPE_UINT8
	TYPE_UINT16  DATATYPE = C.JXL_TYPE_UINT16
	TYPE_FLOAT16 DATATYPE = C.JXL_TYPE_FLOAT16
)

type ENDIANNESS int

const (
	NATIVE_ENDIAN ENDIANNESS = C.JXL_NATIVE_ENDIAN
	LITTLE_ENDIAN ENDIANNESS = C.JXL_LITTLE_ENDIAN
	BIG_ENDIAN    ENDIANNESS = C.JXL_BIG_ENDIAN
)

// API is a simple wrapper struct for the butteraugli api.
type API struct {
	jxlAPI *C.JxlButteraugliApi
}

// PixelFormat is a simple wrapper struct for butteraugli JxlPixelFormat.
type PixelFormat struct {
	NumChannels NUM_CHANNELS
	DataType    DATATYPE
	Endianness  ENDIANNESS
	Align       uint32
}

// Result is a simple wrapper struct for results from the butteraugli api.
type Result struct {
	jxlResult *C.JxlButteraugliResult
}

type ComputeTask struct {
	Width, Height        uint32
	RefBytes, DisBytes   []byte
	RefPixFmt, DisPixFmt PixelFormat
}
