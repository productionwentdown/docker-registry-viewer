// reference: github.com/pivotal-golang/bytefmt
package client

import (
	"fmt"
	"strings"
)

const (
	s_BYTE     = 1.0
	s_KILOBYTE = 1024 * s_BYTE
	s_MEGABYTE = 1024 * s_KILOBYTE
	s_GIGABYTE = 1024 * s_MEGABYTE
	s_TERABYTE = 1024 * s_GIGABYTE
)

func humanSize(bytes uint64) string {
	unit := ""
	value := float32(bytes)

	switch {
	case bytes >= s_TERABYTE:
		unit = "T"
		value = value / s_TERABYTE
	case bytes >= s_GIGABYTE:
		unit = "G"
		value = value / s_GIGABYTE
	case bytes >= s_MEGABYTE:
		unit = "M"
		value = value / s_MEGABYTE
	case bytes >= s_KILOBYTE:
		unit = "K"
		value = value / s_KILOBYTE
	case bytes >= s_BYTE:
		unit = "B"
	case bytes == 0:
		return "0"
	}

	s := fmt.Sprintf("%.1f", value)
	s = strings.TrimSuffix(s, ".0")
	return fmt.Sprintf("%s%s", s, unit)
}
