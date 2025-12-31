//go:build linux

package mem

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// TODO: TEST LINUX MEMORY READING

// Linux: read /proc/meminfo
func GetMemoryUsedPercent() (float64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var total, available uint64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		val, _ := strconv.ParseUint(fields[1], 10, 64)
		switch fields[0] {
		case "MemTotal:":
			total = val
		case "MemAvailable:":
			available = val
		}
	}
	return float64(total-available) / float64(total) * 100, nil
}
