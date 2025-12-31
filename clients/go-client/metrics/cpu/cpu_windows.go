//go:build windows

package cpu

import (
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// TODO: TEST WINDOWS CPU READING

var (
	modkernel32        = windows.NewLazySystemDLL("kernel32.dll")
	procGetSystemTimes = modkernel32.NewProc("GetSystemTimes")
)

func GetCpuPercent(interval time.Duration) (float64, error) {
	idle1, kernel1, user1, err := getSystemTimes()
	if err != nil {
		return 0, err
	}

	time.Sleep(interval)

	idle2, kernel2, user2, err := getSystemTimes()
	if err != nil {
		return 0, err
	}

	idleDelta := idle2 - idle1
	kernelDelta := kernel2 - kernel1
	userDelta := user2 - user1

	totalDelta := kernelDelta + userDelta
	if totalDelta == 0 {
		return 0, nil
	}

	// kernel time includes idle time
	return float64(totalDelta-idleDelta) / float64(totalDelta) * 100, nil
}

func getSystemTimes() (idle, kernel, user uint64, err error) {
	var idleTime, kernelTime, userTime windows.Filetime

	ret, _, err := procGetSystemTimes.Call(
		uintptr(unsafe.Pointer(&idleTime)),
		uintptr(unsafe.Pointer(&kernelTime)),
		uintptr(unsafe.Pointer(&userTime)),
	)
	if ret == 0 {
		return 0, 0, 0, err
	}

	return fileTimeToUint64(idleTime), fileTimeToUint64(kernelTime), fileTimeToUint64(userTime), nil
}

func fileTimeToUint64(ft windows.Filetime) uint64 {
	return uint64(ft.HighDateTime)<<32 | uint64(ft.LowDateTime)
}
