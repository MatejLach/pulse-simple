package pulse

/*
#cgo pkg-config: libpulse-simple

#include <stdlib.h>
#include <pulse/sample.h>
*/
import "C"
import "unsafe"

const (
	CHANNELS_MAX = C.PA_CHANNELS_MAX // Maximum number of allowed channels
	RATE_MAX     = C.PA_RATE_MAX     // Maximum allowed sample rate
)

type SampleFormat C.pa_sample_format_t

const (
	SAMPLE_U8        SampleFormat = C.PA_SAMPLE_U8
	SAMPLE_ALAW      SampleFormat = C.PA_SAMPLE_ALAW
	SAMPLE_ULAW      SampleFormat = C.PA_SAMPLE_ULAW
	SAMPLE_S16LE     SampleFormat = C.PA_SAMPLE_S16LE
	SAMPLE_S16BE     SampleFormat = C.PA_SAMPLE_S16BE
	SAMPLE_FLOAT32LE SampleFormat = C.PA_SAMPLE_FLOAT32LE
	SAMPLE_FLOAT32BE SampleFormat = C.PA_SAMPLE_FLOAT32BE
	SAMPLE_S32LE     SampleFormat = C.PA_SAMPLE_S32LE
	SAMPLE_S32BE     SampleFormat = C.PA_SAMPLE_S32BE
	SAMPLE_S24LE     SampleFormat = C.PA_SAMPLE_S24LE
	SAMPLE_S24BE     SampleFormat = C.PA_SAMPLE_S24BE
	SAMPLE_S24_32LE  SampleFormat = C.PA_SAMPLE_S24_32LE
	SAMPLE_S24_32BE  SampleFormat = C.PA_SAMPLE_S24_32BE
	SAMPLE_MAX       SampleFormat = C.PA_SAMPLE_MAX
	SAMPLE_INVALID   SampleFormat = C.PA_SAMPLE_INVALID
)

const (
	SAMPLE_S16NE     SampleFormat = C.PA_SAMPLE_S16NE
	SAMPLE_FLOAT32NE SampleFormat = C.PA_SAMPLE_FLOAT32NE
	SAMPLE_S32NE     SampleFormat = C.PA_SAMPLE_S32NE
	SAMPLE_S24NE     SampleFormat = C.PA_SAMPLE_S24NE
	SAMPLE_S24_32NE  SampleFormat = C.PA_SAMPLE_S24_32NE
	SAMPLE_S16RE     SampleFormat = C.PA_SAMPLE_S16RE
	SAMPLE_FLOAT32RE SampleFormat = C.PA_SAMPLE_FLOAT32RE
	SAMPLE_S32RE     SampleFormat = C.PA_SAMPLE_S32RE
	SAMPLE_S24RE     SampleFormat = C.PA_SAMPLE_S24RE
	SAMPLE_S24_32RE  SampleFormat = C.PA_SAMPLE_S24_32RE
)

const SAMPLE_FLOAT32 SampleFormat = C.PA_SAMPLE_FLOAT32

type SampleSpec struct {
	Format   SampleFormat
	Rate     uint32
	Channels uint8
}

func cSampleSpec(spec *SampleSpec) *C.pa_sample_spec {
	return &C.pa_sample_spec{
		format:   C.pa_sample_format_t(spec.Format),
		rate:     C.uint32_t(spec.Rate),
		channels: C.uint8_t(spec.Channels),
	}
}

// SampleSpec.BytesPerSecond returns the number of bytes per second of audio.
func (spec *SampleSpec) BytesPerSecond() uint {
	return uint(C.pa_bytes_per_second(cSampleSpec(spec)))
}

// SampleSpec.FrameSize returns the size of a single audio frame in bytes.
func (spec *SampleSpec) FrameSize() uint {
	return uint(C.pa_frame_size(cSampleSpec(spec)))
}

// SampleSpec.SampleSize returns the size of a single sample in bytes.
func (spec *SampleSpec) SampleSize() uint {
	return uint(C.pa_sample_size(cSampleSpec(spec)))
}

// SampleFormat.SampleSize returns the size of a single sample in bytes.
func (f SampleFormat) SampleSize() uint {
	return uint(C.pa_sample_size_of_format(C.pa_sample_format_t(f)))
}

