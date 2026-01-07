//go:build linux

package mem

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// GetTotalMemory returns total physical memory in bytes
func GetTotalMemory() (uint64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		if fields[0] == "MemTotal:" {
			val, _ := strconv.ParseUint(fields[1], 10, 64)
			return val * 1024, nil // Convert from KB to bytes
		}
	}
	return 0, nil
}
