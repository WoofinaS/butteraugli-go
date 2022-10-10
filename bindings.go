package butteraugli_go

// #cgo CFLAGS:
// #cgo LDFLAGS: -ljxl -shared
// #include <jxl/butteraugli.h>
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
	TYPE_FLOAT   DATATYPE = 1
	TYPE_UINT8   DATATYPE = 2
	TYPE_UINT16  DATATYPE = 3
	TYPE_FLOAT16 DATATYPE = 4
)

type ENDIANNESS int

const (
	NATIVE_ENDIAN ENDIANNESS = 1
	LITTLE_ENDIAN ENDIANNESS = 2
	BIG_ENDIAN    ENDIANNESS = 3
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
