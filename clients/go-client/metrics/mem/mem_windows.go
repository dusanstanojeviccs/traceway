package mem

import (
	"syscall"
	"unsafe"
)

// GetTotalMemory returns total physical memory in bytes
func GetTotalMemory() (uint64, error) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("GlobalMemoryStatusEx")

	var mem struct {
		dwLength     uint32
		dwMemoryLoad uint32
		ullTotalPhys uint64
		ullAvailPhys uint64
		_            [4]uint64 // remaining fields we don't need
	}
	mem.dwLength = uint32(unsafe.Sizeof(mem))

	ret, _, err := proc.Call(uintptr(unsafe.Pointer(&mem)))
	if ret == 0 {
		return 0, err
	}
	return mem.ullTotalPhys, nil
}
