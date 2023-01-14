package readers

import (
	"encoding/binary"
	"errors"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

func FileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

// ParseUint32s parses a slice of strings into a slice of uint32s.
func ParseUint32s(ss []string) ([]uint32, error) {
	us := make([]uint32, 0, len(ss))
	for _, s := range ss {
		u, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return nil, err
		}

		us = append(us, uint32(u))
	}

	return us, nil
}

// ParseUint64s parses a slice of strings into a slice of uint64s.
func ParseUint64s(ss []string) ([]uint64, error) {
	us := make([]uint64, 0, len(ss))
	for _, s := range ss {
		u, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, err
		}

		us = append(us, u)
	}

	return us, nil
}

// ParsePInt64s parses a slice of strings into a slice of int64 pointers.
func ParsePInt64s(ss []string) ([]*int64, error) {
	us := make([]*int64, 0, len(ss))
	for _, s := range ss {
		u, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}

		us = append(us, &u)
	}

	return us, nil
}

// ReadUintFromFile reads a file and attempts to parse a uint64 from it.
func ReadUintFromFile(path string) (uint64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
}

// ReadIntFromFile reads a file and attempts to parse a int64 from it.
func ReadIntFromFile(path string) (int64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
}

// ReadStringFromFile reads a file and attempts to trim a string from it.
func ReadStringFromFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// ParseBool parses a string into a boolean pointer.
func ParseBool(b string) *bool {
	var result bool
	switch b {
	case "enabled":
		result = true
	case "disabled":
		result = false
	default:
		return nil
	}
	return &result
}

func GetEndianness() (binary.ByteOrder, error) {
	buffer := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buffer[0])) = uint16(0xABCD)

	switch buffer {
	case [2]byte{0xCD, 0xAB}:
		return binary.LittleEndian, nil
	case [2]byte{0xAB, 0xCD}:
		return binary.BigEndian, nil
	default:
		return nil, errors.New("could not determine native endianness")
	}
}

//const (
//	CPU_UNKNOWN_MODEL    = -1
//	CPU_SANDYBRIDGE      = 42
//	CPU_SANDYBRIDGE_EP   = 45
//	CPU_IVYBRIDGE        = 58
//	CPU_IVYBRIDGE_EP     = 62
//	CPU_HASWELL          = 60
//	CPU_HASWELL_ULT      = 69
//	CPU_HASWELL_GT3E     = 70
//	CPU_HASWELL_EP       = 63
//	CPU_BROADWELL        = 61
//	CPU_BROADWELL_GT3E   = 71
//	CPU_BROADWELL_EP     = 79
//	CPU_BROADWELL_DE     = 86
//	CPU_SKYLAKE          = 78
//	CPU_SKYLAKE_HS       = 94
//	CPU_SKYLAKE_X        = 85
//	CPU_KNIGHTS_LANDING  = 87
//	CPU_KNIGHTS_MILL     = 133
//	CPU_KABYLAKE_MOBILE  = 142
//	CPU_KABYLAKE         = 158
//	CPU_ATOM_SILVERMONT  = 55
//	CPU_ATOM_AIRMONT     = 76
//	CPU_ATOM_MERRIFIELD  = 74
//	CPU_ATOM_MOOREFIELD  = 90
//	CPU_ATOM_GOLDMONT    = 92
//	CPU_ATOM_GEMINI_LAKE = 122
//	CPU_ATOM_DENVERTON   = 95
//)