// SampleSpec.BytesToUsec returns the number of microseconds taken
// to play the given number of bytes as audio.
//
// The return value will always be rounded down for non-integral values.
func (spec *SampleSpec) BytesToUsec(bytes uint) uint64 {
	return uint64(C.pa_bytes_to_usec(C.uint64_t(bytes), cSampleSpec(spec)))
}

// SampleSpec.UsecToBytes returns the number of bytes required
// for the given number of microseconds of audio.
//
// The return value will always be rounded down for non-integral values.
func (spec *SampleSpec) UsecToBytes(usec uint64) uint {
	return uint(C.pa_usec_to_bytes(C.pa_usec_t(usec), cSampleSpec(spec)))
}

/*
// Initialize the specified sample spec and return a pointer to
// it. The sample spec will have a defined state but
// pa_sample_spec_valid() will fail for it. \since 0.9.13
pa_sample_spec* pa_sample_spec_init(pa_sample_spec *spec);
*/

// SampleSpec.Valid returns whether or not the given sample spec is valid.
func (spec *SampleSpec) Valid() bool {
	if C.pa_sample_spec_valid(cSampleSpec(spec)) == 0 {
		return false
	}
	return true
}

// SampleSpec.Equal returns whether or not the given sample specs match.
func (spec *SampleSpec) Equal(other *SampleSpec) bool {
	if C.pa_sample_spec_equal(cSampleSpec(spec), cSampleSpec(other)) == 0 {
		return false
	}
	return true
}

// SampleFormat.String returns a string describing the format.
func (f SampleFormat) String() string {
	cstr := C.pa_sample_format_to_string(C.pa_sample_format_t(f))
	ret := C.GoString(cstr)
	return ret
}

// ParseSampleFormat returns the SampleFormat described by the given string.
//
// The string should be as returned by SampleFormat.String().
func ParseSampleFormat(s string) SampleFormat {
    cstr := C.CString(s)
    defer C.free(unsafe.Pointer(cstr))
    ret := SampleFormat(C.pa_parse_sample_format(cstr))
    return ret
}

/*
// Maximum required string length for
// pa_sample_spec_snprint(). Please note that this value can change
// with any release without warning and without being considered API
// or ABI breakage. You should not use this definition anywhere where
// it might become part of an ABI.
#define PA_SAMPLE_SPEC_SNPRINT_MAX 32

// Pretty print a sample type specification to a string
char* pa_sample_spec_snprint(char *s, size_t l, const pa_sample_spec *spec);

// Maximum required string length for pa_bytes_snprint(). Please note
// that this value can change with any release without warning and
// without being considered API or ABI breakage. You should not use
// this definition anywhere where it might become part of an
// ABI. \since 0.9.16
#define PA_BYTES_SNPRINT_MAX 11

// Pretty print a byte size value (i.e.\ "2.5 MiB")
char* pa_bytes_snprint(char *s, size_t l, unsigned v);
*/

// SampleFormat.IsLe returns 1 when the format is little endian.
//
// Returns -1 when endianness does not apply to this format.
func (f SampleFormat) IsLe() int {
    return int(C.pa_sample_format_is_le(C.pa_sample_format_t(f)))
}

// SampleFormat.IsBe returns 1 when the format is big endian.
//
// Returns -1 when endianness does not apply to this format.
func (f SampleFormat) IsBe() int {
    return int(C.pa_sample_format_is_be(C.pa_sample_format_t(f)))
}

// SampleFormat.IsNe returns 1 when the format is native endian.
//
// Returns -1 when endianness does not apply to this format.
func (f SampleFormat) IsNe() int {
    // note: C.pa_sample_format_is_ne() doesn't seem to work
    if SAMPLE_S16NE == SAMPLE_S16LE {
        return int(C.pa_sample_format_is_le(C.pa_sample_format_t(f)))
    }
    return int(C.pa_sample_format_is_be(C.pa_sample_format_t(f)))
}

// SampleFormat.IsRe returns 1 when the format is reverse endian.
//
// Returns -1 when endianness does not apply to this format.
func (f SampleFormat) IsRe() int {
    // note: C.pa_sample_format_is_re() doesn't seem to work
    if SAMPLE_S16NE == SAMPLE_S16LE {
        return int(C.pa_sample_format_is_be(C.pa_sample_format_t(f)))
    }
    return int(C.pa_sample_format_is_le(C.pa_sample_format_t(f)))
}

