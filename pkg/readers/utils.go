package readers

import (
	"errors"
	"os"
	"strconv"
	"strings"
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
	var truth bool
	switch b {
	case "enabled":
		truth = true
	case "disabled":
		truth = false
	default:
		return nil
	}
	return &truth
}
