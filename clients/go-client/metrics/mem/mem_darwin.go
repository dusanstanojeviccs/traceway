//go:build darwin

package mem

import (
	"encoding/binary"
	"fmt"
	"syscall"
	"unsafe"
)

// GetTotalMemory returns total physical memory in bytes
func GetTotalMemory() (uint64, error) {
	return getTotalVirtualMemory()
}

func getTotalVirtualMemory() (uint64, error) {
	// Get physical memory size using sysctl
	physMem, err := syscall.Sysctl("hw.memsize")
	if err != nil {
		return 0, fmt.Errorf("failed to get hw.memsize: %w", err)
	}

	// The result is a binary representation of uint64
	// syscall.Sysctl returns a string, but the data is actually binary
	if len(physMem) < 8 {
		// Pad if necessary
		padded := make([]byte, 8)
		copy(padded, []byte(physMem))
		return binary.LittleEndian.Uint64(padded), nil
	}

	return binary.LittleEndian.Uint64([]byte(physMem)[:8]), nil
}

// Alternative using SysctlUint64 (cleaner approach)
func getTotalPhysicalMemory() (uint64, error) {
	mib := []int32{6, 24} // CTL_HW, HW_MEMSIZE

	var value uint64
	size := uintptr(unsafe.Sizeof(value))

	_, _, errno := syscall.Syscall6(
		syscall.SYS___SYSCTL,
		uintptr(unsafe.Pointer(&mib[0])),
		uintptr(len(mib)),
		uintptr(unsafe.Pointer(&value)),
		uintptr(unsafe.Pointer(&size)),
		0,
		0,
	)

	if errno != 0 {
		return 0, errno
	}

	return value, nil
}
